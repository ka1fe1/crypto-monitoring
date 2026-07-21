package stock

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ka1fe1/crypto-monitoring/config"
)

var (
	cfg *config.Config
)

// loadTestConfig resolves the absolute path to config.yaml and loads it.
func loadTestConfig() (*config.Config, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("failed to get current file path")
	}

	// Current file is in <ProjectRoot>/pkg/utils/stock/stock_test.go
	// Go up 4 levels to get <ProjectRoot>
	rootDir := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(filename))))

	configPath := filepath.Join(rootDir, "config", "config.yaml")

	return config.LoadConfig(configPath)
}

func TestMain(m *testing.M) {
	var err error
	cfg, err = loadTestConfig()
	if err != nil {
		log.Fatalf("Critical: Could not load config.yaml: %v", err)
	}

	os.Exit(m.Run())
}

func TestGetStockPrices(t *testing.T) {
	t.Log("xxx")
	if cfg == nil {
		t.Fatalf("cfg is nil! Test environment failed to load config.yaml")
	}

	t.Logf("Successfully loaded config.yaml. HkStockIDs (raw: %q): %v, AStockIDs (raw: %q): %v",
		cfg.TokenPriceMonitor.HkStockIds, cfg.TokenPriceMonitor.HkStockIDs,
		cfg.TokenPriceMonitor.AStockIds, cfg.TokenPriceMonitor.AStockIDs)

	if len(cfg.TokenPriceMonitor.HkStockIDs) == 0 && len(cfg.TokenPriceMonitor.AStockIDs) == 0 {
		t.Fatalf("Both HkStockIDs and AStockIDs are empty in config.yaml under token_price_monitor!")
	}

	client := NewStockClient()

	var symbols []string
	symbols = append(symbols, cfg.TokenPriceMonitor.AStockIDs...)
	symbols = append(symbols, cfg.TokenPriceMonitor.HkStockIDs...)

	prices, err := client.GetStockPrices(context.Background(), symbols)
	if err != nil {
		t.Fatalf("failed to get stock prices: %v", err)
	}

	if len(prices) == 0 {
		t.Fatalf("expected non-empty stock prices, got empty")
	}

	for code, info := range prices {
		displayName := info.Name
		if cfg != nil {
			if name, ok := cfg.TokenPriceMonitor.AStockNames[code]; ok && name != "" {
				displayName = name
			} else if name, ok := cfg.TokenPriceMonitor.HkStockNames[code]; ok && name != "" {
				displayName = name
			}
		}
		t.Logf("Stock Code: %s, DisplayName: %s (API Name: %s), Price: %.2f, Change: %.2f%%, UpdatedAt: %s",
			code, displayName, info.Name, info.Price, info.PercentChange, info.LastUpdated.Format("15:04:05"))
	}
}
