package tasks

import (
	"log"
	"strings"
	"testing"

	"github.com/ka1fe1/crypto-monitoring/internal/service"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/alter/dingding"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/polymarket"
)

func TestGeneralMonitorTask_Run(t *testing.T) {
	if cfg == nil {
		t.Skip("Config not loaded")
	}

	// 1. Initialize CoinMarketCap
	cmcClient := utils.NewCoinMarketClient(cfg.CoinMarketCap.APIKey)
	tokenService := service.NewTokenService(cmcClient)

	// 2. Initialize Polymarket
	polyClient := polymarket.NewClient(cfg.Polymarket.APIKey)
	polymarketService := service.NewPolymarketMonitorService(polyClient)

	// 3. Initialize DingBot
	botName := cfg.GeneralMonitor.BotName
	// If bot name is empty, try to find a valid one from config
	if botName == "" && len(cfg.DingTalk) > 0 {
		for k := range cfg.DingTalk {
			botName = k
			break
		}
	}

	botConfig, ok := cfg.DingTalk[botName]
	if !ok {
		// Log warning and maybe skip if no bot configured?
		log.Printf("Warning: Bot %s not found in config", botName)
		// We might fail here or proceed if we accept nil bot (but Start() might panic or log error)
		if len(cfg.DingTalk) > 0 {
			for _, v := range cfg.DingTalk {
				botConfig = v
				break
			}
		} else {
			t.Skip("No valid DingTalk config found")
		}
	}

	dingBot := dingding.NewDingBot(botConfig.AccessToken, botConfig.Secret, botConfig.Keyword)

	// 4. Parse Token IDs from TokenPriceMonitor config
	var tokenIds []string
	if cfg.TokenPriceMonitor.TokenIds != "" {
		for _, p := range strings.Split(cfg.TokenPriceMonitor.TokenIds, ",") {
			if t := strings.TrimSpace(p); t != "" {
				tokenIds = append(tokenIds, t)
			}
		}
	}

	// 5. Create Task
	task := NewGeneralMonitorTask(
		tokenService,
		polymarketService,
		dingBot,
		cfg.GeneralMonitor.Modules,
		tokenIds,
		cfg.PolymarketMonitor.MarketIDs,
		60,
		utils.QuietHoursParams{Enabled: false},
	)

	// 6. Run Task
	task.run()
}
