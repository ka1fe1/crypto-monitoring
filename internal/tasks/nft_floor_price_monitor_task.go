package tasks

import (
	"fmt"
	"strings"
	"time"

	"github.com/ka1fe1/crypto-monitoring/pkg/logger"

	"github.com/ka1fe1/crypto-monitoring/internal/service"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/alter/dingding"
)

type NFTFloorPriceMonitorTask struct {
	openSeaService   service.OpenSeaService
	dingBot          *dingding.DingBot
	ticker           *time.Ticker
	stop             chan bool
	collections      []string // Slugs from config
	interval         time.Duration
	quietHoursParams utils.QuietHoursParams
	lastRunTime      time.Time
}

func NewNFTFloorPriceMonitorTask(openSeaService service.OpenSeaService, dingBot *dingding.DingBot, collections []string, intervalSeconds int, quietHoursParams utils.QuietHoursParams) *NFTFloorPriceMonitorTask {
	interval := time.Duration(intervalSeconds) * time.Second
	if interval <= 0 {
		interval = 3600 * time.Second // Default to 1 hour
	}

	return &NFTFloorPriceMonitorTask{
		openSeaService:   openSeaService,
		dingBot:          dingBot,
		stop:             make(chan bool),
		collections:      collections,
		interval:         interval,
		quietHoursParams: quietHoursParams,
	}

}

func (t *NFTFloorPriceMonitorTask) Start() {
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

func (t *NFTFloorPriceMonitorTask) Stop() {
	t.stop <- true
}

func (t *NFTFloorPriceMonitorTask) run() {
	if len(t.collections) == 0 {
		return
	}

	if !utils.ShouldExecTask(t.quietHoursParams, t.lastRunTime, t.interval) {
		return
	}
	t.lastRunTime = time.Now()

	// Fetch prices with USD conversion enabled
	prices, err := t.openSeaService.GetNFTFloorPrices(t.collections, true)
	if err != nil {
		logger.Error("Error fetching NFT floor prices: %v", err)
		return
	}

	var allTexts []string

	for _, info := range prices {
		// Format: - **Slug**: ***1.5 ETH*** ($3000.00)
		priceStr := fmt.Sprintf("%s %s", utils.FormatPrice(info.FloorPrice), info.FloorPriceSymbol)
		usdStr := ""
		if info.FloorPriceUSD > 0 {
			usdStr = fmt.Sprintf(" ($%s)", utils.FormatPrice(info.FloorPriceUSD))
		}

		text := fmt.Sprintf(
			"- **%s**: ***%s***%s\n",
			info.CollectionSlug, priceStr, usdStr)

		allTexts = append(allTexts, text)
	}

	if len(allTexts) == 0 {
		return
	}

	// Aggregate messages
	unifiedTitle := fmt.Sprintf("%s floor price", t.dingBot.Keyword)

	unifiedText := fmt.Sprintf("#### %s\n\n%s", unifiedTitle, strings.Join(allTexts, "\n---\n"))
	unifiedText += fmt.Sprintf("\n\n---\n**Last Updated**: %s", utils.FormatBJTime(time.Now()))

	err = t.dingBot.SendMarkdown(unifiedTitle, unifiedText, nil, false)
	if err != nil {
		logger.Error("Error sending dingtalk message for NFT alerts: %v", err)
	} else {
		logger.Info("Sent batch NFT floor price alerts for %d collections", len(allTexts))
	}
}
