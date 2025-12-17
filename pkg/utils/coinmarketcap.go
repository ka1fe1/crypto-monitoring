package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const DefaultBaseURL = "https://pro-api.coinmarketcap.com/v2"

type CoinMarketClient struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
}

func NewCoinMarketClient(apiKey string) *CoinMarketClient {
	return &CoinMarketClient{
		APIKey:  apiKey,
		BaseURL: DefaultBaseURL,
		Client:  &http.Client{},
	}
}

type QuoteResponse struct {
	Status Status            `json:"status"`
	Data   map[string]Crypto `json:"data"`
}

type Status struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

type Crypto struct {
	ID     int              `json:"id"`
	Name   string           `json:"name"`
	Symbol string           `json:"symbol"`
	Quote  map[string]Quote `json:"quote"`
}

type Quote struct {
	Price            float64 `json:"price"`
	PercentChange1h  float64 `json:"percent_change_1h"`
	PercentChange24h float64 `json:"percent_change_24h"`
	LastUpdated      string  `json:"last_updated"`
}

type TokenInfo struct {
	Price            float64
	Symbol           string
	PercentChange1h  float64
	PercentChange24h float64
	LastUpdated      time.Time
}

// GetPrice fetches the prices of multiple cryptocurrency IDs in USD.
// It uses the /v2/cryptocurrency/quotes/latest endpoint.
func (c *CoinMarketClient) GetPrice(ids []string) (map[string]TokenInfo, error) {
	u, err := url.Parse(c.BaseURL + "/cryptocurrency/quotes/latest")
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	q := u.Query()
	q.Set("id", strings.Join(ids, ","))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-CMC_PRO_API_KEY", c.APIKey)
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

	var quoteResp QuoteResponse
	if err := json.Unmarshal(body, &quoteResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if quoteResp.Status.ErrorCode != 0 {
		return nil, fmt.Errorf("API error: %s (code: %d)", quoteResp.Status.ErrorMessage, quoteResp.Status.ErrorCode)
	}

	result := make(map[string]TokenInfo)
	for id, crypto := range quoteResp.Data {
		usdQuote := crypto.Quote["USD"]
		lastUpdated, _ := time.Parse(time.RFC3339, usdQuote.LastUpdated)
		result[id] = TokenInfo{
			Price:            usdQuote.Price,
			Symbol:           crypto.Symbol,
			PercentChange1h:  usdQuote.PercentChange1h,
			PercentChange24h: usdQuote.PercentChange24h,
			LastUpdated:      lastUpdated,
		}
	}

	return result, nil
}

// GetPriceBySymbol fetches the prices of multiple cryptocurrency symbols in USD.
// It uses the /v2/cryptocurrency/quotes/latest endpoint.
func (c *CoinMarketClient) GetPriceBySymbol(symbols []string) (map[string]TokenInfo, error) {
	u, err := url.Parse(c.BaseURL + "/cryptocurrency/quotes/latest")
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	q := u.Query()
	q.Set("symbol", strings.Join(symbols, ","))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-CMC_PRO_API_KEY", c.APIKey)
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

	// CMC V2 API response for symbol query puts the data in an array for each symbol because symbols can be duplicated?
	// Actually, for quotes/latest v2:
	// "Each map key will be a symbol. The value will be an array of currency objects..."
	// Wait, the existing QuoteResponse has Data map[string]Crypto. Crypto matches a singe object?
	// Let's verify standard response format for symbol query.
	// Documentation says: "If one of "slug" or "symbol" follows: response.data.<symbol/slug> is an Array of objects"
	// So my struct for Data might be wrong for symbol query if it expects a single object.
	// Existing GetPrice uses ID, where response is map[string]Crypto (object).
	// For symbol, it is map[string][]Crypto.

	// I need to define a different response struct or handle dynamic unmarshalling.
	// Since I'm in Go, defining a new struct is safer.

	type QuoteResponseBySymbol struct {
		Status Status              `json:"status"`
		Data   map[string][]Crypto `json:"data"`
	}

	var quoteResp QuoteResponseBySymbol
	if err := json.Unmarshal(body, &quoteResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if quoteResp.Status.ErrorCode != 0 {
		return nil, fmt.Errorf("API error: %s (code: %d)", quoteResp.Status.ErrorMessage, quoteResp.Status.ErrorCode)
	}

	result := make(map[string]TokenInfo)
	for symbol, cryptoList := range quoteResp.Data {
		if len(cryptoList) > 0 {
			// Just take the first one if multiple exist for the same symbol (e.g. different networks/tokens)
			// Ideally we'd filter, but for major coins like ETH/SOL it's usually the first.
			crypto := cryptoList[0]
			usdQuote := crypto.Quote["USD"]
			lastUpdated, _ := time.Parse(time.RFC3339, usdQuote.LastUpdated)
			result[symbol] = TokenInfo{
				Price:            usdQuote.Price,
				Symbol:           crypto.Symbol,
				PercentChange1h:  usdQuote.PercentChange1h,
				PercentChange24h: usdQuote.PercentChange24h,
				LastUpdated:      lastUpdated,
			}
		}
	}

	return result, nil
}

type DexQuoteResponse struct {
	Status DexStatus `json:"status"`
	Data   []DexPair `json:"data"`
}

type DexStatus struct {
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

type DexPair struct {
	ContractAddress           string         `json:"contract_address"`
	Name                      string         `json:"name"`
	BaseAssetSymbol           string         `json:"base_asset_symbol"`
	BaseAssetContractAddress  string         `json:"base_asset_contract_address"`
	QuoteAssetSymbol          string         `json:"quote_asset_symbol"`
	QuoteAssetContractAddress string         `json:"quote_asset_contract_address"`
	DexSlug                   string         `json:"dex_slug"`
	NetworkSlug               string         `json:"network_slug"`
	Quote                     []DexPairQuote `json:"quote"`
}

type DexPairQuote struct {
	Price            float64 `json:"price"`
	PriceByQuote     float64 `json:"price_by_quote"`
	Volume24h        float64 `json:"volume_24h"`
	PercentChange1h  float64 `json:"percent_change_1h"`
	PercentChange24h float64 `json:"percent_change_24h"`
	Liquidity        float64 `json:"liquidity"`
	LastUpdated      string  `json:"last_updated"`
}

// GetDexPairQuotes fetches the quotes for DEX pairs using the /v4/dex/pairs/quotes/latest endpoint.
func (c *CoinMarketClient) GetDexPairQuotes(contractAddresses []string, networkSlug, networkId string) (map[string]DexPair, error) {
	// Note: The BaseURL constant points to v2, but this endpoint is v4.
	// We need to construct the URL carefully.
	// BaseURL is "https://pro-api.coinmarketcap.com/v2"
	// We want "https://pro-api.coinmarketcap.com/v4/dex/pairs/quotes/latest"

	// Let's just use the full URL for v4 to be safe and clear, or replace v2 with v4.
	v4URL := strings.Replace(c.BaseURL, "/v2", "/v4", 1) + "/dex/pairs/quotes/latest"

	u, err := url.Parse(v4URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	q := u.Query()
	q.Set("contract_address", strings.Join(contractAddresses, ","))
	if networkId != "" {
		q.Set("network_id", networkId)
	} else if networkSlug != "" {
		q.Set("network_slug", networkSlug)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-CMC_PRO_API_KEY", c.APIKey)
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

	var dexResp DexQuoteResponse
	if err := json.Unmarshal(body, &dexResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if dexResp.Status.ErrorCode != "0" {
		return nil, fmt.Errorf("API error: %s (code: %s)", dexResp.Status.ErrorMessage, dexResp.Status.ErrorCode)
	}

	pairs := make(map[string]DexPair)
	for _, pair := range dexResp.Data {
		pairs[pair.ContractAddress] = pair
	}

	return pairs, nil
}
