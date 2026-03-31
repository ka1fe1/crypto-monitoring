package tasks

import (
	"testing"
	"time"

	"github.com/ka1fe1/crypto-monitoring/internal/service"
	"github.com/ka1fe1/crypto-monitoring/pkg/logger"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/alter/dingding"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/alternative"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/bgeometrics"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/binance"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/constant"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/mempool"
)

type MockBinanceProvider struct{}

func (m *MockBinanceProvider) GetKlines(symbol, interval string, limit int) ([]binance.Kline, error) {
	klines := make([]binance.Kline, limit)
	for i := 0; i < limit; i++ {
		klines[i] = binance.Kline{
			OpenTime: time.Now().Add(-time.Duration(limit-i) * 24 * time.Hour).UnixMilli(),
			Close:    60000.0,
		}
	}
	return klines, nil
}

type MockMempoolProvider struct{}

func (m *MockMempoolProvider) GetTipHeight() (int64, error) {
	return 840000, nil
}

type MockAlternativeProvider struct{}

func (m *MockAlternativeProvider) GetFng(limit int) (*alternative.FngResponse, error) {
	return &alternative.FngResponse{
		Data: []alternative.FngData{
			{Value: "70", ValueClassification: "Greed"},
		},
	}, nil
}

type MockBalancedPriceProvider struct{}

func (m *MockBalancedPriceProvider) GetBalancedPrice() (float64, error) {
	return 40000.0, nil
}

func TestBtcDashboardMonitorTask_Run_Mock(t *testing.T) {
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

	svc := service.NewBtcDashboardService(
		&MockBinanceProvider{},
		&MockMempoolProvider{},
		&MockAlternativeProvider{},
		&MockBalancedPriceProvider{},
	)

	qh := utils.QuietHoursParams{Enabled: true, StartHour: 11, EndHour: 12, Behavior: constant.QUIET_HOURS_BEHAVIOR_PAUSE}
	task := NewBtcDashboardMonitorTask(svc, bot, cfg.BtcDashboardMonitor.IntervalSeconds, qh)

	logger.Info("Running BtcDashboardMonitorTask (Mock Data) for testing...")
	task.run()

	// Wait momentarily to ensure logs/requests complete in test context
	time.Sleep(1 * time.Second)
}

func TestBtcDashboardMonitorTask_Run_Real(t *testing.T) {
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
	if cfg.BtcDashboardMonitor.BinanceApiUrl != "" {
		binApi = cfg.BtcDashboardMonitor.BinanceApiUrl
	}
	if cfg.BtcDashboardMonitor.MempoolApiUrl != "" {
		memApi = cfg.BtcDashboardMonitor.MempoolApiUrl
	}
	if cfg.BtcDashboardMonitor.AlternativeApiUrl != "" {
		altApi = cfg.BtcDashboardMonitor.AlternativeApiUrl
	}

	bCli := binance.NewClient(binApi)
	mCli := mempool.NewClient(memApi)
	aCli := alternative.NewClient(altApi)

	bgApi := cfg.BtcDashboardMonitor.BgeometricsApiUrl
	bgKey := cfg.BtcDashboardMonitor.BgeometricsApiKey
	bpCli := bgeometrics.NewClient(bgApi, bgKey)

	svc := service.NewBtcDashboardService(bCli, mCli, aCli, bpCli)

	qh := utils.QuietHoursParams{Enabled: true, StartHour: 11, EndHour: 12, Behavior: constant.QUIET_HOURS_BEHAVIOR_PAUSE}
	task := NewBtcDashboardMonitorTask(svc, bot, cfg.BtcDashboardMonitor.IntervalSeconds, qh)

	logger.Info("Running BtcDashboardMonitorTask (Real API) for testing...")
	task.run()

	time.Sleep(1 * time.Second)
}
