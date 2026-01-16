package tasks

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ka1fe1/crypto-monitoring/config"
	"github.com/ka1fe1/crypto-monitoring/pkg/logger"
)

var (
	cfg *config.Config
)

// loadTestConfig resolves the absolute path to config.yaml and loads it.
func loadTestConfig() (*config.Config, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, nil
	}

	// Current file is in <ProjectRoot>/internal/tasks/main_test.go
	// Project root is two levels up.
	rootDir := filepath.Dir(filepath.Dir(filepath.Dir(filename)))

	configPath := filepath.Join(rootDir, "config", "config.yaml")

	return config.LoadConfig(configPath)
}

func TestMain(m *testing.M) {
	var err error
	cfg, err = loadTestConfig()
	if err != nil {
		log.Printf("Warning: Could not load config: %v", err)
	}

	// Setup Logger with DEBUG level for tests
	// If config bas log level, use it, otherwise default to "debug" for tests
	logLevel := "debug"
	if cfg != nil && cfg.Log.Level != "" {
		logLevel = cfg.Log.Level
	}
	logger.Setup(logLevel)

	os.Exit(m.Run())
}
