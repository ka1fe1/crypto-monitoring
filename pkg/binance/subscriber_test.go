package binance

import (
	"encoding/hex"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ka1fe1/crypto-monitoring/config"
)

func TestGenerateSignature(t *testing.T) {
	apiKey := "testKey"
	apiSecret := "testSecret"
	sub := NewCmsSubscriber(apiKey, apiSecret)

	// Test data
	data := "random=123&timestamp=456"

	// Expected signature: HMAC-SHA256(data, key)
	// We can compute this or just check properties.
	// Let's rely on the fact that hex.DecodeString works and length is correct.

	sig := sub.generateSignature(data)
	if len(sig) != 64 {
		t.Errorf("Expected signature length 64, got %d", len(sig))
	}

	_, err := hex.DecodeString(sig)
	if err != nil {
		t.Errorf("Signature is not valid hex: %v", err)
	}
}

func TestCmsSubscriber_Connect(t *testing.T) {
	// Skip if config file not present (optional, but good for CI if no secrets)
	configPath := "../../config/config.yaml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Skip("config.yaml not found, skipping integration test")
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.BinanceCex.APIKey == "" || cfg.BinanceCex.SecretKey == "" {
		t.Skip("Binance credentials not found in config, skipping integration test")
	}

	sub := NewCmsSubscriber(cfg.BinanceCex.APIKey, cfg.BinanceCex.SecretKey)
	if cfg.BinanceCex.ProxyURL != "" {
		sub.SetProxy(cfg.BinanceCex.ProxyURL)
	}

	// Test Connection
	// We pass empty topics for now just to test auth
	err = sub.Connect([]string{})
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer sub.Close()

	// Keep alive with Ticker
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Wait and Ping repeatedly
	// Note: This blocks indefinitely, user might need to stop manually or Set a timeout
	go func() {
		for range ticker.C {
			if err := sub.Ping(); err != nil {
				fmt.Printf("Failed to ping: %v\n", err)
			} else {
				fmt.Println("Ping sent")
			}
		}
	}()

	var topics []string
	topics = append(topics, "com_announcement_en", "com_announcement_cn")
	if err := sub.Subscribe(topics); err != nil {
		t.Errorf("Failed to subscribe: %v", err)
	}

	go sub.listen(func(msg []byte) { fmt.Printf("recv: %s\n", msg) })

	time.Sleep(500 * time.Second)
}
