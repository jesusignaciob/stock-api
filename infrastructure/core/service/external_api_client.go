package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"stock-api/infrastructure/core/domain"
)

type ExternalAPIClient struct {
	baseURL string
	client  *http.Client
}

func NewExternalAPIClient(baseURL string) *ExternalAPIClient {
	return &ExternalAPIClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

type StockAPIResponse struct {
	Items    []*domain.Stock `json:"items"`
	NextPage string          `json:"next_page"`
}

func (c *ExternalAPIClient) FetchStocks(ctx context.Context, jwtToken, lastTicker string) ([]*domain.Stock, string, error) {
	url := c.baseURL
	if lastTicker != "" {
		url += fmt.Sprintf("?next_page=%s", lastTicker)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
	if err != nil {
		return nil, "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+jwtToken)
	req.Header.Add("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("API request failed: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	var apiResponse StockAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, "", fmt.Errorf("error decoding response: %w", err)
	}

	return apiResponse.Items, apiResponse.NextPage, nil
}
