package tasks

import (
	"testing"
	"time"

	"github.com/ka1fe1/crypto-monitoring/internal/service"
	"github.com/ka1fe1/crypto-monitoring/pkg/logger"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/alter/dingding"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/constant"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/alternative"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/binance"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/mempool"
)

func TestBtcDashboardMonitorTask_Run(t *testing.T) {
	if cfg == nil {
		t.Skip("Config not loaded, skipping test")
	}

	botName := cfg.BtcDashboardMonitor.BotName
	if botName == "" {
		botName = "token" // Default fallback
	}

	botCfg, ok := cfg.DingTalk[botName]
	if !ok {
		t.Skipf("DingTalk bot %s not configured, skipping test", botName)
	}

	bot := dingding.NewDingBot(botCfg.AccessToken, botCfg.Secret, botCfg.Keyword)

	binApi, memApi, altApi := "https://api.binance.com", "https://mempool.space", "https://api.alternative.me"
	if cfg != nil {
		if cfg.BtcDashboardMonitor.BinanceApiUrl != "" {
			binApi = cfg.BtcDashboardMonitor.BinanceApiUrl
		}
		if cfg.BtcDashboardMonitor.MempoolApiUrl != "" {
			memApi = cfg.BtcDashboardMonitor.MempoolApiUrl
		}
		if cfg.BtcDashboardMonitor.AlternativeApiUrl != "" {
			altApi = cfg.BtcDashboardMonitor.AlternativeApiUrl
		}
	}
	bCli := binance.NewClient(binApi)
	mCli := mempool.NewClient(memApi)
	aCli := alternative.NewClient(altApi)
	svc := service.NewBtcDashboardService(bCli, mCli, aCli)

	qh := utils.QuietHoursParams{Enabled: true, StartHour: 11, EndHour: 12, Behavior: constant.QUIET_HOURS_BEHAVIOR_PAUSE}
	task := NewBtcDashboardMonitorTask(svc, bot, cfg.BtcDashboardMonitor.IntervalSeconds, qh)

	logger.Info("Running BtcDashboardMonitorTask for testing...")
	task.run()
	
	// Wait momentarily to ensure logs/requests complete in test context
	time.Sleep(1 * time.Second)
}
