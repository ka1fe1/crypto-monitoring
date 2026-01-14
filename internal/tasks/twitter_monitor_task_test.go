package tasks

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ka1fe1/crypto-monitoring/config"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/alter/dingding"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/constant"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/twitter"
)

// loadTwitterTestConfig resolves the absolute path to config.yaml and loads it.
func loadTwitterTestConfig() (*config.Config, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("failed to get current file path")
	}

	// Current file is in <ProjectRoot>/internal/tasks/twitter_monitor_task_test.go
	// Project root is two levels up.
	rootDir := filepath.Dir(filepath.Dir(filepath.Dir(filename)))

	configPath := filepath.Join(rootDir, "config", "config.yaml")

	return config.LoadConfig(configPath)
}

func TestTwitterMonitorTask_Run(t *testing.T) {
	cfg, err := loadTwitterTestConfig()
	if err != nil {
		t.Skipf("Warning: Could not load config: %v", err)
	}

	botName := cfg.TwitterMonitor.BotName
	if botName == "" {
		botName = constant.DEFAULT_BOT_NAME
	}

	botCfg, ok := cfg.DingTalk[botName]
	if !ok {
		t.Skipf("DingTalk bot %s not configured, skipping test", botName)
	}

	bot := dingding.NewDingBot(botCfg.AccessToken, botCfg.Secret, botCfg.Keyword)
	client := twitter.NewTwitterClient(cfg.Twitter.APIKey)
	usernames := cfg.TwitterMonitor.Usernames
	if len(usernames) == 0 {
		// Use a default one for testing if not configured
		t.Error("No Twitter usernames configured")
		return
	}

	qh := utils.QuietHoursParams{Enabled: true, StartHour: 11, EndHour: 12, Behavior: constant.QUIET_HOURS_BEHAVIOR_PAUSE}
	task := NewTwitterMonitorTask(client, bot, usernames, cfg.TwitterMonitor.IntervalSeconds, qh)

	// Manually trigger run to test logic and notification
	// First run initializes the lastTweetIDs map
	log.Printf("Running first time to initialize...")
	task.run()

	// If you want to test the multi-run logic, you can modify the map or wait for new tweets
	// log.Printf("Running second time to check for new tweets...")
	// task.run()
}
