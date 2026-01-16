package tasks

import (
	"testing"

	"github.com/ka1fe1/crypto-monitoring/internal/service"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/alter/dingding"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/constant"
)

func TestTokenPriceMonitorTask_Run(t *testing.T) {
	if cfg == nil {
		t.Skip("Config not loaded, skipping test")
	}

	botName := cfg.TokenPriceMonitor.BotName
	if botName == "" {
		botName = constant.DEFAULT_BOT_NAME
	}

	botCfg, ok := cfg.DingTalk[botName]
	if !ok {
		t.Skipf("DingTalk bot %s not configured, skipping test", botName)
	}

	bot := dingding.NewDingBot(botCfg.AccessToken, botCfg.Secret, botCfg.Keyword)

	tokenClient := utils.NewCoinMarketClient(cfg.CoinMarketCap.APIKey)
	tokenSvc := service.NewTokenService(tokenClient)

	tokenIdsStr := cfg.TokenPriceMonitor.TokenIds
	if tokenIdsStr == "" {
		// Default to BTC if not configured
		tokenIdsStr = "1"
	}

	qh := utils.QuietHoursParams{Enabled: true, StartHour: 0, EndHour: 7, Behavior: constant.QUIET_HOURS_BEHAVIOR_PAUSE}

	// Create task with a short interval for testing, though we call run() manually
	task := NewTokenPriceMonitorTask(tokenSvc, bot, tokenIdsStr, 60, qh)

	// Manually trigger run to test logic and notification
	// This will call the real API and send a real DingTalk message is configured
	task.run()
}
