package tasks

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/alter/dingding"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/polymarket"
)

type PolymarketMonitorTask struct {
	client    *polymarket.Client
	dingBot   *dingding.DingBot
	ticker    *time.Ticker
	stop      chan bool
	marketIDs []string
	interval  time.Duration
}

func NewPolymarketMonitorTask(client *polymarket.Client, dingBot *dingding.DingBot, marketIDs []string, intervalSeconds int) *PolymarketMonitorTask {
	interval := time.Duration(intervalSeconds) * time.Second
	if interval <= 0 {
		interval = 3600 * time.Second
	}

	return &PolymarketMonitorTask{
		client:    client,
		dingBot:   dingBot,
		stop:      make(chan bool),
		marketIDs: marketIDs,
		interval:  interval,
	}
}

func (t *PolymarketMonitorTask) Start() {
	t.ticker = time.NewTicker(t.interval)
	log.Printf("Starting Polymarket Monitor Task with interval %v, monitoring %d markets", t.interval, len(t.marketIDs))
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

	var allTexts []string
	for _, id := range t.marketIDs {
		market, err := t.client.GetMarketDetail(id)
		if err != nil {
			log.Printf("Error fetching Polymarket detail for ID %s: %v", id, err)
			continue
		}

		// Skip closed markets
		if market.Closed {
			continue
		}

		// Build price string: Yes: 0.75, No: 0.25
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

		allTexts = append(allTexts, text)
	}

	if len(allTexts) == 0 {
		return
	}

	title := fmt.Sprintf("%s Polymarket Monitor", t.dingBot.Keyword)
	content := fmt.Sprintf("## %s\n\n --- \n\n%s\n\n---\n**Last Updated**: %s",
		title,
		strings.Join(allTexts, "\n\n---\n\n"),
		utils.FormatBJTime(time.Now()),
	)

	err := t.dingBot.SendMarkdown(title, content, nil, false)
	if err != nil {
		log.Printf("Error sending DingTalk message for Polymarket monitor: %v", err)
	} else {
		log.Printf("Sent Polymarket monitor update for %d markets", len(allTexts))
	}
}

func (t *PolymarketMonitorTask) getClosedStr(closed bool) string {
	if closed {
		return "Closed"
	}
	return "Active"
}
