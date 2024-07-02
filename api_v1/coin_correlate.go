package api_v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

const (
	apiV1CoinCorrelate = "/api/v1/coin/correlate"
)

var (
	ErrPostCoinCorrelateInputExchangeFromNotValid = errors.New("input not valid: one of ExchangeFrom or ExchangeFromID should be provided")
	ErrPostCoinCorrelateInputExchangeToNotValid   = errors.New("input not valid: one of ExchangeTo or ExchangeToID should be provided")
	ErrPostCoinCorrelateInputExchangeCoinNotValid = errors.New("input not valid: one of ExchangeCoin, ExchangeCoinBase or ExchangeCoinID should be provided")
)

// PostCoinCorrelateInput should be filled with at least the name or the ID of all.
// For example, you can fill ExchangeFrom, ExchangeToID, CoinID, or ExchangeFromID, ExchangeToID, Coin, etc.
// We recommend using the *ID values.
type PostCoinCorrelateInput struct {
	ExchangeFrom   string    `json:"ExchangeFrom,omitempty"`
	ExchangeFromID uuid.UUID `json:"ExchangeFromID,omitempty"`

	ExchangeTo   []string    `json:"ExchangeTo,omitempty"`
	ExchangeToID []uuid.UUID `json:"ExchangeToID,omitempty"`

	ExchangeCoin     string    `json:"ExchangeCoin,omitempty"`
	ExchangeCoinBase string    `json:"ExchangeCoinBase,omitempty"`
	ExchangeCoinID   uuid.UUID `json:"ExchangeCoinID,omitempty"`
}

// Note: It is being checked serverside too.
func validateCoinCorrelateInput(input PostCoinCorrelateInput) error {
	if len(input.ExchangeFrom) == 0 && input.ExchangeFromID == uuid.Nil {
		return ErrPostCoinCorrelateInputExchangeFromNotValid
	}
	if len(input.ExchangeTo) == 0 && len(input.ExchangeToID) == 0 {
		return ErrPostCoinCorrelateInputExchangeToNotValid
	}
	if len(input.ExchangeCoin) == 0 && len(input.ExchangeCoinBase) == 0 && input.ExchangeCoinID == uuid.Nil {
		return ErrPostCoinCorrelateInputExchangeCoinNotValid
	}

	if len(input.ExchangeFrom) != 0 && input.ExchangeFromID != uuid.Nil {
		return ErrPostCoinCorrelateInputExchangeFromNotValid
	}
	if len(input.ExchangeTo) != 0 && len(input.ExchangeToID) != 0 {
		return ErrPostCoinCorrelateInputExchangeToNotValid
	}

	filled := 0
	if len(input.ExchangeCoin) != 0 {
		filled++
	}
	if len(input.ExchangeCoinBase) != 0 {
		filled++
	}
	if input.ExchangeCoinID != uuid.Nil {
		filled++
	}

	if filled != 1 {
		return ErrPostCoinCorrelateInputExchangeCoinNotValid
	}

	return nil
}

// PostCoinCorrelateOutput is returned by the function PostCoinCorrelate.
// ExchangeName, ExchangeID represents the name and ID of the exchange
// ExchangeCoin, ExchangeCoinBase represents the name and base of the coin of the exchange - not the standard
type PostCoinCorrelateOutput struct {
	ExchangeName              string
	ExchangeID                uuid.UUID
	ExchangeCoinID            uuid.UUID
	ExchangeCoin              string
	ExchangeCoinBase          string
	ExchangeCoinUnsafetyScore int
}

// PostCoinCorrelate is used when you want to know the name and the base of a coin of another exchange,
// based on the value of the same coin of another exchange.
// Ex: You want to know the equivalent of name=BENQI base=QI from Binance, on MEXC and Gate, input could be:
//
//	CoinCorrelateInput{
//		ExchangeFrom: "Binance",
//		ExchangeTo: [
//			"MEXC",
//			"Gate"
//		],
//		ExchangeCoin: "BENQI"
//	}
//
// and the output will be:
//
//	[]PostCoinCorrelateOutput{
//		{
//			ExchangeName: "MEXC",
//			ExchangeID: "12345678-1234-1234-1234-123456789012",
//			ExchangeCoin: "BENQI",
//			ExchangeCoinBase: "BENQI",
//			ExchangeCoinUnsafetyScore: 0
//		},
//		{
//			ExchangeName: "Gate",
//			ExchangeID: "abcdef12-1234-1234-1234-123456789012",
//			ExchangeCoin: "BENQI",
//			ExchangeCoinBase: "BENQI",
//			ExchangeCoinUnsafetyScore: 0
//		}
//	}
func (c *Client) PostCoinCorrelate(ctx context.Context, input PostCoinCorrelateInput) ([]PostCoinCorrelateOutput, error) {
	if err := validateCoinCorrelateInput(input); err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.baseUrl+apiV1CoinCorrelate, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("CES_API_KEY", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: Received status code: %v. Body: %v", resp.StatusCode, string(bodyBytes))
	}

	var output []PostCoinCorrelateOutput
	if err := json.Unmarshal(bodyBytes, &output); err != nil {
		return nil, err
	}

	return output, nil
}
