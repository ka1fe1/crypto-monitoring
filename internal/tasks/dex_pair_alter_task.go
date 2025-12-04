package tasks

import (
	"fmt"
	"log"
	"time"

	"github.com/ka1fe1/crypto-monitoring/internal/service"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/alter/dingding"
)

type DexPairAlterTask struct {
	dexService      service.DexPairService
	dingBot         *dingding.DingBot
	ticker          *time.Ticker
	stop            chan bool
	contractAddress string
	networkSlug     string
	interval        time.Duration
}

func NewDexPairAlterTask(dexService service.DexPairService, dingBot *dingding.DingBot, contractAddress, networkSlug string, intervalSeconds int) *DexPairAlterTask {
	interval := time.Duration(intervalSeconds) * time.Second
	if interval <= 0 {
		interval = 60 * time.Second
	}
	return &DexPairAlterTask{
		dexService:      dexService,
		dingBot:         dingBot,
		stop:            make(chan bool),
		contractAddress: contractAddress,
		networkSlug:     networkSlug,
		interval:        interval,
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
	info, err := t.dexService.GetDexPairInfo(t.contractAddress, t.networkSlug)
	if err != nil {
		log.Printf("Error fetching dex pair info: %v", err)
		return
	}

	title := fmt.Sprintf("Price Alert: %s", info.Name)
	text := fmt.Sprintf("### %s Price Alert\n\n"+
		"- **Price**: $%.6f\n"+
		"- **1h Change**: %.4f%%\n"+
		"- **24h Change**: %.4f%%\n"+
		"- **Liquidity**: $%s\n"+
		"- **DEX**: %s\n",
		info.Name, info.Price, info.PercentChange1h, info.PercentChange24h, formatLiquidity(info.Liquidity), info.DexSlug)

	err = t.dingBot.SendMarkdown(title, text, nil, false)
	if err != nil {
		log.Printf("Error sending dingtalk message: %v", err)
	} else {
		log.Printf("Sent price alert for %s", info.Name)
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
