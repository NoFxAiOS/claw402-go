# Architecture

## Overview

```
claw402-go (go get)
├── client.go              ← Core client: x402 V2 payment via Coinbase Go lib
├── errors.go              ← Error types
├── coinank.go             ← 78 endpoints: market data, ETF, liquidations, etc.
├── nofxos.go              ← 18 endpoints: AI signals, rankings, Upbit
└── examples/
    └── basic/main.go      ← Usage example
```

## Payment Flow (x402 V2)

```
Client                          claw402.ai                    Base L2
  │                                │                            │
  │─── GET /api/v1/... ──────────▶│                            │
  │◀── 402 + Payment-Required ────│                            │
  │                                │                            │
  │  [sign EIP-3009 locally]       │                            │
  │                                │                            │
  │─── GET + PAYMENT-SIGNATURE ──▶│                            │
  │                                │── verify + settle ────────▶│
  │                                │◀── tx confirmed ──────────│
  │◀── 200 + data ────────────────│                            │
```

The Go SDK uses `github.com/coinbase/x402/go` which natively implements the
V2 protocol, wrapping `http.Client` with automatic payment handling.

## Code Generation

The generated files (`coinank.go`, `nofxos.go`) are produced by `sdks/codegen/`
which reads `providers/*.yaml` and emits typed SDK methods for Go, TypeScript,
and Python.

Each YAML route becomes a typed method:

```yaml
# providers/coinank.yaml
- gateway_path: /api/v1/coinank/fund/realtime
  category: Fund
  allowed_params: [sortBy, productType, page, size]
```

Becomes:

```go
// coinank.go
type CoinankFundRealtimeParams struct {
    SortBy      string
    ProductType string
    Page        string
    Size        string
}

func (r *CoinankFundResource) Realtime(ctx context.Context, p *CoinankFundRealtimeParams) (json.RawMessage, error) {
    params := map[string]string{}
    if p != nil { params["sortBy"] = p.SortBy; ... }
    return r.client.get(ctx, "/api/v1/coinank/fund/realtime", params)
}
```
