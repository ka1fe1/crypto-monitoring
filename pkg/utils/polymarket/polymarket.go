package polymarket

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	BaseURL = "https://gamma-api.polymarket.com"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) SetHttpClient(client *http.Client) {
	c.httpClient = client
}

func (c *Client) GetMarketDetail(marketID string) (*MarketDetail, error) {
	body, err := c.fetchMarketRaw(marketID)
	if err != nil {
		return nil, err
	}

	var market Market
	if err := json.Unmarshal(body, &market); err != nil {
		return nil, fmt.Errorf("failed to unmarshal market data: %w", err)
	}

	return c.refineMarketData(&market), nil
}

func (c *Client) fetchMarketRaw(marketID string) ([]byte, error) {
	url := fmt.Sprintf("%s/markets/%s", BaseURL, marketID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("accept", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	return io.ReadAll(resp.Body)
}

func (c *Client) refineMarketData(market *Market) *MarketDetail {
	detail := &MarketDetail{
		Question:           market.Question,
		Slug:               market.Slug,
		Closed:             market.Closed,
		OneHourPriceChange: market.OneHourPriceChange,
		OneWeekPriceChange: market.OneWeekPriceChange,
	}

	// Parse Volume
	if market.Volume != "" {
		if vol, err := strconv.ParseFloat(market.Volume, 64); err == nil {
			detail.Volume = vol
		}
	}

	// Parse Outcomes and Prices
	detail.OutcomePrices = c.parseOutcomePrices(market.Outcomes, market.OutcomePrices)

	return detail
}

func (c *Client) parseOutcomePrices(outcomesStr, pricesStr string) map[string]float64 {
	var names []string
	var prices []string

	if err := json.Unmarshal([]byte(outcomesStr), &names); err != nil {
		return nil
	}
	if err := json.Unmarshal([]byte(pricesStr), &prices); err != nil {
		return nil
	}

	result := make(map[string]float64)
	for i, name := range names {
		if i < len(prices) {
			if p, err := strconv.ParseFloat(prices[i], 64); err == nil {
				result[name] = p
			}
		}
	}
	return result
}
