package tasks

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/ka1fe1/crypto-monitoring/internal/service"
	"github.com/ka1fe1/crypto-monitoring/pkg/logger"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/alter/dingding"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/constant"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/polymarket"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/stock"
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
	rwaTokenIds       []string
	rwaTokenNames     map[string]string
	hkStockIds        []string
	hkStockNames      map[string]string
	aStockIds         []string
	aStockNames       map[string]string
	marketIds         []string
	stockClient       *stock.StockClient
}

func NewGeneralMonitorTask(
	tokenService service.TokenService,
	polymarketService service.PolymarketMonitorService,
	dingBot *dingding.DingBot,
	modules []string,
	tokenIds []string,
	rwaTokenIds []string,
	rwaTokenNames map[string]string,
	hkStockIds []string,
	hkStockNames map[string]string,
	aStockIds []string,
	aStockNames map[string]string,
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
		rwaTokenIds:       rwaTokenIds,
		rwaTokenNames:     rwaTokenNames,
		hkStockIds:        hkStockIds,
		hkStockNames:      hkStockNames,
		aStockIds:         aStockIds,
		aStockNames:       aStockNames,
		marketIds:         marketIds,
		stockClient:       stock.NewStockClient(),
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
	if t.isModuleEnabled("token_price") && (len(t.tokenIds) > 0 || len(t.rwaTokenIds) > 0 || len(t.hkStockIds) > 0 || len(t.aStockIds) > 0) {
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
			logger.Error("Error fetching PAXG token price in CNY for GeneralMonitor: %v", err)
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

	formatted, maxUpdated := t.formatTokenPricesSimple(prices, cnyPrices, stockPrices, t.tokenIds, t.rwaTokenIds, t.rwaTokenNames, t.hkStockIds, t.hkStockNames, t.aStockIds, t.aStockNames)
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

		names := make([]string, 0, len(market.OutcomePrices))
		for name := range market.OutcomePrices {
			names = append(names, name)
		}
		sort.Sort(sort.Reverse(sort.StringSlice(names)))

		prices := make([]string, 0, len(names))
		for _, name := range names {
			prices = append(prices, fmt.Sprintf("%s: %s", name, utils.FormatPrice(market.OutcomePrices[name])))
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

func (t *GeneralMonitorTask) formatTokenPricesSimple(
	prices map[string]utils.TokenInfo,
	cnyPrices map[string]utils.TokenInfo,
	stockPrices map[string]stock.StockInfo,
	tokenIds []string,
	rwaTokenIds []string,
	rwaTokenNames map[string]string,
	hkStockIds []string,
	hkStockNames map[string]string,
	aStockIds []string,
	aStockNames map[string]string,
) (string, time.Time) {
	var parts []string
	var maxUpdated time.Time

	// Format Crypto Assets
	if len(tokenIds) > 0 && prices != nil {
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
			parts = append(parts, "#### Crypto Assets\n"+strings.Join(cryptoTexts, "\n"))
		}
	}

	// Format RWA Assets
	if len(rwaTokenIds) > 0 && prices != nil {
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
			parts = append(parts, "#### RWA Assets\n"+strings.Join(rwaTexts, "\n"))
		}
	}

	// Format HK Stocks
	if len(hkStockIds) > 0 && stockPrices != nil {
		var hkTexts []string
		for _, code := range hkStockIds {
			info, ok := stockPrices[code]
			if !ok {
				continue
			}
			if info.LastUpdated.After(maxUpdated) {
				maxUpdated = info.LastUpdated
			}

			displayName := info.Name
			if overrideName, exists := hkStockNames[code]; exists && overrideName != "" {
				displayName = overrideName
			}

			text := fmt.Sprintf(
				"- **%s**: ***HK$%.2f*** (%.2f%%)",
				displayName, info.Price, info.PercentChange)
			hkTexts = append(hkTexts, text)
		}
		if len(hkTexts) > 0 {
			parts = append(parts, "#### HK Stocks\n"+strings.Join(hkTexts, "\n"))
		}
	}

	// Format A-Shares
	if len(aStockIds) > 0 && stockPrices != nil {
		var aTexts []string
		for _, code := range aStockIds {
			info, ok := stockPrices[code]
			if !ok {
				continue
			}
			if info.LastUpdated.After(maxUpdated) {
				maxUpdated = info.LastUpdated
			}

			displayName := info.Name
			if overrideName, exists := aStockNames[code]; exists && overrideName != "" {
				displayName = overrideName
			}

			text := fmt.Sprintf(
				"- **%s**: ***¥%.2f*** (%.2f%%)",
				displayName, info.Price, info.PercentChange)
			aTexts = append(aTexts, text)
		}
		if len(aTexts) > 0 {
			parts = append(parts, "#### A-Shares\n"+strings.Join(aTexts, "\n"))
		}
	}

	return strings.Join(parts, "\n\n"), maxUpdated
}
