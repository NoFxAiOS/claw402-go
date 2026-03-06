# claw402-go

[![Go Reference](https://pkg.go.dev/badge/github.com/NoFxAiOS/claw402-go.svg)](https://pkg.go.dev/github.com/NoFxAiOS/claw402-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Typed Go SDK for [claw402.ai](https://claw402.ai) — pay-per-call data APIs via [x402](https://www.x402.org/) micropayments.

**200+ endpoints** covering crypto market data, US stocks, China A-shares, forex, global time-series, and AI (OpenAI/Anthropic/DeepSeek/Qwen). No API key, no signup, no subscription — just a Base wallet with USDC.

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

    // Crypto: Fund flow — $0.001/call
    flow, _ := client.Coinank.Fund.Realtime(ctx, &claw402.CoinankFundRealtimeParams{
        ProductType: "SWAP",
    })
    fmt.Println(string(flow))

    // US Stocks: Latest quote — $0.001/call
    quote, _ := client.Alpaca.Quotes.Latest(ctx, &claw402.AlpacaQuotesLatestParams{
        Symbols: "AAPL,TSLA",
    })
    fmt.Println(string(quote))

    // China A-shares — $0.001/call  (sub-resource is Cn, not Tushare)
    stocks, _ := client.Tushare.Cn.StockBasic(ctx, &claw402.TushareCnStockBasicParams{
        ListStatus: "L",
    })
    fmt.Println(string(stocks))

    // Forex time-series — $0.001/call  (use GetTimeSeries, TimeSeries is a sub-resource)
    ts, _ := client.Twelvedata.GetTimeSeries(ctx, &claw402.TwelvedataGetTimeSeriesParams{
        Symbol:   "EUR/USD",
        Interval: "1h",
    })
    fmt.Println(string(ts))

    // AI: OpenAI chat — $0.01/call
    resp, _ := client.Openai.Openai.Chat(ctx, map[string]interface{}{
        "messages": []map[string]string{{"role": "user", "content": "Hello"}},
    })
    fmt.Println(string(resp))
}
```

## Features

- **Typed methods** — every endpoint has a dedicated Go function with typed params struct
- **Automatic x402 payment** — signs EIP-3009 USDC transfers locally, never sends your key
- **11 provider groups** — crypto, US stocks, China stocks, forex, global data, and AI
- **Base mainnet** — pays USDC per call on Coinbase L2

## API Overview

### Crypto Market Data

#### Coinank

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

#### Nofxos (AI Signals)

| Resource | Methods | Description |
|----------|---------|-------------|
| `Nofxos.Ai500` | `List`, `Stats` | AI500 high-potential coin signals |
| `Nofxos.Ai300` | `List`, `Stats` | AI300 quant model rankings |
| `Nofxos.Netflow` | `TopRanking`, `LowRanking` | Net capital flow rankings |
| `Nofxos.Oi` | `TopRanking`, `LowRanking` | OI change rankings |
| `Nofxos.FundingRate` | `Top`, `Low` | Extreme funding rate coins |
| `Nofxos.Price` | `Ranking` | Price change rankings |
| `Nofxos.Upbit` | `Hot`, `NetflowTopRanking`, `NetflowLowRanking` | Korean market data |

### US Stock & Options Market

#### Alpaca

| Resource | Methods | Description |
|----------|---------|-------------|
| `Alpaca.Quotes` | `Latest`, `History` | Real-time & historical quotes — $0.001–0.002/call |
| `Alpaca.Bars` | `Latest` | Latest OHLCV bar — $0.001/call |
| `Alpaca.Trades` | `Latest`, `History` | Real-time & historical trades — $0.001–0.002/call |
| `Alpaca.Options` | `Bars`, `QuotesLatest`, `Snapshots` | Options chain data — $0.003/call |
| `Alpaca` | `GetBars`, `Snapshots` (multi), `Snapshot` (single), `Movers`, `MostActives`, `News`, `CorporateActions` | Direct market endpoints — $0.001–0.002/call |

```go
// Latest quotes
quote, _ := client.Alpaca.Quotes.Latest(ctx, &claw402.AlpacaQuotesLatestParams{
    Symbols: "AAPL,MSFT,TSLA",
})

// Historical bars (GetBars, not Bars, to avoid collision with Alpaca.Bars sub-resource)
bars, _ := client.Alpaca.GetBars(ctx, &claw402.AlpacaGetBarsParams{
    Symbols:   "AAPL",
    Timeframe: "1Day",
    Start:     "2024-01-01",
})

// Top movers
movers, _ := client.Alpaca.Movers(ctx, &claw402.AlpacaMoversParams{
    Top:        "10",
    MarketType: "stocks",
})
```

#### Polygon

| Resource | Methods | Description |
|----------|---------|-------------|
| `Polygon.Aggs` | `Ticker`, `Grouped` | Aggregates / OHLCV bars — $0.001/call |
| `Polygon.Snapshot` | `Ticker`, `All`, `Gainers`, `Losers` | Full market snapshots — $0.002/call |
| `Polygon.Options` | `Contracts`, `ContractDetails`, `Snapshot`, `SnapshotContract` | Options data — $0.002/call |
| `Polygon` | `TickerDetails`, `MarketStatus`, `TickerNews`, `Tickers`, `Exchanges`, `PrevClose`, `Trades`, `LastTrade`, `Quotes`, `LastQuote`, `Sma`, `Ema`, `Rsi`, `Macd` | Reference & technical indicators — $0.001–0.003/call |

```go
// OHLCV bars  (method is Ticker, not Aggs)
bars, _ := client.Polygon.Aggs.Ticker(ctx, &claw402.PolygonAggsTickerParams{
    StocksTicker: "AAPL",
    Multiplier:   "1",
    Timespan:     "day",
    From:         "2024-01-01",
    To:           "2024-01-05",
})

// RSI indicator
rsi, _ := client.Polygon.Rsi(ctx, &claw402.PolygonRsiParams{
    StockTicker: "AAPL",
    Timespan:    "day",
    Window:      "14",
})
```

#### Alpha Vantage

| Resource | Methods | Description |
|----------|---------|-------------|
| `Alphavantage.Us` | `Quote`, `Search`, `Daily`, `DailyAdjusted`, `Intraday`, `Weekly`, `Monthly`, `Overview`, `Earnings`, `Income`, `BalanceSheet`, `CashFlow`, `Movers`, `News`, `Rsi`, `Macd`, `Bbands`, `Sma`, `Ema` | Comprehensive financial data — $0.001–0.003/call |

```go
// Real-time quote
quote, _ := client.Alphavantage.Us.Quote(ctx, &claw402.AlphavantageUsQuoteParams{
    Symbol: "AAPL",
})

// Daily OHLCV
daily, _ := client.Alphavantage.Us.Daily(ctx, &claw402.AlphavantageUsDailyParams{
    Symbol:     "AAPL",
    Outputsize: "compact",
})

// Top movers (no params)
movers, _ := client.Alphavantage.Us.Movers(ctx)
```

### China A-Shares

#### Tushare

| Resource | Methods | Description |
|----------|---------|-------------|
| `Tushare.Cn` | `StockBasic`, `Daily`, `Weekly`, `Monthly`, `DailyBasic`, `TradeCal`, `Income`, `BalanceSheet`, `CashFlow`, `Dividend`, `Northbound`, `Moneyflow`, `Margin`, `MarginDetail`, `TopList`, `TopInst` | China A-share market data — $0.001–0.003/call |

```go
// Stock list  (sub-resource is Cn, not Tushare)
stocks, _ := client.Tushare.Cn.StockBasic(ctx, &claw402.TushareCnStockBasicParams{
    ListStatus: "L",
})

// Daily OHLCV
daily, _ := client.Tushare.Cn.Daily(ctx, &claw402.TushareCnDailyParams{
    TsCode:    "000001.SZ",
    StartDate: "20240101",
    EndDate:   "20240131",
})
```

### Global Time-Series & Forex

#### Twelve Data

| Resource | Methods | Description |
|----------|---------|-------------|
| `Twelvedata.TimeSeries` | `Complex` (POST) | Complex multi-symbol/indicator query — $0.005/call |
| `Twelvedata.Indicator` | `Sma`, `Ema`, `Rsi`, `Macd`, `Bbands`, `Atr` | Technical indicators — $0.002/call |
| `Twelvedata.Metals` | `Price`, `TimeSeries` | Precious metals prices — $0.001/call |
| `Twelvedata.Indices` | `List`, `Quote` | Global index data — $0.001/call |
| `Twelvedata` | `GetTimeSeries`, `Price`, `Quote`, `Eod`, `ExchangeRate`, `ForexPairs`, `EconomicCalendar` | Direct endpoints — $0.001/call |

```go
// Time series (use GetTimeSeries, TimeSeries is a sub-resource with only Complex/POST)
ts, _ := client.Twelvedata.GetTimeSeries(ctx, &claw402.TwelvedataGetTimeSeriesParams{
    Symbol:     "EUR/USD",
    Interval:   "1h",
    Outputsize: "50",
})

// RSI  (sub-resource is Indicator, not TechnicalIndicators)
rsi, _ := client.Twelvedata.Indicator.Rsi(ctx, &claw402.TwelvedataIndicatorRsiParams{
    Symbol:     "AAPL",
    Interval:   "1day",
    TimePeriod: "14",
})

// Real-time price
price, _ := client.Twelvedata.Price(ctx, &claw402.TwelvedataPriceParams{
    Symbol: "BTC/USD",
})
```

### AI Providers

#### OpenAI

| Resource | Methods | Description |
|----------|---------|-------------|
| `Openai.Openai` | `Chat`, `ChatMini`, `Embeddings`, `EmbeddingsLarge`, `Images`, `Models` | OpenAI API — $0.001–0.05/call |

```go
// Chat (POST endpoint — pass body as map)
resp, _ := client.Openai.Openai.Chat(ctx, map[string]interface{}{
    "model": "gpt-4o",
    "messages": []map[string]string{
        {"role": "user", "content": "Analyze AAPL stock trend"},
    },
})
fmt.Println(string(resp))
```

#### Anthropic

| Resource | Methods | Description |
|----------|---------|-------------|
| `Anthropic.Anthropic` | `Messages`, `MessagesExtended`, `CountTokens` | Anthropic Claude API — $0.01–0.015/call |

```go
// Claude messages (POST endpoint)
resp, _ := client.Anthropic.Anthropic.Messages(ctx, map[string]interface{}{
    "model":      "claude-opus-4-6",
    "max_tokens": 1024,
    "messages": []map[string]string{
        {"role": "user", "content": "Summarize this earnings report: ..."},
    },
})
fmt.Println(string(resp))
```

#### DeepSeek

| Resource | Methods | Description |
|----------|---------|-------------|
| `Deepseek.Deepseek` | `Chat`, `ChatReasoner`, `Completions`, `Models` | DeepSeek chat, reasoning, beta completions, model listing — $0.001–0.005/call |

```go
resp, _ := client.Deepseek.Deepseek.Chat(ctx, map[string]interface{}{
    "messages": []map[string]string{
        {"role": "user", "content": "Explain BTC basis trade"},
    },
})
fmt.Println(string(resp))
```

#### Qwen

| Resource | Methods | Description |
|----------|---------|-------------|
| `Qwen.Qwen` | `ChatMax`, `ChatPlus`, `ChatTurbo`, `ChatFlash`, `ChatCoder`, `ChatVl` | Qwen chat, coder, and vision models — $0.002–0.01/call |

```go
resp, _ := client.Qwen.Qwen.ChatMax(ctx, map[string]interface{}{
    "messages": []map[string]string{
        {"role": "user", "content": "Write a Go HTTP middleware"},
    },
})
fmt.Println(string(resp))
```

## Configuration

```go
// Custom base URL
client, err := claw402.New("0xKey", claw402.WithBaseURL("https://custom.gateway"))
```

## How Payment Works

1. SDK sends a GET/POST request to the endpoint
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
