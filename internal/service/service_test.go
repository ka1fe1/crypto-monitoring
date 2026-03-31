package service

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ka1fe1/crypto-monitoring/config"
	"github.com/ka1fe1/crypto-monitoring/pkg/logger"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/opensea"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/polymarket"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/twitter"

	"github.com/ka1fe1/crypto-monitoring/pkg/utils/alternative"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/bgeometrics"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/binance"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/mempool"
)

var (
	polySvc         PolymarketMonitorService
	twitterSvc      TwitterService
	tokenSvc        TokenService
	btcDashboardSvc BtcDashboardService
	osSvc           OpenSeaService
)

func loadServiceTestConfig() (*config.Config, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("failed to get current file path")
	}
	rootDir := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	configPath := filepath.Join(rootDir, "config", "config.yaml")
	return config.LoadConfig(configPath)
}

func TestMain(m *testing.M) {
	cfg, err := loadServiceTestConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		// Proceeding with default/empty config or exit?
		// We'll proceed but keys will be empty.
		if cfg == nil {
			cfg = &config.Config{}
		}
	}

	// Setup Logger with DEBUG level for tests
	// If config bas log level, use it, otherwise default to "debug" for tests
	logLevel := "debug"
	if cfg.Log.Level != "" {
		logLevel = cfg.Log.Level
	}
	logger.Setup(logLevel)

	// 1. Setup Polymarket Service
	polyClient := polymarket.NewClient(cfg.Polymarket.APIKey)
	// User requested to remove mock transport to allow real network interaction
	polySvc = NewPolymarketMonitorService(polyClient)

	// 2. Setup Twitter Service
	twitterClient := twitter.NewTwitterClient(cfg.Twitter.APIKey)
	twitterSvc = NewTwitterService(twitterClient)

	// 3. Setup Token Service
	tokenClient := utils.NewCoinMarketClient(cfg.CoinMarketCap.APIKey)
	tokenSvc = NewTokenService(tokenClient)

	// 4. Setup OpenSea Service (Integrating existing test setup)
	osClient := opensea.NewOpenSeaClient(cfg.OpenSea.APIKey)
	// OpenSea client doesn't seem to have SetHttpClient in the snippet I saw,
	// but hopefully NewOpenSeaService accepts it.
	// opensea_service_test.go uses: svc = NewOpenSeaService(osClient, cmcClient)
	// We use the same cmcClient as tokenSvc? Yes, likely fine.
	svc = NewOpenSeaService(osClient, tokenClient)

	// 5. Setup Btc Dashboard Service
	binApi, memApi, altApi, bgApi, bgKey := "https://api.binance.com", "https://mempool.space", "https://api.alternative.me", "", ""
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
		bgApi = cfg.BtcDashboardMonitor.BgeometricsApiUrl
		bgKey = cfg.BtcDashboardMonitor.BgeometricsApiKey
	}

	bCli := binance.NewClient(binApi)
	mCli := mempool.NewClient(memApi)
	aCli := alternative.NewClient(altApi)
	bpCli := bgeometrics.NewClient(bgApi, bgKey)

	btcDashboardSvc = NewBtcDashboardService(bCli, mCli, aCli, bpCli)

	os.Exit(m.Run())
}
