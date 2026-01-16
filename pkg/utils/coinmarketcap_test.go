package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ka1fe1/crypto-monitoring/config"
	"github.com/ka1fe1/crypto-monitoring/pkg/logger"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/constant"
)

var client *CoinMarketClient

// loadTestConfig resolves the absolute path to config.yaml and loads it.
func loadTestConfig() (*config.Config, error) {
	// 1. Get the absolute path of the current file to determine the project root.
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("failed to get current file path")
	}

	// The current file is in <ProjectRoot>/pkg/utils/coinmarketcap_test.go
	// So we go up two levels to get to <ProjectRoot>
	rootDir := filepath.Dir(filepath.Dir(filepath.Dir(filename)))

	// 2. Construct the absolute path to config.yaml
	configPath := filepath.Join(rootDir, "config", "config.yaml")

	// 3. Load the configuration
	return config.LoadConfig(configPath)
}

// TestMain acts as the entry point for running tests in this package.
// It allows us to perform global setup (like loading config) before running any tests.
func TestMain(m *testing.M) {
	// Load the configuration using the helper method
	cfg, err := loadTestConfig()
	if err != nil {
		logger.Warn("Warning: Could not load config: %v", err)
	}

	// Initialize the client if config was loaded successfully
	if cfg != nil {
		client = NewCoinMarketClient(cfg.CoinMarketCap.APIKey)
	}

	// Run the tests
	exitCode := m.Run()

	// Exit with the code returned by m.Run()
	os.Exit(exitCode)
}

func TestGetPrice(t *testing.T) {
	// If client wasn't initialized (e.g. config missing), skip this test
	if client == nil {
		t.Skip("Skipping TestGetPrice: client not initialized (check config.yaml path and content)")
	}

	// Use the client initialized in TestMain
	// BTC ID: 1, ETH ID: 1027
	ids := []string{"1", "1027"}
	prices, err := client.GetPrice(ids)
	if err != nil {
		t.Fatalf("GetPrice failed: %v", err)
	}

	for id, info := range prices {
		t.Logf("ID: %s, Price: %f, Symbol: %s", id, info.Price, info.Symbol)
		if info.Price <= 0 {
			t.Errorf("Price for ID %s should be greater than 0", id)
		}
		if info.Symbol == "" {
			t.Errorf("Symbol for ID %s should not be empty", id)
		}
	}
}

func TestGetDexPairQuotes(t *testing.T) {
	// If client wasn't initialized (e.g. config missing), skip this test
	if client == nil {
		t.Skip("Skipping TestGetDexPairQuotes: client not initialized")
	}

	// Example: WETH/USDC pair on Ethereum
	// Contract Address: 0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640 (Uniswap V3)
	contractAddress := "0xb67e5eaf770a384ab28029d08b9bc5ebe32beb0f"
	networkId := constant.BNB_NETWORK_ID

	pairs, err := client.GetDexPairQuotes([]string{contractAddress}, "", networkId)
	if err != nil {
		t.Fatalf("GetDexPairQuotes failed: %v", err)
	}

	pair, ok := pairs[contractAddress]
	if !ok {
		t.Fatalf("Pair %s not found in response", contractAddress)
	}

	t.Logf("Pair Name: %s", pair.Name)
	t.Logf("Base Token: %s (%s)", pair.BaseAssetSymbol, pair.BaseAssetContractAddress)
	t.Logf("Quote Token: %s (%s)", pair.QuoteAssetSymbol, pair.QuoteAssetContractAddress)

	if len(pair.Quote) > 0 {
		t.Logf("info: %s", PrintJson(pair))
	} else {
		t.Log("No quote data available")
	}

}
