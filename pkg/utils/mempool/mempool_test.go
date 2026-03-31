package mempool

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ka1fe1/crypto-monitoring/config"
)

var (
	cfg    *config.Config
	client *Client
)

func loadTestConfig() (*config.Config, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("failed to get current file path")
	}

	rootDir := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(filename))))
	configPath := filepath.Join(rootDir, "config", "config.yaml")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = filepath.Join(rootDir, "config", "config.yaml.temp")
	}

	return config.LoadConfig(configPath)
}

func TestMain(m *testing.M) {
	var err error
	cfg, err = loadTestConfig()
	if err != nil {
		log.Printf("Warning: Could not load config: %v", err)
	}

	apiUrl := "https://mempool.space"
	if cfg != nil && cfg.BtcDashboardMonitor.MempoolApiUrl != "" {
		apiUrl = cfg.BtcDashboardMonitor.MempoolApiUrl
	}

	client = NewClient(apiUrl)
	os.Exit(m.Run())
}

func TestGetTipHeight(t *testing.T) {
	height, err := client.GetTipHeight()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if height < 800000 {
		t.Fatalf("expected height > 800000, got %d", height)
	}
	t.Logf("current tip height: %d", height)
}
