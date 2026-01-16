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

type PolymarketMonitorTask struct {
	service          service.PolymarketMonitorService
	dingBot          *dingding.DingBot
	ticker           *time.Ticker
	stop             chan bool
	marketIDs        []string
	interval         time.Duration
	quietHoursParams utils.QuietHoursParams
	lastRunTime      time.Time
}

func NewPolymarketMonitorTask(service service.PolymarketMonitorService, dingBot *dingding.DingBot, marketIDs []string, intervalSeconds int, quietHoursParams utils.QuietHoursParams) *PolymarketMonitorTask {
	interval := time.Duration(intervalSeconds) * time.Second
	if interval <= 0 {
		interval = 3600 * time.Second
	}

	return &PolymarketMonitorTask{
		service:          service,
		dingBot:          dingBot,
		stop:             make(chan bool),
		marketIDs:        marketIDs,
		interval:         interval,
		quietHoursParams: quietHoursParams,
	}

}

func (t *PolymarketMonitorTask) Start() {
	t.ticker = time.NewTicker(t.interval)
	logger.Info("Starting Polymarket Monitor Task with interval %v, monitoring %d markets", t.interval, len(t.marketIDs))
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

func (t *PolymarketMonitorTask) Stop() {
	t.stop <- true
}

func (t *PolymarketMonitorTask) run() {
	if len(t.marketIDs) == 0 {
		return
	}

	if !utils.ShouldExecTask(t.quietHoursParams, t.lastRunTime, t.interval) {
		logger.Info("Skipping Polymarket Monitor Task for %s in quiet hours", t.dingBot.Keyword)
		return
	}
	t.lastRunTime = time.Now()

	markets, err := t.service.GetMarketDetails(t.marketIDs)
	if err != nil {
		logger.Error("Error fetching Polymarket details: %v", err)
	}

	if len(markets) == 0 {
		return
	}

	if len(markets) == 0 {
		return
	}

	content := t.formatMarkets(markets)
	if content == "" {
		return
	}

	title := fmt.Sprintf("%s Polymarket Monitor", t.dingBot.Keyword)
	// We need to construct the full markdown like before
	fullContent := fmt.Sprintf("## %s\n\n --- \n\n%s\n\n---\n**Last Updated**: %s",
		title,
		content,
		utils.FormatBJTime(time.Now()),
	)

	err = t.dingBot.SendMarkdown(title, fullContent, nil, false)
	if err != nil {
		logger.Error("Error sending DingTalk message for Polymarket monitor: %v", err)
	} else {
		logger.Info("Sent Polymarket monitor update for %d markets", len(markets))
	}
}

func (t *PolymarketMonitorTask) formatMarkets(markets []polymarket.MarketDetail) string {
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
			"### %s\n"+
				"- **Status**: %s | **Volume**: $%s\n"+
				"- **Outcome Prices**: %s\n"+
				"- **1H Change**: %s%%",
			market.Question,
			t.getClosedStr(market.Closed),
			utils.FormatPrice(market.Volume),
			strings.Join(prices, " | "),
			utils.FormatPrice(market.OneHourPriceChange*100),
		)
		texts = append(texts, text)
	}
	return strings.Join(texts, "\n\n---\n\n")
}

func (t *PolymarketMonitorTask) getClosedStr(closed bool) string {
	if closed {
		return "Closed"
	}
	return "Active"
}
