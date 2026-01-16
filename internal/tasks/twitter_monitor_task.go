package tasks

import (
	"fmt"
	"sync"
	"time"

	"github.com/ka1fe1/crypto-monitoring/internal/service"
	"github.com/ka1fe1/crypto-monitoring/pkg/logger"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/alter/dingding"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/twitter"
)

type TwitterMonitorTask struct {
	twitterService   service.TwitterMonitorService
	dingBot          *dingding.DingBot
	ticker           *time.Ticker
	stop             chan bool
	usernames        []string
	interval         time.Duration
	lastTweetIDs     map[string]string
	lastTweetLock    sync.RWMutex
	quietHoursParams utils.QuietHoursParams
	lastRunTime      time.Time
}

func NewTwitterMonitorTask(twitterService service.TwitterMonitorService, dingBot *dingding.DingBot, usernames []string, intervalSeconds int, quietHoursParams utils.QuietHoursParams) *TwitterMonitorTask {
	interval := time.Duration(intervalSeconds) * time.Second
	if interval <= 0 {
		interval = 600 * time.Second // Default 10 minutes
	}

	return &TwitterMonitorTask{
		twitterService:   twitterService,
		dingBot:          dingBot,
		stop:             make(chan bool),
		usernames:        usernames,
		interval:         interval,
		lastTweetIDs:     make(map[string]string),
		quietHoursParams: quietHoursParams,
	}

}

func (t *TwitterMonitorTask) Start() {
	t.ticker = time.NewTicker(t.interval)
	logger.Info("Starting Twitter Monitor Task with interval %v, monitoring %d accounts", t.interval, len(t.usernames))

	// Run immediately on start
	go t.run()

	go func() {
		for {
			select {
			case <-t.ticker.C:
				t.run()
			case <-t.stop:
				t.ticker.Stop()
				return
			}
		}
	}()
}

func (t *TwitterMonitorTask) Stop() {
	t.stop <- true
}

func (t *TwitterMonitorTask) run() {
	if len(t.usernames) == 0 {
		return
	}

	// Check for Quiet Hours
	if !utils.ShouldExecTask(t.quietHoursParams, t.lastRunTime, t.interval) {
		logger.Info("Skipping Twitter Monitor Task for %s in quiet hours", t.dingBot.Keyword)
		return
	}
	t.lastRunTime = time.Now()

	for _, username := range t.usernames {
		t.monitorUser(username)
	}
}

func (t *TwitterMonitorTask) monitorUser(username string) {
	t.lastTweetLock.RLock()
	lastID := t.lastTweetIDs[username]
	t.lastTweetLock.RUnlock()

	newTweets, newestID, err := t.twitterService.FetchNewTweets(username, lastID)
	if err != nil {
		logger.Error("Error searching tweets for %s: %v", username, err)
		return
	}

	if len(newTweets) > 0 {
		t.notifyTweets(username, newTweets)

		t.lastTweetLock.Lock()
		t.lastTweetIDs[username] = newestID
		t.lastTweetLock.Unlock()
	}
}

func (t *TwitterMonitorTask) notifyTweets(username string, tweets []twitter.Tweet) {
	title := fmt.Sprintf("%s [%s] New Tweets", t.dingBot.Keyword, username)

	content := t.formatTweets(tweets)

	allTexts := fmt.Sprintf("## %s\n\n --- \n\n%s\n\n---\n**Last Updated**: %s",
		title,
		content,
		utils.FormatBJTime(time.Now()),
	)

	err := t.dingBot.SendMarkdown(title, allTexts, nil, false)
	if err != nil {
		logger.Error("Error sending DingTalk notification for %s: %v", username, err)
	} else {
		logger.Info("Notified %d new tweets for %s", len(tweets), username)
	}
}

func (t *TwitterMonitorTask) formatTweets(tweets []twitter.Tweet) string {
	var content string

	// Display newest first (presumed sorted)
	for i := 0; i < len(tweets); i++ {
		tweet := tweets[i]
		isReplyStr := "No"
		if tweet.IsReply {
			isReplyStr = "Yes"
		}
		content += fmt.Sprintf("- Type: %s | IsReply: %s\n", tweet.Type, isReplyStr)
		content += fmt.Sprintf("- %s\n", tweet.Text)
		if tweet.InReplyToUserName != "" {
			content += fmt.Sprintf("- InReplyTo: %s\n", tweet.InReplyToUserName)
		}
		content += fmt.Sprintf("- [View on Twitter](%s)\n", tweet.URL)
		content += fmt.Sprintf("- %s\n", utils.FormatRelativeTime(tweet.CreatedAt))

		if i < len(tweets)-1 {
			content += "--- \n\n"
		}
	}
	return content
}
