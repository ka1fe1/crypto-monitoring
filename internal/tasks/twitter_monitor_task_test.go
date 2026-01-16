package tasks

import (
	"testing"

	"github.com/ka1fe1/crypto-monitoring/internal/service"
	"github.com/ka1fe1/crypto-monitoring/pkg/logger"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/alter/dingding"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/constant"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/twitter"
)

func TestTwitterMonitorTask_Run(t *testing.T) {
	if cfg == nil {
		t.Skip("Config not loaded, skipping test")
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
	twitterService := service.NewTwitterMonitorService(client)
	usernames := cfg.TwitterMonitor.Usernames
	if len(usernames) == 0 {
		// Use a default one for testing if not configured
		t.Error("No Twitter usernames configured")
		return
	}

	qh := utils.QuietHoursParams{Enabled: true, StartHour: 11, EndHour: 12, Behavior: constant.QUIET_HOURS_BEHAVIOR_PAUSE}
	task := NewTwitterMonitorTask(twitterService, bot, usernames, cfg.TwitterMonitor.IntervalSeconds, qh)

	// Manually trigger run to test logic and notification
	// First run initializes the lastTweetIDs map
	logger.Info("Running first time to initialize...")
	task.run()

	// time.Sleep(60 * time.Second)

	// If you want to test the multi-run logic, you can modify the map or wait for new tweets
	// log.Printf("Running second time to check for new tweets...")
	// task.run()
}
