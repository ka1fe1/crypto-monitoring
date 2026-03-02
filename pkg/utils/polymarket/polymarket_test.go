package polymarket

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/ka1fe1/crypto-monitoring/config"
	"github.com/ka1fe1/crypto-monitoring/pkg/logger"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
)

var (
	cfg    *config.Config
	client *Client
)

// loadTestConfig resolves the absolute path to config.yaml and loads it.
func loadTestConfig() (*config.Config, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("failed to get current file path")
	}

	// Current file: /pkg/utils/polymarket/polymarket_test.go
	// Project root: /
	// pkg -> utils -> polymarket -> polymarket_test.go (3 levels up)
	rootDir := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(filename))))

	configPath := filepath.Join(rootDir, "config", "config.yaml")

	return config.LoadConfig(configPath)
}

func TestMain(m *testing.M) {
	var err error
	cfg, err = loadTestConfig()
	if err != nil {
		logger.Warn("Warning: Could not load config: %v", err)
	}

	var apiKey string
	if cfg != nil {
		apiKey = cfg.Polymarket.APIKey
	}

	client = NewClient(apiKey)
	os.Exit(m.Run())
}

func TestGetMarketDetail(t *testing.T) {
	// Use a known market ID (updated by user)
	marketID := "983678"
	market, err := client.GetMarketDetail(marketID)

	if err != nil {
		t.Fatalf("Failed to get market detail: %v", err)
	}

	if market == nil {
		t.Fatal("Market is nil")
	}

	if market.Question == "" {
		t.Error("Market question is empty")
	}

	if market.Slug == "" {
		t.Error("Market slug is empty")
	}

	if len(market.OutcomePrices) == 0 {
		t.Error("Outcome prices are empty")
	}

	t.Log(utils.PrintJson(market))
}

func TestGetTraderLeaderboardRankings(t *testing.T) {
	// Our client uses a hardcoded BaseURL, so we can't easily swap it with httptest.NewServer.
	// But since this is just a unit test and we already have a real API key, we can test it live with a known address
	// or we can just test if the endpoint parses correctly. Let's do a live test for a dummy address,
	// or just skip if it works. Here we'll do a live test for a common address or just check for no-error on zero volume.

	address := "0xfe9b5b2109a59138e4e645167a8384aa420dcb15"
	res, err := client.GetTraderLeaderboardRankings(address)
	if err != nil {
		t.Fatalf("Failed to GetTraderLeaderboardRankings: %v", err)
	}

	if res == nil {
		t.Fatal("LeaderboardResponse is nil")
	}

	t.Log(utils.PrintJson(res))

}

func TestGetTotalValueOfUserPositions(t *testing.T) {
	address := "0xfe9b5b2109a59138e4e645167a8384aa420dcb15"
	res, err := client.GetTotalValueOfUserPositions(address)
	if err != nil {
		t.Fatalf("Failed to GetTotalValueOfUserPositions: %v", err)
	}

	if res == nil {
		t.Fatal("TotalValueResponse is nil")
	}

	t.Log(utils.PrintJson(res))

	t.Logf("Total Value for %s: %v", address, res.Value)
}

func TestGetCurrentPositionsForUser(t *testing.T) {
	address := "0x11b6916fe7212b596e093d631d0d30ea80b1971d"
	res, err := client.GetCurrentPositionsForUser(address)
	if err != nil {
		t.Fatalf("Failed to GetCurrentPositionsForUser: %v", err)
	}

	if res == nil {
		t.Fatal("CurrentPositionsResponse is nil")
	}

	t.Logf("Current Positions for %s: %v items", address, len(*res))
	t.Log(utils.PrintJson(*res))
}

func TestResolveProxyWallet(t *testing.T) {
	// EOA address -> should resolve to proxyWallet
	address := "0xfbdcc3c6469b21c273517a971aa08400f272e514"
	proxyWallet, err := client.ResolveProxyWallet(address)
	if err != nil {
		t.Fatalf("Failed to ResolveProxyWallet: %v", err)
	}

	if proxyWallet == "" {
		t.Fatal("proxyWallet is empty")
	}

	t.Logf("EOA %s -> proxyWallet %s", address, proxyWallet)
}

func TestGetUserActivity(t *testing.T) {
	address := "0xfe9b5b2109a59138e4e645167a8384aa420dcb15"
	res, err := client.GetUserActivity(address)
	if err != nil {
		t.Fatalf("Failed to GetUserActivity: %v", err)
	}

	if res == nil {
		t.Fatal("ActivityResponse is nil")
	}

	t.Logf("Activity for %s: %v items", address, len(res))
	if len(res) > 0 {
		ts := res[0].Timestamp
		loc := time.FixedZone("UTC+8", 8*3600)
		formatted := time.Unix(ts, 0).In(loc).Format("2006-01-02 15:04:05")
		t.Logf("Latest Activity: timestamp=%v, type=%s, UTC+8=%s", res[0].Timestamp, res[0].Type, formatted)
	}
}
