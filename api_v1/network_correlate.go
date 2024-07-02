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
	apiV1NetworkCorrelate = "/api/v1/network/correlate"
)

var (
	ErrPostNetworkCorrelateInputExchangeFromNotValid    = errors.New("input not valid: one of ExchangeFrom or ExchangeFromID should be provided")
	ErrPostNetworkCorrelateInputExchangeToNotValid      = errors.New("input not valid: one of ExchangeTo or ExchangeToID should be provided")
	ErrPostNetworkCorrelateInputExchangeNetworkNotValid = errors.New("input not valid: one of ExchangeNetwork, ExchangeNetworkCode or ExchangeNetworkID should be provided")
)

type PostNetworkCorrelateInput struct {
	ExchangeFrom   string    `json:"ExchangeFrom,omitempty"`
	ExchangeFromID uuid.UUID `json:"ExchangeFromID,omitempty"`

	ExchangeTo   []string    `json:"ExchangeTo,omitempty"`
	ExchangeToID []uuid.UUID `json:"ExchangeToID,omitempty"`

	ExchangeNetwork     string    `json:"ExchangeNetwork,omitempty"`
	ExchangeNetworkCode string    `json:"ExchangeNetworkCode,omitempty"`
	ExchangeNetworkID   uuid.UUID `json:"ExchangeNetworkID,omitempty"`
}

// Note: It is being checked serverside too.
func validatePostNetworkCorrelateInput(input PostNetworkCorrelateInput) error {
	if len(input.ExchangeFrom) == 0 && input.ExchangeFromID == uuid.Nil {
		return ErrPostNetworkCorrelateInputExchangeFromNotValid
	}
	if len(input.ExchangeTo) == 0 && len(input.ExchangeToID) == 0 {
		return ErrPostNetworkCorrelateInputExchangeToNotValid
	}
	if len(input.ExchangeNetwork) == 0 && len(input.ExchangeNetworkCode) == 0 && input.ExchangeNetworkID == uuid.Nil {
		return ErrPostNetworkCorrelateInputExchangeNetworkNotValid
	}

	if len(input.ExchangeFrom) != 0 && input.ExchangeFromID != uuid.Nil {
		return ErrPostNetworkCorrelateInputExchangeFromNotValid
	}
	if len(input.ExchangeTo) != 0 && len(input.ExchangeToID) != 0 {
		return ErrPostNetworkCorrelateInputExchangeToNotValid
	}

	filled := 0
	if len(input.ExchangeNetwork) != 0 {
		filled++
	}
	if len(input.ExchangeNetworkCode) != 0 {
		filled++
	}
	if input.ExchangeNetworkID != uuid.Nil {
		filled++
	}

	if filled != 1 {
		return ErrPostNetworkCorrelateInputExchangeNetworkNotValid
	}

	return nil
}

type PostNetworkCorrelateOutput struct {
	ExchangeName                 string
	ExchangeID                   uuid.UUID
	ExchangeNetworkID            uuid.UUID
	ExchangeNetwork              string
	ExchangeNetworkCode          string
	ExchangeNetworkUnsafetyScore int
}

// PostNetworkCorrelate is used when you want to know the name and the code of a network of another exchange,
// based on the value of the same network of another exchange.
// Ex: You want to know the equivalent of name=Ethereum(ERC20) code=ETH from Binance, on MEXC and Gate, input could be:
//
//	PostNetworkCorrelateInput{
//		ExchangeFrom: "Binance",
//		ExchangeTo: [
//			"MEXC",
//			"Gate"
//		],
//		ExchangeNetworkCode: "ETH"
//	}
//
// and the output will be:
//
//	[]PostNetworkCorrelateOutput{
//		{
//			ExchangeName: "MEXC",
//			ExchangeID: "12345678-1234-1234-1234-123456789012",
//			ExchangeNetworkID: "11111111-1234-1234-1234-123456789012",
//			ExchangeNetwork: "Ethereum",
//			ExchangeNetworkCode: "ETH"
//		},
//		{
//			ExchangeName: "Gate",
//			ExchangeID: "abcdef12-1234-1234-1234-123456789012",
//			ExchangeNetworkID: "99999999-1234-1234-1234-123456789012",
//			ExchangeNetwork: "",
//			ExchangeNetworkCode: "ETH"
//		}
//	}
func (c *Client) PostNetworkCorrelate(ctx context.Context, input PostNetworkCorrelateInput) ([]PostNetworkCorrelateOutput, error) {
	if err := validatePostNetworkCorrelateInput(input); err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.baseUrl+apiV1NetworkCorrelate, bytes.NewBuffer(jsonData))
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

	var output []PostNetworkCorrelateOutput
	if err := json.Unmarshal(bodyBytes, &output); err != nil {
		return nil, err
	}

	return output, nil
}
