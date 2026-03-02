// Package claw402 provides a typed SDK for claw402.ai — pay-per-call crypto data APIs via x402.
//
// Usage:
//
//	client, err := claw402.New("0xPrivateKey")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	data, err := client.Coinank.Fund.Realtime(ctx, &claw402.CoinankFundRealtimeParams{
//	    ProductType: "SWAP",
//	})
package claw402

import (
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

	Coinank *CoinankResource
	Nofxos  *NofxosResource
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
	c.Coinank = newCoinankResource(c)
	c.Nofxos = newNofxosResource(c)
	return c, nil
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
