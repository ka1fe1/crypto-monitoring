package binance

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

const (
	BaseURL = "wss://api.binance.com/sapi/wss"
)

// CmsSubscriber handles WebSocket subscriptions for Binance announcements.
type CmsSubscriber struct {
	apiKey    string
	apiSecret string
	conn      *websocket.Conn
	done      chan struct{}
	proxyURL  string
}

// NewCmsSubscriber creates a new CmsSubscriber.
func NewCmsSubscriber(apiKey, apiSecret string) *CmsSubscriber {
	return &CmsSubscriber{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		done:      make(chan struct{}),
	}
}

// SetProxy sets the proxy URL for the subscriber.
func (s *CmsSubscriber) SetProxy(proxyURL string) {
	s.proxyURL = proxyURL
}

// Connect establishes the WebSocket connection.
// If topics are provided, they are included in the initial connection URL.
func (s *CmsSubscriber) Connect(topics []string) error {
	u, err := url.Parse(BaseURL)
	if err != nil {
		return err
	}

	// Build query parameters
	values := url.Values{}
	values.Set("random", generateRandomString(32))
	if len(topics) > 0 {
		topicStr := ""
		for i, t := range topics {
			if i > 0 {
				topicStr += "|"
			}
			topicStr += t
		}
		values.Set("topic", topicStr)
	}
	values.Set("recvWindow", "50000")
	values.Set("timestamp", fmt.Sprintf("%d", time.Now().UnixNano()/1e6))

	// Generate signature from the encoded query string
	queryString := values.Encode()
	signature := s.generateSignature(queryString)
	values.Set("signature", signature)

	u.RawQuery = values.Encode()

	log.Printf("Connecting to %s", u.String())

	// Set API Key in headers
	headers := http.Header{}
	headers.Add("X-MBX-APIKEY", s.apiKey)

	dialer := websocket.DefaultDialer
	if s.proxyURL != "" {
		proxy, err := url.Parse(s.proxyURL)
		if err != nil {
			return err
		}
		dialer = &websocket.Dialer{
			Proxy:            http.ProxyURL(proxy),
			HandshakeTimeout: 45 * time.Second,
		}
	}

	c, _, err := dialer.Dial(u.String(), headers)
	if err != nil {
		return err
	}
	s.conn = c

	return nil
}

func (s *CmsSubscriber) generateSignature(data string) string {
	h := hmac.New(sha256.New, []byte(s.apiSecret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func generateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (s *CmsSubscriber) listen(handler func([]byte)) {
	defer close(s.done)
	for {
		_, message, err := s.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		handler(message)
	}
}

// Subscribe sends a subscription message for topics.
func (s *CmsSubscriber) Subscribe(topics []string) error {
	if len(topics) == 0 {
		return nil
	}
	topicStr := ""
	for i, t := range topics {
		if i > 0 {
			topicStr += "|"
		}
		topicStr += t
	}

	msg := map[string]string{
		"command": "SUBSCRIBE",
		"value":   topicStr,
	}
	return s.conn.WriteJSON(msg)
}

// Unsubscribe sends an unsubscription message.
func (s *CmsSubscriber) Unsubscribe(topics []string) error {
	if len(topics) == 0 {
		return nil
	}
	topicStr := ""
	for i, t := range topics {
		if i > 0 {
			topicStr += "|"
		}
		topicStr += t
	}

	msg := map[string]string{
		"command": "UNSUBSCRIBE",
		"value":   topicStr,
	}
	return s.conn.WriteJSON(msg)
}

// Ping sends a ping message to the server.
func (s *CmsSubscriber) Ping() error {
	return s.conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second))
}

func (s *CmsSubscriber) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

// Run maintains the connection and handles messages.
func (s *CmsSubscriber) Run(ctx context.Context, topics []string, handler func([]byte)) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := s.Connect(topics); err != nil {
				log.Printf("Failed to connect: %v, retrying in 5s...", err)
				time.Sleep(5 * time.Second)
				continue
			}

			// Start Pinger
			ticker := time.NewTicker(30 * time.Second)
			go func() {
				defer ticker.Stop()
				for {
					select {
					case <-ctx.Done():
						return
					case <-s.done:
						return
					case <-ticker.C:
						if err := s.Ping(); err != nil {
							log.Printf("Ping failed: %v", err)
							s.Close()
							return
						}
					}
				}
			}()

			log.Println("Connected to Binance CEX")
			s.listen(handler)
			log.Println("Connection lost, reconnecting...")
			s.Close()
			time.Sleep(1 * time.Second)
		}
	}
}
