package tron

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/nidemidovich/trontracker/internal/dto/tron"
)

const (
	baseURL = "https://api.trongrid.io"

	apiKeyHeaderName = "TRON-PRO-API-KEY"

	ownerAddress = "TWST6WByrwwo132hqoKQUe9zzj5UXmBRtv"
)

var ErrFailedRequest = errors.New("failed request")

type Client struct {
	client http.Client
	apiKey string
}

func NewClient(client http.Client, apiKey string) *Client {
	return &Client{
		client: client,
		apiKey: apiKey,
	}
}

func (c *Client) GetBlock(ctx context.Context) (tron.GetBlockResponseDto, error) {
	path := "/wallet/getblock"

	in := tron.GetBlockRequestDto{
		Detail: true,
	}

	body, err := json.Marshal(in)
	if err != nil {
		return tron.GetBlockResponseDto{}, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", baseURL+path, bytes.NewReader(body))
	if err != nil {
		return tron.GetBlockResponseDto{}, err
	}

	req.Header.Add(apiKeyHeaderName, c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return tron.GetBlockResponseDto{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return tron.GetBlockResponseDto{}, ErrFailedRequest
	}

	var out tron.GetBlockResponseDto
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return tron.GetBlockResponseDto{}, err
	}

	return out, nil
}

func (c *Client) GetBlockInfo(ctx context.Context, blockNum int64) (tron.GetTransactionInfoByBlockNumResponseDto, error) {
	path := "/wallet/gettransactioninfobyblocknum"

	in := tron.GetTransactionInfoByBlockNumRequestDto{
		Num: blockNum,
	}

	body, err := json.Marshal(in)
	if err != nil {
		return tron.GetTransactionInfoByBlockNumResponseDto{}, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", baseURL+path, bytes.NewReader(body))
	if err != nil {
		return tron.GetTransactionInfoByBlockNumResponseDto{}, err
	}

	req.Header.Add(apiKeyHeaderName, c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return tron.GetTransactionInfoByBlockNumResponseDto{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return tron.GetTransactionInfoByBlockNumResponseDto{}, ErrFailedRequest
	}

	var out tron.GetTransactionInfoByBlockNumResponseDto
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return tron.GetTransactionInfoByBlockNumResponseDto{}, err
	}

	return out, nil
}

func (c *Client) GetTokenSymbol(ctx context.Context, tokenAddress string) (tron.GetTokenSymbolResponseDto, error) {
	path := "/wallet/triggerconstantcontract"

	in := tron.GetTokenAddressRequestDto{
		ContractAddress:  tokenAddress,
		FunctionSelector: "symbol()",
		OwnerAddress:     ownerAddress,
		Visible:          true,
	}

	body, err := json.Marshal(in)
	if err != nil {
		return tron.GetTokenSymbolResponseDto{}, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", baseURL+path, bytes.NewReader(body))
	if err != nil {
		return tron.GetTokenSymbolResponseDto{}, err
	}

	req.Header.Add(apiKeyHeaderName, c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return tron.GetTokenSymbolResponseDto{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return tron.GetTokenSymbolResponseDto{}, ErrFailedRequest
	}

	var out tron.GetTokenSymbolResponseDto
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return tron.GetTokenSymbolResponseDto{}, err
	}

	return out, nil
}
