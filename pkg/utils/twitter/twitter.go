package twitter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	BaseURL = "https://api.twitterapi.io/twitter/tweet/advanced_search"
)

// TwitterClient provides methods to interact with the twitterapi.io service.
type TwitterClient struct {
	APIKey string
	Client *http.Client
}

// NewTwitterClient creates a new TwitterClient.
func NewTwitterClient(apiKey string) *TwitterClient {
	return &TwitterClient{
		APIKey: apiKey,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Search performs an advanced search for tweets.
func (c *TwitterClient) Search(req AdvancedSearchRequest) (*SearchResponse, error) {
	apiUrl, err := url.Parse(BaseURL)
	if err != nil {
		return nil, err
	}

	q := apiUrl.Query()
	q.Set("query", req.Query)
	if req.Cursor != "" {
		q.Set("cursor", req.Cursor)
	}
	apiUrl.RawQuery = q.Encode()

	httpReq, err := http.NewRequest("GET", apiUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("X-API-Key", c.APIKey)

	resp, err := c.Client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var rawResponse RawTwitterSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return nil, err
	}

	return c.mapToSimplifiedResponse(rawResponse), nil
}

func (c *TwitterClient) mapToSimplifiedResponse(raw RawTwitterSearchResponse) *SearchResponse {
	tweets := make([]Tweet, 0, len(raw.Tweets))
	for _, rt := range raw.Tweets {
		// Parse CreatedAt: "Wed Jan 08 07:11:05 +0000 2025" or similar ISO string
		// The API docs say "<string>", usually Twitter API uses "Mon Jan 02 15:04:05 -0700 2006"
		// or ISO8601. We'll try common formats.
		createdAt, _ := time.Parse(time.RubyDate, rt.CreatedAt)
		if createdAt.IsZero() {
			createdAt, _ = time.Parse(time.RFC3339, rt.CreatedAt)
		}

		mentions := make([]string, 0, len(rt.Entities.UserMentions))
		for _, m := range rt.Entities.UserMentions {
			mentions = append(mentions, m.Name)
		}

		inReplyToUserName := ""
		if rt.IsReply && rt.InReplyToUserId != "" {
			// Find name in mentions
			for _, m := range rt.Entities.UserMentions {
				if m.IDStr == rt.InReplyToUserId {
					inReplyToUserName = m.Name
					break
				}
			}
			// Fallback to InReplyToUsername from API if not found in mentions
			if inReplyToUserName == "" {
				inReplyToUserName = rt.InReplyToUsername
			}
		}

		tweets = append(tweets, Tweet{
			Type:              rt.Type,
			ID:                rt.ID,
			URL:               rt.URL,
			Text:              rt.Text,
			CreatedAt:         createdAt,
			IsReply:           rt.IsReply,
			InReplyToUserName: inReplyToUserName,
			AuthorID:          rt.Author.ID,
			AuthorName:        rt.Author.Name,
			AuthorHandle:      rt.Author.UserName,
		})
	}

	return &SearchResponse{
		Tweets:      tweets,
		NextCursor:  raw.NextCursor,
		HasNextPage: raw.HasNextPage,
	}
}
