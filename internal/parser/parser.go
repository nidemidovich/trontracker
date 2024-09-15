package parser

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	tele "gopkg.in/telebot.v3"

	tron_dto "github.com/nidemidovich/trontracker/internal/dto/tron"
	"github.com/nidemidovich/trontracker/internal/infrastructure/tron"
)

const (
	pumpSwapContract string = "TZFs5ch1R1C4mmjwrrmZqeqbUgGpxY1yWB"
	swapAbiSignature string = "d78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822"
)

type Parser struct {
	tronClient        tron.Client
	bot               *tele.Bot
	pumpSwapRouterABI abi.ABI
	uniswapABI        abi.ABI
	trc20ABI          abi.ABI
}

func New(client tron.Client, bot *tele.Bot, pumpSwapRouterABI abi.ABI, uniswapABI abi.ABI, trc20ABI abi.ABI) *Parser {
	return &Parser{
		tronClient:        client,
		bot:               bot,
		pumpSwapRouterABI: pumpSwapRouterABI,
		uniswapABI:        uniswapABI,
		trc20ABI:          trc20ABI,
	}
}

func (p *Parser) ParseBlock(ctx context.Context) error {
	blockOut, err := p.tronClient.GetBlock(ctx)
	if err != nil {
		log.Printf("error while retrieving block: %s", err)
		return err
	}

	blockInfoOut, err := p.tronClient.GetBlockInfo(ctx, blockOut.BlockHeader.RawData.Number)
	if err != nil {
		log.Printf("error while retrieving block info: %s", err)
		return err
	}

	var wg sync.WaitGroup

	messagesChan := make(chan string, len(blockOut.Transactions))

	defer close(messagesChan)

	quit := make(chan struct{})

	defer close(quit)

	go func(ctx context.Context) {
		for {
			select {
			case msg := <-messagesChan:
				if _, err := p.bot.Send(user{ChatID: "292177347"}, msg, &tele.SendOptions{
					ParseMode:             tele.ModeMarkdownV2,
					DisableWebPagePreview: true,
				}); err != nil {
					log.Printf("error while sending message: %s", err)
				}
				wg.Done()
			case <-ctx.Done():
				return
			case <-quit:
				return
			default:
			}
		}
	}(ctx)

	for i, trx := range blockOut.Transactions {
		if blockInfoOut[i].Result != nil {
			continue
		}

		if len(trx.RawData.Contract) == 0 {
			continue
		}

		c := trx.RawData.Contract[0]
		if c.Type == "TriggerSmartContract" {
			var param tron_dto.TriggerSmartContractParameter
			if err := json.NewDecoder(bytes.NewBuffer(c.Parameter)).Decode(&param); err != nil {
				log.Printf("error while decoding parameter: %s", err)
				continue
			}

			if ok, err := compareAddresses(param.Value.ContractAddress); !ok || err != nil {
				continue
			}

			decodedSelector, err := hex.DecodeString(param.Value.Data[:8])
			if err != nil {
				log.Printf("error while decoding selector: %s", err)
				continue
			}

			method, err := p.pumpSwapRouterABI.MethodById(decodedSelector)
			if err != nil {
				log.Printf("error while retrieving method: %s", err)
				continue
			}

			decodedData, err := hex.DecodeString(param.Value.Data[8:])
			if err != nil {
				log.Printf("error while decoding args: %s", err)
				continue
			}

			args, err := method.Inputs.Unpack(decodedData)
			if err != nil {
				log.Printf("error while unpacking args: %s", err)
				continue
			}

			swapRes, err := p.getSwapResultFromEvent(blockInfoOut[i])
			if err != nil {
				continue
			}

			var msg string

			switch method.Name {

			// Sell
			case "swapTokensForExactETH", "swapExactTokensForETH", "swapExactTokensForETHSupportingFeeOnTransferTokens":
				soldTokenETHAddr := args[2].([]common.Address)[0]
				soldTokenTronAddr := address.Address(append([]byte{byte(0x41)}, soldTokenETHAddr.Bytes()...))

				symbol, err := p.getSwappedTokenSymbol(ctx, soldTokenTronAddr)
				if err != nil {
					log.Printf("error while retrieving token symbol: %s", err)
					continue
				}

				msg = p.prepareSellMessage(args, swapRes, trx.TxID, soldTokenTronAddr, symbol)

			// Buy
			case "swapExactETHForTokens", "swapETHForExactTokens", "swapExactETHForTokensSupportingFeeOnTransferTokens":
				boughtTokenETHAddr := args[1].([]common.Address)[1]
				boughtTokenTronAddr := address.Address(append([]byte{byte(0x41)}, boughtTokenETHAddr.Bytes()...))

				symbol, err := p.getSwappedTokenSymbol(ctx, boughtTokenTronAddr)
				if err != nil {
					log.Printf("error while retrieving token symbol: %s", err)
					continue
				}

				msg = p.prepareBuyMessage(args, swapRes, trx.TxID, boughtTokenTronAddr, symbol)
			default:
				continue
			}

			wg.Add(1)

			messagesChan <- msg
		}
	}

	wg.Wait()

	quit <- struct{}{}

	return nil
}

type user struct {
	ChatID string
}

func (u user) Recipient() string {
	return u.ChatID
}

func compareAddresses(addr string) (bool, error) {
	triggeredAddr := address.HexToAddress(addr)

	pumpSwapAddr, err := address.Base58ToAddress(pumpSwapContract)
	if err != nil {
		return false, err
	}

	return bytes.Equal(triggeredAddr.Bytes(), pumpSwapAddr.Bytes()), nil
}

type sell struct {
	Seller             string
	TokenSymbol        string
	TokenAddr          string
	SoldTokenAmount    *big.Int
	RecievedTronAmount *big.Int
	TrxID              string
}

func (s *sell) String() string {
	return fmt.Sprintf(
		"🔴SELL\n\n`%s` [sold](https://tronscan.org/#/transaction/%s) %s of [%s](https://tronscan.org/#/token20/%s) for %s",
		s.Seller,
		s.TrxID,
		s.SoldTokenAmount,
		s.TokenSymbol,
		s.TokenAddr,
		s.RecievedTronAmount,
	)
}

type buy struct {
	Buyer             string
	TokenAddr         string
	TokenSymbol       string
	TrxID             string
	BoughtTokenAmount *big.Int
	SpentTronAmount   *big.Int
}

func (b *buy) String() string {
	return fmt.Sprintf(
		"🟢BUY\n\n`%s` [bought](https://tronscan.org/#/transaction/%s) %s of [%s](https://tronscan.org/#/token20/%s) for %s",
		b.Buyer,
		b.TrxID,
		b.BoughtTokenAmount,
		b.TokenSymbol,
		b.TokenAddr,
		b.SpentTronAmount,
	)
}

func (p *Parser) getSwapResultFromEvent(trxInfo tron_dto.TransactionInfo) ([]*big.Int, error) {
	var data string

	for _, logEntry := range trxInfo.Log {
		for _, topic := range logEntry.Topics {
			if topic == swapAbiSignature {
				data = logEntry.Data
			}
		}
	}

	if data == "" {
		msg := "swap event not found in logs"
		log.Println(msg)
		return nil, errors.New(msg)
	}

	decodedId, err := hex.DecodeString(swapAbiSignature)
	if err != nil {
		log.Printf("error while decoding swap abi signature: %s", err)
		return nil, err
	}

	event, err := p.uniswapABI.EventByID(common.BytesToHash(decodedId))
	if err != nil {
		log.Printf("error while retrieving event: %s", err)
		return nil, err
	}

	decodedData, err := hex.DecodeString(data)
	if err != nil {
		log.Printf("error while decoding event data: %s", err)
		return nil, err
	}

	eventArgs, err := event.Inputs.Unpack(decodedData)
	if err != nil {
		log.Printf("error while unpacking event args: %s", err)
		return nil, err
	}

	var in, out *big.Int

	in0, in1 := eventArgs[0].(*big.Int), eventArgs[1].(*big.Int)
	if in0.Cmp(in1) == 1 {
		in = in0
	} else {
		in = in1
	}

	out0, out1 := eventArgs[2].(*big.Int), eventArgs[3].(*big.Int)
	if out0.Cmp(out1) == 1 {
		out = out0
	} else {
		out = out1
	}

	return []*big.Int{in, out}, nil
}

func (p *Parser) getSwappedTokenSymbol(ctx context.Context, tokenAddress address.Address) (string, error) {
	out, err := p.tronClient.GetTokenSymbol(ctx, tokenAddress.String())
	if err != nil {
		return "", err
	}

	res := out.ConstantResult[0]

	decodedRes, err := hex.DecodeString(res)
	if err != nil {
		return "", err
	}

	symbol, err := p.trc20ABI.Methods["symbol"].Outputs.Unpack(decodedRes)
	if err != nil {
		return "", err
	}

	return symbol[0].(string), nil
}

func (p *Parser) prepareSellMessage(methodArgs []interface{}, swapRes []*big.Int, trxID string, tokenAddr address.Address, symbol string) string {
	soldAmount := swapRes[0]
	receivedTronAmount := swapRes[1]

	sellerETHAddr := methodArgs[3].(common.Address)
	sellerTronAddr := address.Address(append([]byte{byte(0x41)}, sellerETHAddr.Bytes()...))

	s := sell{
		Seller:             sellerTronAddr.String(),
		TokenAddr:          tokenAddr.String(),
		TokenSymbol:        symbol,
		SoldTokenAmount:    soldAmount,
		RecievedTronAmount: receivedTronAmount,
		TrxID:              trxID,
	}

	return s.String()
}

func (p *Parser) prepareBuyMessage(methodArgs []interface{}, swapRes []*big.Int, trxID string, tokenAddr address.Address, symbol string) string {
	boughtAmount := swapRes[1]
	spentTronAmount := swapRes[0]

	buyerETHAddr := methodArgs[2].(common.Address)
	buyerTronAddr := address.Address(append([]byte{byte(0x41)}, buyerETHAddr.Bytes()...))

	b := buy{
		Buyer:             buyerTronAddr.String(),
		TokenAddr:         tokenAddr.String(),
		TokenSymbol:       symbol,
		TrxID:             trxID,
		BoughtTokenAmount: boughtAmount,
		SpentTronAmount:   spentTronAmount,
	}

	return b.String()
}
