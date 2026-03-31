package binance

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ka1fe1/crypto-monitoring/config"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
)

var (
	cfg    *config.Config
	client *Client
)

// loadTestConfig resolves the absolute path to config.yaml and loads it.
func loadTestConfig() (*config.Config, error) {
	// 1. Get the absolute path of the current file to determine the project root.
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("failed to get current file path")
	}

	// The current file is in <ProjectRoot>/pkg/utils/binance/binance_test.go
	rootDir := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(filename))))

	// 2. Construct the absolute path to config.yaml
	configPath := filepath.Join(rootDir, "config", "config.yaml")

	// Fallback to config.yaml.temp if config.yaml does not exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = filepath.Join(rootDir, "config", "config.yaml.temp")
	}

	// 3. Load the configuration
	return config.LoadConfig(configPath)
}

func TestMain(m *testing.M) {
	var err error
	cfg, err = loadTestConfig()
	if err != nil {
		log.Printf("Warning: Could not load config: %v", err)
	}

	apiUrl := "https://api.binance.com"
	if cfg != nil && cfg.BtcDashboardMonitor.BinanceApiUrl != "" {
		apiUrl = cfg.BtcDashboardMonitor.BinanceApiUrl
	}

	client = NewClient(apiUrl)
	os.Exit(m.Run())
}

func TestGetKlines(t *testing.T) {
	// Test fetching BTCUSDT daily klines, limit 5
	klines, err := client.GetKlines("BTCUSDT", "1d", 5)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(klines) != 5 {
		t.Fatalf("expected 5 klines, got %d", len(klines))
	} else {
		t.Log(utils.PrintJson(klines))
	}

	for i, k := range klines {
		if k.OpenTime == 0 {
			t.Errorf("kline[%d] missing open time", i)
		}
		if k.Close == 0 {
			t.Errorf("kline[%d] missing close price", i)
		}
		if k.Volume == 0 {
			t.Errorf("kline[%d] missing volume", i)
		}
	}
}
