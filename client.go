// Package claw402 provides a typed SDK for claw402.ai — pay-per-call data APIs via x402.
//
// All API calls cost micro-amounts of USDC on Base mainnet (Coinbase L2).
// No API key, no account, no subscription needed — just a wallet with USDC.
//
// Usage:
//
//	client, err := claw402.New("0xYourPrivateKey")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Crypto market data
//	data, err := client.Coinank.Fund.Realtime(ctx, &claw402.CoinankFundRealtimeParams{
//	    ProductType: "SWAP",
//	})
//
//	// US stocks
//	quote, err := client.Alphavantage.StocksUs.Quote(ctx, &claw402.AlphavantageStocksUsQuoteParams{
//	    Symbol: "AAPL",
//	})
//
//	// AI models (POST endpoints)
//	resp, err := client.Openai.Openai.Chat(ctx, map[string]interface{}{
//	    "model": "gpt-4o",
//	    "messages": []map[string]string{{"role": "user", "content": "Hello!"}},
//	})
//
//	// Forex & metals
//	price, err := client.Twelvedata.Price.Price(ctx, &claw402.TwelvedataPricePriceParams{
//	    Symbol: "EUR/USD",
//	})
package claw402

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	x402 "github.com/coinbase/x402/go"
	x402http "github.com/coinbase/x402/go/http"
	evmclient "github.com/coinbase/x402/go/mechanisms/evm/exact/client"
	evmsigner "github.com/coinbase/x402/go/signers/evm"
)

// Client is the top-level claw402 SDK client.
type Client struct {
	http    *http.Client
	baseURL string

	// Crypto data providers
	Coinank       *CoinankResource
	Nofxos        *NofxosResource
	Coinmarketcap *CoinmarketcapResource
	// US stock market data
	Alphavantage *AlphavantageResource
	Polygon      *PolygonResource
	Alpaca       *AlpacaResource
	// Chinese A-share market data
	Tushare *TushareResource
	// Forex, metals, indices
	Twelvedata *TwelvedataResource
	// AI model providers
	Openai    *OpenaiResource
	Anthropic *AnthropicResource
	Deepseek  *DeepseekResource
	Qwen      *QwenResource
	Gemini    *GeminiResource
	Grok      *GrokResource
	Kimi      *KimiResource
	// Web3 intelligence
	Rootdata *RootdataResource
}

// New creates a new claw402 client with the given private key.
func New(privateKey string, opts ...Option) (*Client, error) {
	cfg := &config{baseURL: "https://claw402.ai"}
	for _, o := range opts {
		o(cfg)
	}

	signer, err := evmsigner.NewClientSignerFromPrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("claw402: invalid private key: %w", err)
	}

	x402Client := x402.Newx402Client().
		Register(x402.Network("eip155:8453"), evmclient.NewExactEvmScheme(signer))
	httpClient := x402http.WrapHTTPClientWithPayment(
		&http.Client{Timeout: 30 * time.Second},
		x402http.Newx402HTTPClient(x402Client),
	)

	c := &Client{
		http:    httpClient,
		baseURL: cfg.baseURL,
	}
	// Crypto data
	c.Coinank = newCoinankResource(c)
	c.Nofxos = newNofxosResource(c)
	c.Coinmarketcap = newCoinmarketcapResource(c)
	// US stocks
	c.Alphavantage = newAlphavantageResource(c)
	c.Polygon = newPolygonResource(c)
	c.Alpaca = newAlpacaResource(c)
	// Chinese A-shares
	c.Tushare = newTushareResource(c)
	// Forex, metals, indices
	c.Twelvedata = newTwelvedataResource(c)
	// AI models
	c.Openai = newOpenaiResource(c)
	c.Anthropic = newAnthropicResource(c)
	c.Deepseek = newDeepseekResource(c)
	c.Qwen = newQwenResource(c)
	c.Gemini = newGeminiResource(c)
	c.Grok = newGrokResource(c)
	c.Kimi = newKimiResource(c)
	// Web3 intelligence
	c.Rootdata = newRootdataResource(c)
	return c, nil
}

func (c *Client) post(ctx context.Context, path string, body map[string]interface{}) (json.RawMessage, error) {
	urlStr := c.baseURL + path

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("claw402: marshal body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlStr, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("claw402: new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("claw402: request: %w", err)
	}
	defer resp.Body.Close()

	const maxBodySize = 10 * 1024 * 1024 // 10 MB
	respBody, err := io.ReadAll(io.LimitReader(resp.Body, maxBodySize))
	if err != nil {
		return nil, fmt.Errorf("claw402: read body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, &Error{Status: resp.StatusCode, Body: string(respBody)}
	}

	return json.RawMessage(respBody), nil
}

func (c *Client) get(ctx context.Context, path string, params map[string]string) (json.RawMessage, error) {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("claw402: bad url: %w", err)
	}

	if params != nil {
		q := u.Query()
		for k, v := range params {
			if v != "" {
				q.Set(k, v)
			}
		}
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("claw402: new request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("claw402: request: %w", err)
	}
	defer resp.Body.Close()

	const maxBodySize = 10 * 1024 * 1024 // 10 MB
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBodySize))
	if err != nil {
		return nil, fmt.Errorf("claw402: read body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, &Error{Status: resp.StatusCode, Body: string(body)}
	}

	return json.RawMessage(body), nil
}

// Option configures the claw402 client.
type Option func(*config)

type config struct {
	baseURL string
}

// WithBaseURL sets a custom base URL for the claw402 gateway.
func WithBaseURL(u string) Option {
	return func(c *config) { c.baseURL = u }
}
