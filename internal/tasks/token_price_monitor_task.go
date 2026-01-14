package tasks

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ka1fe1/crypto-monitoring/internal/service"
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
		log.Printf("Error fetching token prices: %v", err)
		return
	}

	var allTexts []string

	var lastUpdated time.Time
	for _, tokenId := range t.tokenIds {
		tokenInfo, ok := prices[tokenId]
		if !ok {
			continue
		}
		if tokenInfo.LastUpdated.After(lastUpdated) {
			lastUpdated = tokenInfo.LastUpdated
		}
		text := fmt.Sprintf(
			"- **%s**: ***$%s*** (%.2f%%)\n",
			tokenInfo.Symbol, utils.FormatPrice(tokenInfo.Price), tokenInfo.PercentChange1h)

		allTexts = append(allTexts, text)
	}

	if len(allTexts) == 0 {
		return
	}

	// Aggregate messages
	unifiedTitle := fmt.Sprintf("%s Price Alerts", t.dingBot.Keyword)

	unifiedText := fmt.Sprintf("#### %s\n\n%s", unifiedTitle, strings.Join(allTexts, "\n---\n"))
	unifiedText += fmt.Sprintf("\n\n---\n**Last Updated**: %s", utils.FormatBJTime(lastUpdated))

	err = t.dingBot.SendMarkdown(unifiedTitle, unifiedText, nil, false)
	if err != nil {
		log.Printf("Error sending dingtalk message: %v", err)
	} else {
		log.Printf("Sent batch token price alerts for %d tokens", len(allTexts))
	}
}
