package tasks

import (
	"strings"
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

	qh := utils.QuietHoursParams{Enabled: false, StartHour: 0, EndHour: 7, Behavior: constant.QUIET_HOURS_BEHAVIOR_PAUSE}

	// Create task with a short interval for testing, though we call run() manually
	task := NewTokenPriceMonitorTask(tokenSvc, bot, tokenIdsStr, cfg.TokenPriceMonitor.RwaTokenIDs, cfg.TokenPriceMonitor.RwaTokenNames, 60, qh)

	// Manually trigger run to test logic and notification
	// This will call the real API and send a real DingTalk message is configured
	task.run()
}

func TestFormatTokenPricesDetailed_Paxg(t *testing.T) {
	task := &TokenPriceMonitorTask{}

	prices := map[string]utils.TokenInfo{
		"1":                    {Symbol: "BTC", Price: 60000.0, PercentChange1h: 0.5},
		constant.PAXG_TOKEN_ID: {Symbol: "PAXG", Price: 2150.0, PercentChange1h: -0.1},
	}

	cnyPrices := map[string]utils.TokenInfo{
		constant.PAXG_TOKEN_ID: {Symbol: "PAXG", Price: 15551.7384, PercentChange1h: -0.1},
	}

	tokenIds := []string{"1", constant.PAXG_TOKEN_ID}
	rwaTokenIds := []string{}
	rwaTokenNames := map[string]string{}

	formatted, _ := task.formatTokenPricesDetailed(prices, cnyPrices, tokenIds, rwaTokenIds, rwaTokenNames)

	expectedPaxg := "PAXG**: ***$2150.00*** | ***¥500.00/克*** (-0.10%)"

	if formatted == "" {
		t.Errorf("Expected formatted string, got empty")
	}

	if !strings.Contains(formatted, expectedPaxg) {
		t.Errorf("Expected PAXG to be formatted as $2150.00 | ¥500.00/克, got: %s", formatted)
	}
}

func TestFormatTokenPricesDetailed_RWA(t *testing.T) {
	task := &TokenPriceMonitorTask{}

	prices := map[string]utils.TokenInfo{
		"1":     {Symbol: "BTC", Price: 60000.0, PercentChange1h: 0.5},
		"12345": {Symbol: "ONDO", Price: 1.5, PercentChange1h: 2.3, PercentChange24h: 5.6},
	}

	tokenIds := []string{"1"}
	rwaTokenIds := []string{"12345"}
	rwaTokenNames := map[string]string{
		"12345": "ondo代币",
	}

	formatted, _ := task.formatTokenPricesDetailed(prices, nil, tokenIds, rwaTokenIds, rwaTokenNames)

	if formatted == "" {
		t.Errorf("Expected formatted string, got empty")
	} else {
		t.Log(formatted)
	}

	if !strings.Contains(formatted, "### Crypto Assets") {
		t.Errorf("Expected Crypto Assets section, got: %s", formatted)
	}

	if !strings.Contains(formatted, "### RWA Assets") {
		t.Errorf("Expected RWA Assets section, got: %s", formatted)
	}

	if !strings.Contains(formatted, "ONDO (ondo代币)") {
		t.Errorf("Expected RWA token with Chinese name, got: %s", formatted)
	}

	if !strings.Contains(formatted, "5.60%") {
		t.Errorf("Expected RWA token to show 24h change (5.60%%), got: %s", formatted)
	}
}
