package tasks

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ka1fe1/crypto-monitoring/internal/service"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/alter/dingding"
)

type DexPairAlterTask struct {
	dexService       service.DexPairService
	dingBot          *dingding.DingBot
	ticker           *time.Ticker
	stop             chan bool
	contractAddrInfo map[string][]string
	interval         time.Duration
}

func NewDexPairAlterTask(dexService service.DexPairService, dingBot *dingding.DingBot, contractAddrInfo map[string][]string, intervalSeconds int) *DexPairAlterTask {
	interval := time.Duration(intervalSeconds) * time.Second
	if interval <= 0 {
		interval = 60 * time.Second
	}
	return &DexPairAlterTask{
		dexService:       dexService,
		dingBot:          dingBot,
		stop:             make(chan bool),
		contractAddrInfo: contractAddrInfo,
		interval:         interval,
	}
}

func (t *DexPairAlterTask) Start() {
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

func (t *DexPairAlterTask) Stop() {
	t.stop <- true
}

func (t *DexPairAlterTask) run() {
	var allTexts []string
	var allTitles []string

	for networkId, addrs := range t.contractAddrInfo {
		if len(addrs) == 0 {
			continue
		}

		infos, err := t.dexService.GetDexPairInfo(addrs, "", networkId)
		if err != nil {
			log.Printf("Error fetching dex pair info for network %s: %v", networkId, err)
			continue
		}

		for _, info := range infos {
			if info == nil {
				continue
			}
			title := fmt.Sprintf("Price Alert: %s", info.Name)
			text := fmt.Sprintf("### %s Price Alert\n\n"+
				"- **Price**: $%.6f\n"+
				"- **Liquidity**: $%s\n"+
				"- **1h Change**: %.8f%%\n"+
				"- **DEX & Network**: %s\n"+
				"- **Last Updated**: %s\n",
				info.Name, info.Price, formatLiquidity(info.Liquidity), info.PercentChange1h,
				fmt.Sprintf("%s & %s", info.DexSlug, info.NetworkSlug), info.LastUpdated)

			allTitles = append(allTitles, title)
			allTexts = append(allTexts, text)
		}
	}

	if len(allTexts) == 0 {
		return
	}

	// Aggregate messages
	unifiedTitle := "Batch Price Alerts"
	if len(allTitles) > 0 {
		unifiedTitle = allTitles[0] + "..." // Simple summary title
	}

	unifiedText := strings.Join(allTexts, "\n---\n")

	err := t.dingBot.SendMarkdown(unifiedTitle, unifiedText, nil, false)
	if err != nil {
		log.Printf("Error sending dingtalk message: %v", err)
	} else {
		log.Printf("Sent batch price alerts for %d pairs", len(allTexts))
	}
}

func formatLiquidity(val float64) string {
	var formatted string
	if val >= 1e9 {
		formatted = fmt.Sprintf("%.2fB", val/1e9)
	} else if val >= 1e6 {
		formatted = fmt.Sprintf("%.2fM", val/1e6)
	} else if val >= 1e3 {
		formatted = fmt.Sprintf("%.2fK", val/1e3)
	} else {
		formatted = fmt.Sprintf("%.2f", val)
	}
	return fmt.Sprintf("%s(%.0f)", formatted, val)
}
