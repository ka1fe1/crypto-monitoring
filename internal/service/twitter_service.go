package service

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ka1fe1/crypto-monitoring/pkg/logger"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/twitter"
)

type TwitterService interface {
	FetchNewTweets(username string, withinTime string, lastID string, keywords []string) ([]twitter.Tweet, string, error)
}

type twitterService struct {
	client *twitter.TwitterClient
}

func NewTwitterService(client *twitter.TwitterClient) TwitterService {
	return &twitterService{
		client: client,
	}
}

func (s *twitterService) FetchNewTweets(username string, withinTime string, lastID string, keywords []string) ([]twitter.Tweet, string, error) {
	if withinTime == "" {
		withinTime = "2h"
	}
	query := fmt.Sprintf("from:%s within_time:%s", username, withinTime)
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

	// Step 1: Filter tweets to only those from the specified user
	var userTweets []twitter.Tweet
	for _, t := range resp.Tweets {
		if strings.EqualFold(t.AuthorHandle, username) {
			userTweets = append(userTweets, t)
		}
	}

	if len(userTweets) == 0 {
		logger.Info("No tweets found, %s", logQuery)
		return nil, "", nil
	} else {
		logger.Info("Found %d tweets, %s", len(userTweets), logQuery)
		logger.Debug("userTweets: %s", utils.PrintJson(userTweets))
	}

	// Step 2: Sort all user tweets by ID descending (Newest first) to determine newestID
	sort.Slice(userTweets, func(i, j int) bool {
		return userTweets[i].ID > userTweets[j].ID
	})

	// Newest ID is the first one after sorting (User's latest tweet, regardless of keywords)
	newestID := userTweets[0].ID

	// Step 3: Apply Keyword Filtering on pending tweets
	var filteredTweets []twitter.Tweet
	// userKeywords, hasKeywords := s.keywords[username]

	for _, t := range userTweets {
		// Only consider new tweets
		if lastID != "" && t.ID <= lastID {
			continue
		}

		if len(keywords) > 0 {
			match := false
			for _, kw := range keywords {
				if strings.Contains(strings.ToLower(t.Text), strings.ToLower(kw)) {
					match = true
					break
				}
			}
			if !match {
				truncatedText := t.Text
				if len([]rune(truncatedText)) > 50 {
					truncatedText = string([]rune(truncatedText)[:50]) + "..."
				}
				logger.Debug("Filtered out tweet: %s, username: %s, missing keywords: [%s]", truncatedText, username, strings.Join(keywords, ", "))
				continue
			}
		}

		filteredTweets = append(filteredTweets, t)
	}

	if len(filteredTweets) == 0 {
		logger.Info("No new tweets for query %s, keywords: [%s]", logQuery, strings.Join(keywords, ", "))
	} else {
		logger.Info("Found %d new tweets for query %s, keywords: [%s]", len(filteredTweets), logQuery, strings.Join(keywords, ", "))
	}

	return filteredTweets, newestID, nil
}
