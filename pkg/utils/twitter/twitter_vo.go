package twitter

import "time"

// AdvancedSearchRequest represents the request parameters for the Twitter advanced search API.
type AdvancedSearchRequest struct {
	Query  string `json:"query"`
	Cursor string `json:"cursor,omitempty"`
}

// RawTwitterSearchResponse represents the raw response structure from the twitterapi.io endpoint.
type RawTwitterSearchResponse struct {
	Tweets      []RawTweet `json:"tweets"`
	HasNextPage bool       `json:"has_next_page"`
	NextCursor  string     `json:"next_cursor"`
}

// RawTweet represents the raw tweet structure from the API.
type RawTweet struct {
	Type              string      `json:"type"`
	ID                string      `json:"id"`
	URL               string      `json:"url"`
	Text              string      `json:"text"`
	Source            string      `json:"source"`
	RetweetCount      int         `json:"retweetCount"`
	ReplyCount        int         `json:"replyCount"`
	LikeCount         int         `json:"likeCount"`
	QuoteCount        int         `json:"quoteCount"`
	ViewCount         int         `json:"viewCount"`
	CreatedAt         string      `json:"createdAt"`
	Lang              string      `json:"lang"`
	BookmarkCount     int         `json:"bookmarkCount"`
	IsReply           bool        `json:"isReply"`
	InReplyToId       string      `json:"inReplyToId"`
	ConversationId    string      `json:"conversationId"`
	InReplyToUserId   string      `json:"inReplyToUserId"`
	InReplyToUsername string      `json:"inReplyToUsername"`
	Author            RawAuthor   `json:"author"`
	Entities          RawEntities `json:"entities"`
}

// RawAuthor represents the raw author structure from the API.
type RawAuthor struct {
	ID       string `json:"id"`
	UserName string `json:"userName"`
	Name     string `json:"name"`
	URL      string `json:"url"`
}

// RawEntities represents the raw entities structure from the API.
type RawEntities struct {
	Hashtags     []RawHashtag     `json:"hashtags"`
	Urls         []RawUrl         `json:"urls"`
	UserMentions []RawUserMention `json:"user_mentions"`
}

// RawHashtag represents a raw hashtag.
type RawHashtag struct {
	Text string `json:"text"`
}

// RawUrl represents a raw URL.
type RawUrl struct {
	URL         string `json:"url"`
	ExpandedURL string `json:"expanded_url"`
}

// RawUserMention represents a raw user mention.
type RawUserMention struct {
	IDStr      string `json:"id_str"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

// Tweet represents the simplified tweet structure requested by the user.
type Tweet struct {
	Type              string    `json:"type"`
	ID                string    `json:"id"`
	URL               string    `json:"url"`
	Text              string    `json:"text"`
	CreatedAt         time.Time `json:"createdAt"`
	IsReply           bool      `json:"isReply"`
	InReplyToUserName string    `json:"inReplyToUserName"`
	AuthorID          string    `json:"authorId"`
	AuthorName        string    `json:"authorName"`
	AuthorHandle      string    `json:"authorHandle"`
}

// SearchResponse is the simplified response from our utility.
type SearchResponse struct {
	Tweets      []Tweet `json:"tweets"`
	NextCursor  string  `json:"next_cursor"`
	HasNextPage bool    `json:"has_next_page"`
}
