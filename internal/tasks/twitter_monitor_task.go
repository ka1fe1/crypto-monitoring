package tasks

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/alter/dingding"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/twitter"
)

type TwitterMonitorTask struct {
	client           *twitter.TwitterClient
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

func NewTwitterMonitorTask(client *twitter.TwitterClient, dingBot *dingding.DingBot, usernames []string, intervalSeconds int, quietHoursParams utils.QuietHoursParams) *TwitterMonitorTask {
	interval := time.Duration(intervalSeconds) * time.Second
	if interval <= 0 {
		interval = 600 * time.Second // Default 10 minutes
	}

	return &TwitterMonitorTask{
		client:           client,
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
	log.Printf("Starting Twitter Monitor Task with interval %v, monitoring %d accounts", t.interval, len(t.usernames))

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
		log.Printf("Skipping Twitter Monitor Task for %s in quiet hours", t.dingBot.Keyword)
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

	query := fmt.Sprintf("from:%s within_time:%s", username, "2h")
	if lastID != "" {
		query = fmt.Sprintf("from:%s since_id:%s", username, lastID)
	}

	req := twitter.AdvancedSearchRequest{
		Query: query,
	}

	resp, err := t.client.Search(req)
	if err != nil {
		log.Printf("Error searching tweets, %s: %v", t.formatQueryForLog(req.Query, lastID), err)
		return
	}

	if len(resp.Tweets) == 0 {
		log.Printf("No new tweets, %s", t.formatQueryForLog(req.Query, lastID))
		return
	}

	// Filter tweets to only those from the specified user
	var filteredTweets []twitter.Tweet
	for _, t := range resp.Tweets {
		if strings.EqualFold(t.AuthorHandle, username) {
			filteredTweets = append(filteredTweets, t)
		}
	}

	if len(filteredTweets) == 0 {
		log.Printf("No tweets, query: %s found in results (filtered out mentions)", t.formatQueryForLog(req.Query, lastID))
		return
	} else {
		log.Printf("Found %d tweets, query: %s", len(filteredTweets), t.formatQueryForLog(req.Query, lastID))
	}

	// Sort filtered tweets by ID descending (Newest first) to ensure correct processing
	// Twitter IDs are snowflake (time-ordered), so string comparison works if lengths are consistent.
	// We'll trust the string comparison for standard Twitter IDs.
	sort.Slice(filteredTweets, func(i, j int) bool {
		return filteredTweets[i].ID > filteredTweets[j].ID
	})

	// Update newestID to the newest tweet from the USER
	newestID := filteredTweets[0].ID

	// Notify for new tweets
	var newTweets []twitter.Tweet
	for _, tweet := range filteredTweets {
		// Only collect tweets strictly newer than lastID
		if lastID == "" || tweet.ID > lastID {
			newTweets = append(newTweets, tweet)
		}
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
	var content string

	// Display newest first
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

	allTexts := fmt.Sprintf("## %s\n\n --- \n\n%s\n\n---\n**Last Updated**: %s",
		title,
		content,
		utils.FormatBJTime(time.Now()),
	)

	err := t.dingBot.SendMarkdown(title, allTexts, nil, false)
	if err != nil {
		log.Printf("Error sending DingTalk notification for %s: %v", username, err)
	} else {
		log.Printf("Notified %d new tweets for %s", len(tweets), username)
	}
}

func (t *TwitterMonitorTask) formatQueryForLog(query string, lastID string) string {
	res := fmt.Sprintf("query: %s", query)
	if lastID != "" {
		if st, err := utils.SnowflakeToTime(lastID); err == nil {
			res += fmt.Sprintf(" (since %s)", utils.FormatBJTime(st))
		}
	}
	return res
}
