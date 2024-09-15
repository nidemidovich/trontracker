package tronscan

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/nidemidovich/trontracker/internal/dto/tronscan"
)

const (
	baseURL = "https://apilist.tronscanapi.com"

	apiKeyHeaderName = "TRON-PRO-API-KEY"
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

func (c *Client) GetToken(ctx context.Context, tokenAddress string) (tronscan.GetTokenResponseDto, error) {
	path := "/api/token_trc20"

	params := url.Values{}
	params.Add("contract", tokenAddress)

	u, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return tronscan.GetTokenResponseDto{}, err
	}

	u.Path = path
	u.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return tronscan.GetTokenResponseDto{}, err
	}

	req.Header.Add(apiKeyHeaderName, c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return tronscan.GetTokenResponseDto{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return tronscan.GetTokenResponseDto{}, ErrFailedRequest
	}

	var out tronscan.GetTokenResponseDto
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return tronscan.GetTokenResponseDto{}, err
	}

	return out, nil
}
