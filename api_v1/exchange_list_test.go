package api_v1_test

import (
	"context"
	"fmt"
	"testing"

	api "github.com/CryptoExchangeStandard/cessdk-go/api_v1"
)

func TestGetExchangeList(t *testing.T) {
	client := api.NewClient("your_key", api.DEFAULT_BASE_API_ENDPOINT, nil)

	outputExchanges, err := client.GetExchangeList(context.Background())

	if err != nil {
		t.Fatalf("GetExchangeList failed: %v", err)
	}

	fmt.Println("GetExchangeList:\n", outputExchanges)
}
