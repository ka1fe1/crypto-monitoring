package dingding

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	BaseURL = "https://oapi.dingtalk.com/robot/send"
)

type DingBot struct {
	Token   string
	Secret  string
	Keyword string
	BaseURL string
	Client  *http.Client
}

func NewDingBot(token, secret, keyword string) *DingBot {
	return &DingBot{
		Token:   token,
		Secret:  secret,
		Keyword: keyword,
		BaseURL: BaseURL,
		Client:  &http.Client{Timeout: 10 * time.Second},
	}
}

type At struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

type TextMessage struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
	At At `json:"at"`
}

type MarkdownMessage struct {
	MsgType  string `json:"msgtype"`
	Markdown struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"markdown"`
	At At `json:"at"`
}

type Response struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (bot *DingBot) sign(t int64) string {
	if bot.Secret == "" {
		return ""
	}
	stringToSign := fmt.Sprintf("%d\n%s", t, bot.Secret)
	h := hmac.New(sha256.New, []byte(bot.Secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (bot *DingBot) send(msg interface{}) error {
	u, err := url.Parse(bot.BaseURL)
	if err != nil {
		return fmt.Errorf("failed to parse base url: %w", err)
	}

	q := u.Query()
	q.Set("access_token", bot.Token)

	if bot.Secret != "" {
		timestamp := time.Now().UnixMilli()
		sign := bot.sign(timestamp)
		q.Set("timestamp", fmt.Sprintf("%d", timestamp))
		q.Set("sign", sign)
	}

	u.RawQuery = q.Encode()

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	resp, err := bot.Client.Post(u.String(), "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var dingResp Response
	if err := json.Unmarshal(respBody, &dingResp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if dingResp.ErrCode != 0 {
		return fmt.Errorf("dingtalk api error: %s (code: %d)", dingResp.ErrMsg, dingResp.ErrCode)
	}

	return nil
}

func (bot *DingBot) SendText(content string, atMobiles []string, isAtAll bool) error {
	if bot.Keyword != "" {
		content = fmt.Sprintf("[%s] %s", bot.Keyword, content)
	}
	msg := TextMessage{
		MsgType: "text",
		At: At{
			AtMobiles: atMobiles,
			IsAtAll:   isAtAll,
		},
	}
	msg.Text.Content = content
	return bot.send(msg)
}

func (bot *DingBot) SendMarkdown(title, text string, atMobiles []string, isAtAll bool) error {
	if bot.Keyword != "" {
		if !strings.Contains(text, bot.Keyword) {
			text = fmt.Sprintf("[%s]\n%s", bot.Keyword, text)
		}
	}

	msg := MarkdownMessage{
		MsgType: "markdown",
		At: At{
			AtMobiles: atMobiles,
			IsAtAll:   isAtAll,
		},
	}
	msg.Markdown.Title = title
	msg.Markdown.Text = text
	return bot.send(msg)
}
