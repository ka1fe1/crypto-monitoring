package opensea

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const OpenSeaBaseURL = "https://api.opensea.io/api/v2"

type OpenSeaClient struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
}

func NewOpenSeaClient(apiKey string) *OpenSeaClient {
	return &OpenSeaClient{
		APIKey:  apiKey,
		BaseURL: OpenSeaBaseURL,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type CollectionStatsResponse struct {
	Total CollectionStats `json:"total"`
}

type CollectionStats struct {
	FloorPrice       float64 `json:"floor_price"`
	FloorPriceSymbol string  `json:"floor_price_symbol"`
}

// GetCollectionStats fetches the stats for a collection including floor price.
// slug: The collection slug (e.g. "infinex-patrons")
func (c *OpenSeaClient) GetCollectionStats(slug string) (*CollectionStats, error) {
	// OpenSea API V2: /collections/{collection_slug}/stats
	u, err := url.Parse(fmt.Sprintf("%s/collections/%s/stats", c.BaseURL, slug))
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.APIKey != "" {
		req.Header.Set("X-API-KEY", c.APIKey)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var statsResp CollectionStatsResponse
	if err := json.Unmarshal(body, &statsResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &statsResp.Total, nil
}
