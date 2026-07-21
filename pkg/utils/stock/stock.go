package stock

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type StockClient struct {
	client *http.Client
}

func NewStockClient() *StockClient {
	return &StockClient{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// GetStockPrices 批量获取股票价格信息（支持港股 hkXXXXX 与 A股 shXXXXXX / szXXXXXX）
func (s *StockClient) GetStockPrices(ctx context.Context, symbols []string) (map[string]StockInfo, error) {
	if len(symbols) == 0 {
		return make(map[string]StockInfo), nil
	}

	url := fmt.Sprintf("http://qt.gtimg.cn/q=%s", strings.Join(symbols, ","))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create stock request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stock prices: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected http status: %d", resp.StatusCode)
	}

	// 腾讯财经返回 GBK 编码，转为 UTF-8
	utf8Reader := transform.NewReader(resp.Body, simplifiedchinese.GBK.NewDecoder())
	bodyBytes, err := io.ReadAll(utf8Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return s.parseResponse(string(bodyBytes), symbols)
}

func (s *StockClient) parseResponse(body string, symbols []string) (map[string]StockInfo, error) {
	results := make(map[string]StockInfo)
	lines := strings.Split(body, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 返回格式形如: v_sh600519="51~贵州茅台~600519~1600.00~...~";
		parts := strings.Split(line, "=")
		if len(parts) < 2 {
			continue
		}

		left := parts[0]
		right := parts[1]

		// 提取 symbol code, 如 v_sh600519 -> sh600519, v_r_hk00700 -> hk00700 等
		code := strings.TrimPrefix(left, "v_")
		code = strings.TrimPrefix(code, "r_") // 某些港股接口会有 r_ 前缀
		code = strings.TrimSpace(code)

		right = strings.Trim(right, "\";")
		right = strings.TrimSpace(right)
		if right == "" {
			continue
		}

		dataFields := strings.Split(right, "~")
		// A 股与港股分割出来的字段略有差异，但名称在 index 1，价格在 index 3，涨跌幅在 index 32
		if len(dataFields) < 33 {
			continue
		}

		name := dataFields[1]
		price, _ := strconv.ParseFloat(dataFields[3], 64)
		percentChange, _ := strconv.ParseFloat(dataFields[32], 64)

		results[code] = StockInfo{
			Code:          code,
			Name:          name,
			Price:         price,
			PercentChange: percentChange,
			LastUpdated:   time.Now(),
		}
	}

	return results, nil
}
