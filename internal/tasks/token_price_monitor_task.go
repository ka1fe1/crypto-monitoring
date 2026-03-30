package tasks

import (
	"fmt"
	"strings"
	"time"

	"github.com/ka1fe1/crypto-monitoring/internal/service"
	"github.com/ka1fe1/crypto-monitoring/pkg/logger"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/alter/dingding"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/constant"
)

type TokenPriceMonitorTask struct {
	tokenService     service.TokenService
	dingBot          *dingding.DingBot
	ticker           *time.Ticker
	stop             chan bool
	tokenIds         []string
	rwaTokenIds      []string
	rwaTokenNames    map[string]string
	interval         time.Duration
	quietHoursParams utils.QuietHoursParams
	lastRunTime      time.Time
}

func NewTokenPriceMonitorTask(tokenService service.TokenService, dingBot *dingding.DingBot, tokenIdsStr string, rwaTokenIds []string, rwaTokenNames map[string]string, intervalSeconds int, quietHoursParams utils.QuietHoursParams) *TokenPriceMonitorTask {
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
		dingBot:          dingBot,
		stop:             make(chan bool),
		tokenIds:         tokenIds,
		rwaTokenIds:      rwaTokenIds,
		rwaTokenNames:    rwaTokenNames,
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
	if len(t.tokenIds) == 0 && len(t.rwaTokenIds) == 0 {
		return
	}

	if !utils.ShouldExecTask(t.quietHoursParams, t.lastRunTime, t.interval) {
		return
	}
	t.lastRunTime = time.Now()

	var allTokenIds []string
	allTokenIds = append(allTokenIds, t.tokenIds...)
	allTokenIds = append(allTokenIds, t.rwaTokenIds...)

	prices, err := t.tokenService.GetTokenPrice(allTokenIds)
	if err != nil {
		logger.Error("Error fetching token prices: %v", err)
		return
	}

	if len(prices) == 0 {
		return
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

	formatted, lastUpdated := t.formatTokenPricesDetailed(prices, cnyPrices, t.tokenIds, t.rwaTokenIds, t.rwaTokenNames)

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
func (t *TokenPriceMonitorTask) formatTokenPricesDetailed(prices map[string]utils.TokenInfo, cnyPrices map[string]utils.TokenInfo, tokenIds []string, rwaTokenIds []string, rwaTokenNames map[string]string) (string, time.Time) {
	var parts []string
	var maxUpdated time.Time

	// Format Crypto Assets
	if len(tokenIds) > 0 {
		var cryptoTexts []string
		for _, tokenId := range tokenIds {
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
	if len(rwaTokenIds) > 0 {
		var rwaTexts []string
		for _, tokenId := range rwaTokenIds {
			tokenInfo, ok := prices[tokenId]
			if !ok {
				continue
			}
			if tokenInfo.LastUpdated.After(maxUpdated) {
				maxUpdated = tokenInfo.LastUpdated
			}

			displayName := tokenInfo.Symbol
			if chineseName, exists := rwaTokenNames[tokenId]; exists && chineseName != "" {
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

	return strings.Join(parts, "\n\n"), maxUpdated
}
