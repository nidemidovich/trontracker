package service

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	tele "gopkg.in/telebot.v3"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/nidemidovich/trontracker/internal/commands/hello"
	"github.com/nidemidovich/trontracker/internal/config"
	"github.com/nidemidovich/trontracker/internal/infrastructure/telegram"
	"github.com/nidemidovich/trontracker/internal/infrastructure/tron"
	"github.com/nidemidovich/trontracker/internal/parser"
)

func Run() error {
	botConfig, err := config.NewBot()
	if err != nil {
		log.Fatalf("error init bot config: %s", err)
	}

	tronGridConfig, err := config.NewTronGrid()
	if err != nil {
		log.Fatalf("error init tron grid config: %s", err)
	}

	curDir, _ := os.Getwd()

	pumpSwapRouterABIFile, err := os.ReadFile(filepath.Join(curDir, "abis/pump_swap_router.json"))
	if err != nil {
		log.Fatalf("error reading abi from file: %s", err)
	}

	pumpSwapRouterABI, err := abi.JSON(bytes.NewReader(pumpSwapRouterABIFile))
	if err != nil {
		log.Fatalf("error parsing abi: %s", err)
	}

	uniswapABIFile, err := os.ReadFile(filepath.Join(curDir, "abis/uniswap_abi.json"))
	if err != nil {
		log.Fatalf("error reading abi from file: %s", err)
	}

	uniswapABI, err := abi.JSON(bytes.NewReader(uniswapABIFile))
	if err != nil {
		log.Fatalf("error parsing abi: %s", err)
	}

	trc20ABIFile, err := os.ReadFile(filepath.Join(curDir, "abis/trc20_abi.json"))
	if err != nil {
		log.Fatalf("error reading abi from file: %s", err)
	}

	trc20ABI, err := abi.JSON(bytes.NewReader(trc20ABIFile))
	if err != nil {
		log.Fatalf("error parsing abi: %s", err)
	}

	pref := tele.Settings{
		Token:  botConfig.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatalf("error init bot: %s", err)
	}

	tronClient := tron.NewClient(http.Client{}, tronGridConfig.APIKey)

	parser := parser.New(*tronClient, b, pumpSwapRouterABI, uniswapABI, trc20ABI)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	b.Handle("/hello", func(c tele.Context) error {
		tgCtx := telegram.NewContext(c)
		return hello.New().Handle(tgCtx)
	})

	go func() {
		b.Start()
	}()

	go func(ctx context.Context) {
		for {
			parser.ParseBlock(ctx)

			timer := time.NewTimer(time.Second * 3)
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
			}
		}
	}(ctx)

	log.Println("App started")

	log.Printf("Got signal %v, attempting graceful shutdown", <-quit)

	cancel()

	b.Stop()

	log.Println("App shutting down")

	return nil
}
