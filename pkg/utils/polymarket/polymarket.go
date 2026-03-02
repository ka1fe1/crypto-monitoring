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

// ResolveProxyWallet converts an EOA wallet address to a Polymarket proxyWallet address
// by querying the gamma-api public-profile endpoint.
// The data-api endpoints (leaderboard, value, positions) require proxyWallet addresses.
func (c *Client) ResolveProxyWallet(address string) (string, error) {
	url := fmt.Sprintf("%s/public-profile?address=%s", BaseURL, address)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var profile PublicProfileResponse
	if err := json.Unmarshal(body, &profile); err != nil {
		return "", fmt.Errorf("failed to unmarshal profile: %w", err)
	}

	if profile.ProxyWallet == "" {
		return "", fmt.Errorf("no proxyWallet found for address: %s", address)
	}

	return profile.ProxyWallet, nil
}

// GetTraderLeaderboardRankings fetches the ranking data for a single trader address.
func (c *Client) GetTraderLeaderboardRankings(address string) (*LeaderboardResponse, error) {
	// API: https://data-api.polymarket.com/v1/leaderboard?user={address}&timePeriod=all&orderBy=vol
	url := fmt.Sprintf("https://data-api.polymarket.com/v1/leaderboard?user=%s&timePeriod=ALL&orderBy=VOL", address)
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Some APIs return lists, assuming it returns list for a user query or a single object.
	var result []LeaderboardResponse
	if err := json.Unmarshal(body, &result); err != nil {
		// fallback to object parsing if it's not array
		var single LeaderboardResponse
		if err2 := json.Unmarshal(body, &single); err2 != nil {
			return nil, fmt.Errorf("failed to unmarshal leaderboard: %v / %v", err, err2)
		}
		return &single, nil
	}

	if len(result) > 0 {
		return &result[0], nil
	}
	return &LeaderboardResponse{Rank: "0", ProxyWallet: address, Vol: 0, Pnl: 0}, nil
}

// GetTotalValueOfUserPositions fetches the total USD value of positions for a given address.
func (c *Client) GetTotalValueOfUserPositions(address string) (*TotalValueResponse, error) {
	// API: https://data-api.polymarket.com/value?user={address}
	url := fmt.Sprintf("https://data-api.polymarket.com/value?user=%s", address)
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// This endpoint returns an array like [{ "user": "0x...", "value": 123 }]
	var result []TotalValueResponse
	if err := json.Unmarshal(body, &result); err != nil {
		// Fallback for single object just in case
		var single TotalValueResponse
		if err2 := json.Unmarshal(body, &single); err2 != nil {
			return nil, fmt.Errorf("failed to unmarshal total value: %w", err)
		}
		return &single, nil
	}

	if len(result) > 0 {
		return &result[0], nil
	}
	// Return default if empty
	return &TotalValueResponse{Value: 0}, nil
}

// GetCurrentPositionsForUser fetches the list of active positions for a user.
func (c *Client) GetCurrentPositionsForUser(address string) (*CurrentPositionsResponse, error) {
	// API: https://data-api.polymarket.com/positions?user={address}
	url := fmt.Sprintf("https://data-api.polymarket.com/positions?user=%s", address)
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result CurrentPositionsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal positions: %w", err)
	}

	return &result, nil
}
