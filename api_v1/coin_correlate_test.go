package api_v1_test

import (
	"context"
	"fmt"
	"testing"

	api "github.com/CryptoExchangeStandard/cessdk-go/api_v1"
)

func TestPostCoinCorrelate(t *testing.T) {
	client := api.NewClient("your_key", "http://159.223.21.91:8080", nil)

	outputCoins, err := client.PostCoinCorrelate(context.Background(), api.PostCoinCorrelateInput{
		ExchangeFrom: "Binance",
		ExchangeTo:   []string{"MEXC", "Gate.io", "Bitrue", "xt.com"},
		ExchangeCoin: "BENQI",
	})

	if err != nil {
		t.Fatalf("PostCoinCorrelate failed: %v", err)
	}

	fmt.Println("PostCoinCorrelate:\n", outputCoins)
}
