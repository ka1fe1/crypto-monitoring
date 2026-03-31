package bgeometrics

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/ka1fe1/crypto-monitoring/config"
)

var cfg *config.Config
var bCli *Client

func loadTestConfig() (*config.Config, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("failed to get current file path")
	}

	rootDir := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(filename))))
	configPath := filepath.Join(rootDir, "config", "config.yaml")

	return config.LoadConfig(configPath)
}

func TestMain(m *testing.M) {
	var err error
	cfg, err = loadTestConfig()
	if err != nil {
		log.Printf("Warning: Could not load config: %v", err)
	}

	var apiKey, apiUrl string
	if cfg != nil {
		apiKey = cfg.BtcDashboardMonitor.BgeometricsApiKey
		apiUrl = cfg.BtcDashboardMonitor.BgeometricsApiUrl
	}

	bCli = NewClient(apiUrl, apiKey)
	os.Exit(m.Run())
}

// TestGetBalancedPrice_Real 使用真实 API 验证端到端连通性
func TestGetBalancedPrice_Real(t *testing.T) {
	if cfg == nil || cfg.BtcDashboardMonitor.BgeometricsApiKey == "" {
		t.Skip("Skipping real API test: BGeometrics API Key not configured.")
	}

	price, err := bCli.GetBalancedPrice()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if price <= 0 {
		t.Errorf("Expected balanced price to be greater than 0, got %f", price)
	}
	t.Logf("Fetched Balanced Price: %.2f", price)
}

// TestGetBalancedPrice_MockSingleObject 使用 mock server 验证单对象响应的解析逻辑
func TestGetBalancedPrice_MockSingleObject(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/balanced-price/1" {
			t.Errorf("Expected path /v1/balanced-price/1, got %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("Expected Authorization header, got %s", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"d":"2026-03-30","unixTs":1774934400,"balancedPrice":42150.75}`))
	}))
	defer server.Close()

	cli := NewClient(server.URL, "test-key", WithTimeout(5*time.Second))
	price, err := cli.GetBalancedPrice()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if price != 42150.75 {
		t.Errorf("Expected 42150.75, got %f", price)
	}
	t.Logf("Mock single object parsed price: %.2f", price)
}

// TestGetBalancedPrice_MockArray 使用 mock server 验证数组响应的 fallback 解析
func TestGetBalancedPrice_MockArray(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"d":"2026-03-29","unixTs":1774848000,"balancedPrice":41000.00},{"d":"2026-03-30","unixTs":1774934400,"balancedPrice":42500.50}]`))
	}))
	defer server.Close()

	cli := NewClient(server.URL, "test-key")
	price, err := cli.GetBalancedPrice()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if price != 42500.50 {
		t.Errorf("Expected 42500.50 (last element), got %f", price)
	}
	t.Logf("Mock array parsed price: %.2f", price)
}

// TestGetBalancedPrice_EmptyAPIKey 验证 API Key 为空时提前返回错误
func TestGetBalancedPrice_EmptyAPIKey(t *testing.T) {
	cli := NewClient("https://bitcoin-data.com", "")
	_, err := cli.GetBalancedPrice()
	if err == nil {
		t.Fatal("Expected error for empty API key, got nil")
	}
	t.Logf("Correctly returned error for empty API key: %v", err)
}

// TestGetBalancedPrice_MockServerError 验证非 200 状态码的错误处理
func TestGetBalancedPrice_MockServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"invalid token"}`))
	}))
	defer server.Close()

	cli := NewClient(server.URL, "bad-key")
	_, err := cli.GetBalancedPrice()
	if err == nil {
		t.Fatal("Expected error for 401 response, got nil")
	}
	t.Logf("Correctly returned error for 401: %v", err)
}
