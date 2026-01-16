package tasks

import (
	"fmt"
	"strings"
	"time"

	"github.com/ka1fe1/crypto-monitoring/internal/service"
	"github.com/ka1fe1/crypto-monitoring/pkg/logger"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/alter/dingding"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/polymarket"
)

type GeneralMonitorTask struct {
	tokenService      service.TokenService
	polymarketService service.PolymarketMonitorService
	dingBot           *dingding.DingBot
	ticker            *time.Ticker
	stop              chan bool
	modules           []string
	interval          time.Duration
	quietHoursParams  utils.QuietHoursParams
	lastRunTime       time.Time
	tokenIds          []string
	marketIds         []string
}

func NewGeneralMonitorTask(
	tokenService service.TokenService,
	polymarketService service.PolymarketMonitorService,
	dingBot *dingding.DingBot,
	modules []string,
	tokenIds []string,
	marketIds []string,
	intervalSeconds int,
	quietHoursParams utils.QuietHoursParams,
) *GeneralMonitorTask {
	interval := time.Duration(intervalSeconds) * time.Second
	if interval <= 0 {
		interval = 60 * time.Second
	}

	return &GeneralMonitorTask{
		tokenService:      tokenService,
		polymarketService: polymarketService,
		dingBot:           dingBot,
		stop:              make(chan bool),
		modules:           modules,
		interval:          interval,
		quietHoursParams:  quietHoursParams,
		tokenIds:          tokenIds,
		marketIds:         marketIds,
	}
}

func (t *GeneralMonitorTask) Start() {
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

func (t *GeneralMonitorTask) Stop() {
	t.stop <- true
}

func (t *GeneralMonitorTask) run() {
	if !utils.ShouldExecTask(t.quietHoursParams, t.lastRunTime, t.interval) {
		return
	}
	t.lastRunTime = time.Now()

	var parts []string
	var lastUpdated time.Time

	// 1. Token Price Module
	if t.isModuleEnabled("token_price") && len(t.tokenIds) > 0 {
		tokenPart, updated, err := t.getTokenPriceContent()
		if err == nil && tokenPart != "" {
			parts = append(parts, tokenPart)
			if updated.After(lastUpdated) {
				lastUpdated = updated
			}
		} else if err != nil {
			logger.Error("Error in GeneralMonitorTask token_price: %v", err)
		}
	}

	// 3. Polymarket Module
	if t.isModuleEnabled("polymarket") && len(t.marketIds) > 0 {
		polyPart, updated, err := t.getPolymarketContent()
		if err == nil && polyPart != "" {
			parts = append(parts, polyPart)
			// Polymarket doesn't explicitly return updated time in market obj easily, usually it's fetch time.
			// getPolymarketContent will return Now() or closest.
			if updated.After(lastUpdated) {
				lastUpdated = updated
			}
		} else if err != nil {
			logger.Error("Error in GeneralMonitorTask polymarket: %v", err)
		}
	}

	if len(parts) == 0 {
		return
	}

	// Aggregate messages
	unifiedTitle := fmt.Sprintf("%s General Update", t.dingBot.Keyword)
	unifiedText := fmt.Sprintf("## %s\n\n%s", unifiedTitle, strings.Join(parts, "\n\n---\n\n"))
	unifiedText += fmt.Sprintf("\n\n---\n**Last Updated**: %s", utils.FormatBJTime(lastUpdated))

	err := t.dingBot.SendMarkdown(unifiedTitle, unifiedText, nil, false)
	if err != nil {
		logger.Error("Error sending general monitor message: %v", err)
	} else {
		logger.Info("Sent general monitor message with %d parts", len(parts))
	}
}

func (t *GeneralMonitorTask) isModuleEnabled(module string) bool {
	for _, m := range t.modules {
		if strings.EqualFold(m, module) {
			return true
		}
	}
	return false
}

func (t *GeneralMonitorTask) getTokenPriceContent() (string, time.Time, error) {
	prices, err := t.tokenService.GetTokenPrice(t.tokenIds)
	if err != nil {
		return "", time.Time{}, err
	}

	formatted, maxUpdated := t.formatTokenPricesSimple(prices, t.tokenIds)
	if formatted == "" {
		return "", time.Time{}, nil
	}

	content := "### Token Prices\n" + formatted
	return content, maxUpdated, nil
}

func (t *GeneralMonitorTask) getPolymarketContent() (string, time.Time, error) {
	markets, err := t.polymarketService.GetMarketDetails(t.marketIds)
	if err != nil {
		// Log but use whatever we got
		logger.Error("GeneralMonitor: Error fetching polymarket: %v", err)
	}

	if len(markets) == 0 {
		return "", time.Time{}, nil
	}

	formatted := t.formatPolymarketMarkets(markets)
	if formatted == "" {
		return "", time.Time{}, nil
	}

	content := "### Polymarket\n" + formatted
	return content, time.Now(), nil
}

func (t *GeneralMonitorTask) formatPolymarketMarkets(markets []polymarket.MarketDetail) string {
	var texts []string
	for _, market := range markets {
		if market.Closed {
			continue
		}

		prices := []string{}
		for name, price := range market.OutcomePrices {
			prices = append(prices, fmt.Sprintf("%s: %s", name, utils.FormatPrice(price)))
		}

		text := fmt.Sprintf(
			"- **%s** ($%s)\n  %s",
			market.Question,
			utils.FormatPrice(market.Volume),
			strings.Join(prices, " | "),
		)
		texts = append(texts, text)
	}
	return strings.Join(texts, "\n")
}

func (t *GeneralMonitorTask) formatTokenPricesSimple(prices map[string]utils.TokenInfo, tokenIds []string) (string, time.Time) {
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
			"- **%s**: ***$%s*** (%.2f%%)",
			tokenInfo.Symbol, utils.FormatPrice(tokenInfo.Price), tokenInfo.PercentChange1h)
		texts = append(texts, text)
	}

	return strings.Join(texts, "\n"), maxUpdated
}
