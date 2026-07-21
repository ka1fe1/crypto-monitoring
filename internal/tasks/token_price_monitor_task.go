package tasks

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ka1fe1/crypto-monitoring/internal/service"
	"github.com/ka1fe1/crypto-monitoring/pkg/logger"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/alter/dingding"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/constant"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/stock"
)

type TokenPriceMonitorTask struct {
	tokenService     service.TokenService
	stockClient      *stock.StockClient
	dingBot          *dingding.DingBot
	ticker           *time.Ticker
	stop             chan bool
	tokenIds         []string
	rwaTokenIds      []string
	rwaTokenNames    map[string]string
	hkStockIds       []string
	hkStockNames     map[string]string
	aStockIds        []string
	aStockNames      map[string]string
	interval         time.Duration
	quietHoursParams utils.QuietHoursParams
	lastRunTime      time.Time
}

func NewTokenPriceMonitorTask(
	tokenService service.TokenService,
	dingBot *dingding.DingBot,
	tokenIdsStr string,
	rwaTokenIds []string,
	rwaTokenNames map[string]string,
	hkStockIds []string,
	hkStockNames map[string]string,
	aStockIds []string,
	aStockNames map[string]string,
	intervalSeconds int,
	quietHoursParams utils.QuietHoursParams,
) *TokenPriceMonitorTask {
	interval := time.Duration(intervalSeconds) * time.Second
	if interval <= 0 {
		interval = 60 * time.Second
	}

	var tokenIds []string
	if tokenIdsStr != "" {
		parts := strings.Split(tokenIdsStr, ",")
		for _, p := range parts {
			if trimmed := strings.TrimSpace(p); trimmed != "" {
				tokenIds = append(tokenIds, trimmed)
			}
		}
	}

	return &TokenPriceMonitorTask{
		tokenService:     tokenService,
		stockClient:      stock.NewStockClient(),
		dingBot:          dingBot,
		stop:             make(chan bool),
		tokenIds:         tokenIds,
		rwaTokenIds:      rwaTokenIds,
		rwaTokenNames:    rwaTokenNames,
		hkStockIds:       hkStockIds,
		hkStockNames:     hkStockNames,
		aStockIds:        aStockIds,
		aStockNames:      aStockNames,
		interval:         interval,
		quietHoursParams: quietHoursParams,
	}
}

func (t *TokenPriceMonitorTask) Start() {
	t.ticker = time.NewTicker(t.interval)
	go func() {
		for {
			select {
			case <-t.ticker.C:
				t.run()
			case <-t.stop:
				t.ticker.Stop()
				return
			}
		}
	}()
}

func (t *TokenPriceMonitorTask) Stop() {
	t.stop <- true
}

func (t *TokenPriceMonitorTask) run() {
	if len(t.tokenIds) == 0 && len(t.rwaTokenIds) == 0 && len(t.hkStockIds) == 0 && len(t.aStockIds) == 0 {
		return
	}

	if !utils.ShouldExecTask(t.quietHoursParams, t.lastRunTime, t.interval) {
		return
	}
	t.lastRunTime = time.Now()

	var allTokenIds []string
	allTokenIds = append(allTokenIds, t.tokenIds...)
	allTokenIds = append(allTokenIds, t.rwaTokenIds...)

	var prices map[string]utils.TokenInfo
	var err error
	if len(allTokenIds) > 0 {
		prices, err = t.tokenService.GetTokenPrice(allTokenIds)
		if err != nil {
			logger.Error("Error fetching token prices: %v", err)
		}
	}

	var cnyPrices map[string]utils.TokenInfo
	hasPaxg := false
	for _, id := range t.tokenIds {
		if id == constant.PAXG_TOKEN_ID {
			hasPaxg = true
			break
		}
	}

	if hasPaxg {
		cnyPrices, err = t.tokenService.GetTokenPrice([]string{constant.PAXG_TOKEN_ID}, "CNY")
		if err != nil {
			logger.Error("Error fetching PAXG token price in CNY: %v", err)
		}
	}

	var stockSymbols []string
	stockSymbols = append(stockSymbols, t.hkStockIds...)
	stockSymbols = append(stockSymbols, t.aStockIds...)

	var stockPrices map[string]stock.StockInfo
	if len(stockSymbols) > 0 {
		stockPrices, err = t.stockClient.GetStockPrices(context.Background(), stockSymbols)
		if err != nil {
			logger.Error("Error fetching stock prices: %v", err)
		}
	}

	formatted, lastUpdated := t.formatTokenPricesDetailed(prices, cnyPrices, stockPrices)

	if formatted == "" {
		return
	}

	// Aggregate messages
	unifiedTitle := fmt.Sprintf("%s Price Alerts", t.dingBot.Keyword)

	unifiedText := fmt.Sprintf("#### %s\n\n%s", unifiedTitle, formatted)
	unifiedText += fmt.Sprintf("\n\n---\n**Last Updated**: %s", utils.FormatBJTime(lastUpdated))

	err = t.dingBot.SendMarkdown(unifiedTitle, unifiedText, nil, false)
	if err != nil {
		logger.Error("Error sending dingtalk message: %v", err)
	} else {
		logger.Info("Sent batch token price alerts")
	}
}

// formatTokenPricesDetailed returns the format used by TokenPriceMonitorTask
func (t *TokenPriceMonitorTask) formatTokenPricesDetailed(
	prices map[string]utils.TokenInfo,
	cnyPrices map[string]utils.TokenInfo,
	stockPrices map[string]stock.StockInfo,
) (string, time.Time) {
	var parts []string
	var maxUpdated time.Time

	// Format Crypto Assets
	if len(t.tokenIds) > 0 && prices != nil {
		var cryptoTexts []string
		for _, tokenId := range t.tokenIds {
			tokenInfo, ok := prices[tokenId]
			if !ok {
				continue
			}
			if tokenInfo.LastUpdated.After(maxUpdated) {
				maxUpdated = tokenInfo.LastUpdated
			}

			if tokenId == constant.PAXG_TOKEN_ID && cnyPrices != nil {
				if cnyInfo, ok := cnyPrices[tokenId]; ok {
					pricePerGram := cnyInfo.Price / 31.1034768
					if cnyInfo.LastUpdated.After(maxUpdated) {
						maxUpdated = cnyInfo.LastUpdated
					}
					text := fmt.Sprintf(
						"- **%s**: ***$%s*** | ***¥%.2f/克*** (%.2f%%)",
						cnyInfo.Symbol, utils.FormatPrice(tokenInfo.Price), pricePerGram, cnyInfo.PercentChange1h)
					cryptoTexts = append(cryptoTexts, text)
					continue
				}
			}

			text := fmt.Sprintf(
				"- **%s**: ***$%s*** (%.2f%%)",
				tokenInfo.Symbol, utils.FormatPrice(tokenInfo.Price), tokenInfo.PercentChange1h)
			cryptoTexts = append(cryptoTexts, text)
		}
		if len(cryptoTexts) > 0 {
			parts = append(parts, "### Crypto Assets\n"+strings.Join(cryptoTexts, "\n---\n"))
		}
	}

	// Format RWA Assets
	if len(t.rwaTokenIds) > 0 && prices != nil {
		var rwaTexts []string
		for _, tokenId := range t.rwaTokenIds {
			tokenInfo, ok := prices[tokenId]
			if !ok {
				continue
			}
			if tokenInfo.LastUpdated.After(maxUpdated) {
				maxUpdated = tokenInfo.LastUpdated
			}

			displayName := tokenInfo.Symbol
			if chineseName, exists := t.rwaTokenNames[tokenId]; exists && chineseName != "" {
				displayName = fmt.Sprintf("%s (%s)", tokenInfo.Symbol, chineseName)
			}

			text := fmt.Sprintf(
				"- **%s**: ***$%s*** (%.2f%%)",
				displayName, utils.FormatPrice(tokenInfo.Price), tokenInfo.PercentChange24h)
			rwaTexts = append(rwaTexts, text)
		}
		if len(rwaTexts) > 0 {
			parts = append(parts, "### RWA Assets\n"+strings.Join(rwaTexts, "\n---\n"))
		}
	}

	// Format HK Stocks
	if len(t.hkStockIds) > 0 && stockPrices != nil {
		var hkTexts []string
		for _, code := range t.hkStockIds {
			info, ok := stockPrices[code]
			if !ok {
				continue
			}
			if info.LastUpdated.After(maxUpdated) {
				maxUpdated = info.LastUpdated
			}

			displayName := info.Name
			if customName, exists := t.hkStockNames[code]; exists && customName != "" {
				displayName = customName
			}

			text := fmt.Sprintf(
				"- **%s (%s)**: ***HK$%.2f*** (%.2f%%)",
				displayName, strings.TrimPrefix(code, "hk"), info.Price, info.PercentChange)
			hkTexts = append(hkTexts, text)
		}
		if len(hkTexts) > 0 {
			parts = append(parts, "### HK Stocks\n"+strings.Join(hkTexts, "\n---\n"))
		}
	}

	// Format A Stocks
	if len(t.aStockIds) > 0 && stockPrices != nil {
		var aTexts []string
		for _, code := range t.aStockIds {
			info, ok := stockPrices[code]
			if !ok {
				continue
			}
			if info.LastUpdated.After(maxUpdated) {
				maxUpdated = info.LastUpdated
			}

			displayName := info.Name
			if customName, exists := t.aStockNames[code]; exists && customName != "" {
				displayName = customName
			}

			text := fmt.Sprintf(
				"- **%s (%s)**: ***¥%.2f*** (%.2f%%)",
				displayName, code, info.Price, info.PercentChange)
			aTexts = append(aTexts, text)
		}
		if len(aTexts) > 0 {
			parts = append(parts, "### A Shares\n"+strings.Join(aTexts, "\n---\n"))
		}
	}

	return strings.Join(parts, "\n\n"), maxUpdated
}
