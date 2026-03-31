package alternative

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"strconv"

	"github.com/ka1fe1/crypto-monitoring/config"
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

	// The current file is in <ProjectRoot>/pkg/utils/alternative/alternative_test.go
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

	apiUrl := "https://api.alternative.me"
	if cfg != nil && cfg.BtcDashboardMonitor.AlternativeApiUrl != "" {
		apiUrl = cfg.BtcDashboardMonitor.AlternativeApiUrl
	}

	client = NewClient(apiUrl)
	os.Exit(m.Run())
}

func TestGetFng(t *testing.T) {
	fng, err := client.GetFng(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(fng.Data) != 1 {
		t.Fatalf("expected 1 record, got %d", len(fng.Data))
	}

	valueInt, err := strconv.Atoi(fng.Data[0].Value)
	if err != nil {
		t.Fatalf("expected valid integer in value, got %v", err)
	}
	
	if valueInt < 0 || valueInt > 100 {
		t.Errorf("fgi should be between 0 and 100, got %d", valueInt)
	}
	t.Logf("Fear & Greed Index: %d, Classification: %s", valueInt, fng.Data[0].ValueClassification)
}
