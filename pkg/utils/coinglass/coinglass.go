package coinglass

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	BaseURL = "https://open-api-v4.coinglass.com"
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

type Response[T any] struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

type AHR999Data struct {
	DateString   string  `json:"date_string"`
	AveragePrice float64 `json:"average_price"`
	AHR999Value  float64 `json:"ahr999_value"`
	CurrentValue float64 `json:"current_value"`
}

type FearGreedData struct {
	Values   []float64 `json:"values"`
	Price    []float64 `json:"price"`
	TimeList []int64   `json:"time_list"`
}

func (c *Client) request(method, path string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", BaseURL, path)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("CG-API-KEY", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// GetAHR999 fetches Bitcoin AHR999 index data.
// Endpoint: /api/index/ahr999
func (c *Client) GetAHR999() ([]AHR999Data, error) {
	body, err := c.request("GET", "/api/index/ahr999")
	if err != nil {
		return nil, err
	}

	var res Response[[]AHR999Data]
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	if res.Code != "0" {
		return nil, fmt.Errorf("api error: %s", res.Msg)
	}

	return res.Data, nil
}

// GetFearGreedIndex fetches Crypto Fear & Greed Index data.
// Endpoint: /api/index/fear-greed-history
func (c *Client) GetFearGreedIndex() ([]FearGreedData, error) {
	body, err := c.request("GET", "/api/index/fear-greed-history")
	if err != nil {
		return nil, err
	}

	var res Response[[]FearGreedData]
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	if res.Code != "0" {
		return nil, fmt.Errorf("api error: %s", res.Msg)
	}

	return res.Data, nil
}
