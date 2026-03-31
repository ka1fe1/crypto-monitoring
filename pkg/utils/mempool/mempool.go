package mempool

import (
	"fmt"
	"io"
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

// GetTipHeight returns the current block height of Bitcoin blockchain
func (c *Client) GetTipHeight() (int64, error) {
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return 0, fmt.Errorf("invalid base url: %w", err)
	}
	u.Path, err = url.JoinPath(u.Path, "/api/blocks/tip/height")
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return 0, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("http do err: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("mempool api error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	height, err := strconv.ParseInt(string(body), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse height %s: %w", string(body), err)
	}

	return height, nil
}
