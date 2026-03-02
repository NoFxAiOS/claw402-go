# claw402-go

[![Go Reference](https://pkg.go.dev/badge/github.com/NoFxAiOS/claw402-go.svg)](https://pkg.go.dev/github.com/NoFxAiOS/claw402-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Typed Go SDK for [claw402.ai](https://claw402.ai) — pay-per-call crypto data APIs via [x402](https://www.x402.org/) micropayments.

**96+ endpoints** covering fund flow, liquidations, ETF flows, AI trading signals, whale tracking, funding rates, open interest, and more. No API key, no signup, no subscription — just a Base wallet with USDC.

## Install

```bash
go get github.com/NoFxAiOS/claw402-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    claw402 "github.com/NoFxAiOS/claw402-go"
)

func main() {
    client, err := claw402.New("0xYourPrivateKey")
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Fund flow — $0.001 per call
    flow, err := client.Coinank.Fund.Realtime(ctx, &claw402.CoinankFundRealtimeParams{
        ProductType: "SWAP",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(flow))

    // AI trading signals
    signals, err := client.Nofxos.Netflow.TopRanking(ctx, &claw402.NofxosNetflowTopRankingParams{
        Limit:    "20",
        Duration: "1h",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(signals))
}
```

## Features

- **Typed methods** — every endpoint has a dedicated Go function with typed params
- **Automatic x402 payment** — signs EIP-3009 USDC transfers locally, never sends your key
- **Two resource groups** — `client.Coinank.*` (market data) and `client.Nofxos.*` (AI signals)
- **Base mainnet** — pays $0.001 USDC per call on Coinbase L2

## API Overview

### Coinank (Market Data)

| Resource | Methods | Description |
|----------|---------|-------------|
| `Coinank.Fund` | `Realtime`, `History` | Real-time & historical fund flow |
| `Coinank.Oi` | `All`, `AggChart`, `SymbolChart`, `Kline`, ... | Open interest data |
| `Coinank.Liquidation` | `Orders`, `Intervals`, `AggHistory`, `LiqMap`, `HeatMap`, ... | Liquidation tracking |
| `Coinank.FundingRate` | `Current`, `Accumulated`, `Hist`, `Weighted`, `Heatmap`, ... | Funding rate analytics |
| `Coinank.Longshort` | `Realtime`, `BuySell`, `Person`, `Position`, ... | Long/short ratios |
| `Coinank.Hyper` | `TopPosition`, `TopAction` | HyperLiquid whale tracking |
| `Coinank.Etf` | `UsBtc`, `UsEth`, `UsBtcInflow`, `UsEthInflow`, `HkInflow` | ETF flow data |
| `Coinank.Indicator` | `FearGreed`, `AltcoinSeason`, `BtcMultiplier`, `Ahr999`, ... | Market cycle indicators |
| `Coinank.MarketOrder` | `Cvd`, `AggCvd`, `BuySellValue`, ... | Taker flow / CVD |
| `Coinank.Kline` | `Lists` | OHLCV candlestick data |
| `Coinank.Price` | `Last` | Real-time price |
| `Coinank.Rank` | `Screener`, `Oi`, `Volume`, `Price`, `Liquidation`, ... | Rankings & screeners |
| `Coinank.News` | `List`, `Detail` | Crypto news & alerts |

### Nofxos (AI Signals)

| Resource | Methods | Description |
|----------|---------|-------------|
| `Nofxos.Ai500` | `List`, `Stats` | AI500 high-potential coin signals |
| `Nofxos.Ai300` | `List`, `Stats` | AI300 quant model rankings |
| `Nofxos.Netflow` | `TopRanking`, `LowRanking` | Net capital flow rankings |
| `Nofxos.Oi` | `TopRanking`, `LowRanking` | OI change rankings |
| `Nofxos.FundingRate` | `Top`, `Low` | Extreme funding rate coins |
| `Nofxos.Price` | `Ranking` | Price change rankings |
| `Nofxos.Upbit` | `Hot`, `NetflowTopRanking`, `NetflowLowRanking` | Korean market data |

## Configuration

```go
// Custom base URL
client, err := claw402.New("0xKey", claw402.WithBaseURL("https://custom.gateway"))
```

## How Payment Works

1. SDK sends a GET request to the endpoint
2. Server responds with `402 Payment Required` + payment details in header
3. SDK signs an EIP-3009 `TransferWithAuthorization` for USDC on Base
4. SDK retries the request with the `PAYMENT-SIGNATURE` header
5. Server verifies payment on-chain and returns the data

Your private key **never leaves your machine** — it only signs the payment locally.

## Requirements

- Go 1.21+
- A wallet with USDC on [Base mainnet](https://base.org)
- Get USDC on Base: [bridge.base.org](https://bridge.base.org)

## License

MIT
