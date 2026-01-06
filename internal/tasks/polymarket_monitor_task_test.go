package tasks

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ka1fe1/crypto-monitoring/config"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/alter/dingding"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/constant"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/polymarket"
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

	// Current file is in <ProjectRoot>/internal/tasks/polymarket_monitor_task_test.go
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

	os.Exit(m.Run())
}

func TestPolymarketMonitorTask_Run(t *testing.T) {
	if cfg == nil {
		t.Skip("Config not loaded, skipping test")
	}

	botName := cfg.PolymarketMonitor.BotName
	if botName == "" {
		botName = constant.DEFAULT_BOT_NAME
	}

	botCfg, ok := cfg.DingTalk[botName]
	if !ok {
		t.Skipf("DingTalk bot %s not configured, skipping test", botName)
	}

	bot := dingding.NewDingBot(botCfg.AccessToken, botCfg.Secret, botCfg.Keyword)
	client := polymarket.NewClient(cfg.Polymarket.APIKey)
	marketIDs := cfg.PolymarketMonitor.MarketIDs
	if len(marketIDs) == 0 {
		// Use a default one for testing if not configured
		marketIDs = []string{"983678"}
	}

	task := NewPolymarketMonitorTask(client, bot, marketIDs, cfg.PolymarketMonitor.IntervalSeconds)

	// Manually trigger run to test logic and notification
	task.run()
}
