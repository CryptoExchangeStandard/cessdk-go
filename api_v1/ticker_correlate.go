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
	apiV1TickerCorrelate = "/api/v1/ticker/correlate"
)

var (
	ErrPostTickerCorrelateInputExchangeFromNotValid   = errors.New("input not valid: one of ExchangeFrom or ExchangeFromID should be provided")
	ErrPostTickerCorrelateInputExchangeToNotValid     = errors.New("input not valid: one of ExchangeTo or ExchangeToID should be provided")
	ErrPostTickerCorrelateInputExchangeTickerNotValid = errors.New("input not valid: one of ExchangeTicker or ExchangeTickerID should be provided")
)

type Ticker struct {
	Base  string `json:"Base"`
	Quote string `json:"Quote"`
}

type PostTickerCorrelateInput struct {
	ExchangeFrom   string    `json:"ExchangeFrom,omitempty"`
	ExchangeFromID uuid.UUID `json:"ExchangeFromID,omitempty"`

	ExchangeTo   []string    `json:"ExchangeTo,omitempty"`
	ExchangeToID []uuid.UUID `json:"ExchangeToID,omitempty"`

	ExchangeTicker   Ticker    `json:"ExchangeTicker,omitempty"`
	ExchangeTickerID uuid.UUID `json:"ExchangeTickerID,omitempty"`
}

// Note: It is being checked serverside too.
func validatePostTickerCorrelateInput(input PostTickerCorrelateInput) error {
	if len(input.ExchangeFrom) == 0 && input.ExchangeFromID == uuid.Nil {
		return ErrPostTickerCorrelateInputExchangeFromNotValid
	}
	if len(input.ExchangeTo) == 0 && len(input.ExchangeToID) == 0 {
		return ErrPostTickerCorrelateInputExchangeToNotValid
	}
	if len(input.ExchangeTicker.Base) == 0 && len(input.ExchangeTicker.Quote) == 0 && input.ExchangeTickerID == uuid.Nil {
		return ErrPostTickerCorrelateInputExchangeTickerNotValid
	}

	if len(input.ExchangeFrom) != 0 && input.ExchangeFromID != uuid.Nil {
		return ErrPostTickerCorrelateInputExchangeFromNotValid
	}
	if len(input.ExchangeTo) != 0 && len(input.ExchangeToID) != 0 {
		return ErrPostTickerCorrelateInputExchangeToNotValid
	}

	filled := 0
	if len(input.ExchangeTicker.Base) != 0 || len(input.ExchangeTicker.Quote) != 0 {
		filled++
	}
	if input.ExchangeTickerID != uuid.Nil {
		filled++
	}

	if filled != 1 {
		return ErrPostTickerCorrelateInputExchangeTickerNotValid
	}

	return nil
}

type PostTickerCorrelateOutput struct {
	ExchangeName     string
	ExchangeID       uuid.UUID
	ExchangeTickerID uuid.UUID
	ExchangeTicker   Ticker
}

// PostTickerCorrelate is used when you want to know the base and the quote of a ticker of another exchange,
// based on the value of the same ticker of another exchange.
// Ex: You want to know the equivalent of BTC_USDT from Binance, on MEXC and Gate, input could be:
//
//	PostTickerCorrelateInput{
//		ExchangeFrom: "Binance",
//		ExchangeTo: [
//			"MEXC",
//			"Gate"
//		],
//		ExchangeTicker: Ticker{
//			Base:"BTC",
//			Quote:"USDT"
//		}
//	}
//
// and the output will be:
//
//	[]PostTickerCorrelateOutput{
//		{
//			ExchangeName: "MEXC",
//			ExchangeID: "12345678-1234-1234-1234-123456789012",
//			ExchangeTickerID: "11111111-1234-1234-1234-123456789012",
//			Ticker: {
//				Base: "BTC",
//				Quote: "USDT"
//			}
//		},
//		{
//			ExchangeName: "Gate",
//			ExchangeID: "abcdef12-1234-1234-1234-123456789012",
//			ExchangeTickerID: "99999999-1234-1234-1234-123456789012",
//			Ticker: {
//				Base: "BTC",
//				Quote: "USDT"
//			}
//		}
//	}
func (c *Client) PostTickerCorrelate(ctx context.Context, input PostTickerCorrelateInput) ([]PostTickerCorrelateOutput, error) {
	if err := validatePostTickerCorrelateInput(input); err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.baseUrl+apiV1TickerCorrelate, bytes.NewBuffer(jsonData))
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

	var output []PostTickerCorrelateOutput
	if err := json.Unmarshal(bodyBytes, &output); err != nil {
		return nil, err
	}

	return output, nil
}
