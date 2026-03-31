package alternative

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetFng fetches the Fear and Greed index from alternative.me
// limit=1 returns only the latest day.
func (c *Client) GetFng(limit int) (*FngResponse, error) {
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base url: %w", err)
	}
	u.Path, err = url.JoinPath(u.Path, "/fng/")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http do err: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("alternative api error: status=%d", resp.StatusCode)
	}

	var fng FngResponse
	if err := json.NewDecoder(resp.Body).Decode(&fng); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if fng.Metadata.Error != nil {
		return nil, fmt.Errorf("api returned error: %s", *fng.Metadata.Error)
	}

	return &fng, nil
}
