package service

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ka1fe1/crypto-monitoring/pkg/logger"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/twitter"
)

type TwitterMonitorService interface {
	FetchNewTweets(username string, lastID string) ([]twitter.Tweet, string, error)
}

type twitterMonitorService struct {
	client *twitter.TwitterClient
}

func NewTwitterMonitorService(client *twitter.TwitterClient) TwitterMonitorService {
	return &twitterMonitorService{
		client: client,
	}
}

func (s *twitterMonitorService) FetchNewTweets(username string, lastID string) ([]twitter.Tweet, string, error) {
	query := fmt.Sprintf("from:%s within_time:%s", username, "2h")
	if lastID != "" {
		query = fmt.Sprintf("from:%s since_id:%s", username, lastID)
	}

	req := twitter.AdvancedSearchRequest{
		Query: query,
	}

	resp, err := s.client.Search(req)
	if err != nil {
		return nil, "", fmt.Errorf("search failed: %w", err)
	}

	// Prepare log query with readable time if since_id is present
	logQuery := query
	if lastID != "" {
		if t, err := utils.SnowflakeToTime(lastID); err == nil {
			logQuery = fmt.Sprintf("%s (since: %s)", query, utils.FormatBJTime(t))
		}
	}

	if len(resp.Tweets) == 0 {
		logger.Debug("No new tweets found, %s", logQuery)
		return nil, "", nil
	} else {
		logger.Debug("Found %d tweets, query: %s", len(resp.Tweets), logQuery)
	}

	// Filter tweets to only those from the specified user
	var filteredTweets []twitter.Tweet
	for _, t := range resp.Tweets {
		if strings.EqualFold(t.AuthorHandle, username) {
			filteredTweets = append(filteredTweets, t)
		}
	}

	if len(filteredTweets) == 0 {
		logger.Info("No tweets found for query %s after filtering", logQuery)
		return nil, "", nil
	} else {
		logger.Info("Found %d tweets for query %s after filtering", len(filteredTweets), logQuery)
	}

	// Sort filtered tweets by ID descending (Newest first)
	sort.Slice(filteredTweets, func(i, j int) bool {
		return filteredTweets[i].ID > filteredTweets[j].ID
	})

	// Newest ID is the first one after sorting
	newestID := filteredTweets[0].ID

	// Select only new tweets (strictly greater than lastID)
	var newTweets []twitter.Tweet
	for _, tweet := range filteredTweets {
		if lastID == "" || tweet.ID > lastID {
			newTweets = append(newTweets, tweet)
		}
	}

	return newTweets, newestID, nil
}
