// Example: basic usage of the claw402 Go SDK.
//
// Usage:
//
//	WALLET_PRIVATE_KEY=0x... go run .
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	claw402 "github.com/NoFxAiOS/claw402-go"
)

func main() {
	key := os.Getenv("WALLET_PRIVATE_KEY")
	if key == "" {
		log.Fatal("set WALLET_PRIVATE_KEY env var")
	}

	client, err := claw402.New(key)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 1. Fear & Greed Index (no params)
	fmt.Println("=== Fear & Greed Index ===")
	data, err := client.Coinank.Indicator.FearGreed(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))

	// 2. Fund flow with params
	fmt.Println("\n=== Fund Flow (SWAP, top 5) ===")
	data, err = client.Coinank.Fund.Realtime(ctx, &claw402.CoinankFundRealtimeParams{
		ProductType: "SWAP",
		Size:        "5",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))

	// 3. AI500 signals
	fmt.Println("\n=== AI500 Top Signals ===")
	data, err = client.Nofxos.Ai500.List(ctx, &claw402.NofxosAi500ListParams{
		Limit: "10",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))

	// 4. Net capital inflow ranking
	fmt.Println("\n=== Net Inflow Top 10 (1h) ===")
	data, err = client.Nofxos.Netflow.TopRanking(ctx, &claw402.NofxosNetflowTopRankingParams{
		Limit:    "10",
		Duration: "1h",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))
}
