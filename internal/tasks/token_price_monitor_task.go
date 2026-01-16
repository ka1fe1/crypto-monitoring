package tasks

import (
	"fmt"
	"strings"
	"time"

	"github.com/ka1fe1/crypto-monitoring/internal/service"
	"github.com/ka1fe1/crypto-monitoring/pkg/logger"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/alter/dingding"
)

type TokenPriceMonitorTask struct {
	tokenService     service.TokenService
	dingBot          *dingding.DingBot
	ticker           *time.Ticker
	stop             chan bool
	tokenIds         []string
	interval         time.Duration
	quietHoursParams utils.QuietHoursParams
	lastRunTime      time.Time
}

func NewTokenPriceMonitorTask(tokenService service.TokenService, dingBot *dingding.DingBot, tokenIdsStr string, intervalSeconds int, quietHoursParams utils.QuietHoursParams) *TokenPriceMonitorTask {
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
	if len(t.tokenIds) == 0 {
		return
	}

	if !utils.ShouldExecTask(t.quietHoursParams, t.lastRunTime, t.interval) {
		return
	}
	t.lastRunTime = time.Now()

	prices, err := t.tokenService.GetTokenPrice(t.tokenIds)
	if err != nil {
		logger.Error("Error fetching token prices: %v", err)
		return
	}

	if len(prices) == 0 {
		return
	}

	// This is using the detailed format which matches the previous logic
	// But wait, the previous logic added newlines and joined with "\n---\n"
	// TokenService.FormatTokenPricesDetailed does exactly that.
	// But it returns the joined string, not the slice.
	formatted, lastUpdated := t.formatTokenPricesDetailed(prices, t.tokenIds)

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
func (t *TokenPriceMonitorTask) formatTokenPricesDetailed(prices map[string]utils.TokenInfo, tokenIds []string) (string, time.Time) {
	var texts []string
	var maxUpdated time.Time

	for _, tokenId := range tokenIds {
		tokenInfo, ok := prices[tokenId]
		if !ok {
			continue
		}
		if tokenInfo.LastUpdated.After(maxUpdated) {
			maxUpdated = tokenInfo.LastUpdated
		}
		text := fmt.Sprintf(
			"- **%s**: ***$%s*** (%.2f%%)\n",
			tokenInfo.Symbol, utils.FormatPrice(tokenInfo.Price), tokenInfo.PercentChange1h)
		texts = append(texts, text)
	}

	return strings.Join(texts, "\n---\n"), maxUpdated
}
