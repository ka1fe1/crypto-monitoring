package polymarket

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ka1fe1/crypto-monitoring/config"
	"github.com/ka1fe1/crypto-monitoring/pkg/logger"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
)

var (
	cfg    *config.Config
	client *Client
)

// loadTestConfig resolves the absolute path to config.yaml and loads it.
func loadTestConfig() (*config.Config, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("failed to get current file path")
	}

	// Current file: /pkg/utils/polymarket/polymarket_test.go
	// Project root: /
	// pkg -> utils -> polymarket -> polymarket_test.go (3 levels up)
	rootDir := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(filename))))

	configPath := filepath.Join(rootDir, "config", "config.yaml")

	return config.LoadConfig(configPath)
}

func TestMain(m *testing.M) {
	var err error
	cfg, err = loadTestConfig()
	if err != nil {
		logger.Warn("Warning: Could not load config: %v", err)
	}

	var apiKey string
	if cfg != nil {
		apiKey = cfg.Polymarket.APIKey
	}

	client = NewClient(apiKey)
	os.Exit(m.Run())
}

func TestGetMarketDetail(t *testing.T) {
	// Use a known market ID (updated by user)
	marketID := "983678"
	market, err := client.GetMarketDetail(marketID)

	if err != nil {
		t.Fatalf("Failed to get market detail: %v", err)
	}

	if market == nil {
		t.Fatal("Market is nil")
	}

	if market.Question == "" {
		t.Error("Market question is empty")
	}

	if market.Slug == "" {
		t.Error("Market slug is empty")
	}

	if len(market.OutcomePrices) == 0 {
		t.Error("Outcome prices are empty")
	}

	t.Log(utils.PrintJson(market))
}
