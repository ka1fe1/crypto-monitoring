package bgeometrics

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ka1fe1/crypto-monitoring/pkg/logger"
)

const defaultTimeout = 10 * time.Second

type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

func NewClient(baseURL, apiKey string, opts ...Option) *Client {
	if baseURL == "" {
		baseURL = "https://bitcoin-data.com"
	}
	c := &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Option 是 Client 的可选配置函数
type Option func(*Client)

// WithTimeout 设置 HTTP 客户端超时时间
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.HTTPClient.Timeout = d
	}
}

func (c *Client) GetBalancedPrice() (float64, error) {
	// P2-3: API Key 为空时提前返回错误，避免发送无效请求
	if c.APIKey == "" {
		return 0, fmt.Errorf("bgeometrics API key is not configured")
	}

	url := fmt.Sprintf("%s/v1/balanced-price/1", c.BaseURL)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// BGeometrics /v1/balanced-price/1 返回单个对象 {"d":"...","unixTs":...,"balancedPrice":...}
	// 保留数组解析作为兼容性 fallback，因部分 endpoint（如无 {last} 参数时）可能返回数组。
	var singleRes MetricData
	if err := json.Unmarshal(body, &singleRes); err == nil && singleRes.BalancedPrice > 0 {
		return singleRes.BalancedPrice, nil
	}

	var arrRes []MetricData
	if err := json.Unmarshal(body, &arrRes); err == nil && len(arrRes) > 0 {
		return arrRes[len(arrRes)-1].BalancedPrice, nil
	}

	logger.Warn("Failed to decode response as single object or array: %s", string(body))
	return 0, fmt.Errorf("could not parse valid balanced price from response")
}
