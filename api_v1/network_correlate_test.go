package api_v1_test

import (
	"context"
	"fmt"
	"testing"

	api "github.com/CryptoExchangeStandard/cessdk-go/api_v1"
)

func TestPostNetworkCorrelate(t *testing.T) {
	client := api.NewClient("your_key", "http://159.223.21.91:8080", nil)

	outputNetworks, err := client.PostNetworkCorrelate(context.Background(), api.PostNetworkCorrelateInput{
		ExchangeFrom:        "Binance",
		ExchangeTo:          []string{"MEXC", "Gate.io", "Bitrue", "xt.com"},
		ExchangeNetworkCode: "ETH",
	})

	if err != nil {
		t.Fatalf("PostNetworkCorrelate failed: %v", err)
	}

	fmt.Println("PostNetworkCorrelate:\n", outputNetworks)
}
