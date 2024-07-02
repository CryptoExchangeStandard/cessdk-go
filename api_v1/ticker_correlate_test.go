package api_v1_test

import (
	"context"
	"fmt"
	"testing"

	api "github.com/CryptoExchangeStandard/cessdk-go/api_v1"
)

func TestPostTickerCorrelate(t *testing.T) {
	client := api.NewClient("your_key", "http://159.223.21.91:8080", nil)

	outputTickers, err := client.PostTickerCorrelate(context.Background(), api.PostTickerCorrelateInput{
		ExchangeFrom: "Binance",
		ExchangeTo:   []string{"MEXC", "Gate.io", "Bitrue", "xt.com"},
		ExchangeTicker: api.Ticker{
			Base:  "QI",
			Quote: "USDT",
		},
	})

	if err != nil {
		t.Fatalf("PostTickerCorrelate failed: %v", err)
	}

	fmt.Println("PostTickerCorrelate:\n", outputTickers)
}
