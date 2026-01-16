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
)

var (
	polySvc    PolymarketMonitorService
	twitterSvc TwitterMonitorService
	tokenSvc   TokenService
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
	twitterSvc = NewTwitterMonitorService(twitterClient)

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

	os.Exit(m.Run())
}
