package api_v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

const (
	apiV1ExchangeList = "/api/v1/exchange/list"
)

// GetExchangeListOutput is returned by the function GetExchangeList.
// It represents the data of an exchange.
type GetExchangeListOutput struct {
	Name         string
	ID           uuid.UUID
	TickerFormat string
	URL          string
}

// GetExchangeList is used when you want to get the list of exchanges and their data,
// the output will be:
//
//	[]GetExchangeListOutput{
//		{
//			Name: "MEXC",
//			ID: "12345678-1234-1234-1234-123456789012",
//			TickerFormat: "%s_%s",
//			Url: "https://mexc.com/",
//		}
//	}
func (c *Client) GetExchangeList(ctx context.Context) ([]GetExchangeListOutput, error) {
	req, err := http.NewRequest("GET", c.baseUrl+apiV1ExchangeList, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
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

	var output []GetExchangeListOutput
	if err := json.Unmarshal(bodyBytes, &output); err != nil {
		return nil, err
	}

	return output, nil
}
