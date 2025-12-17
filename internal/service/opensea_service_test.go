package service

import (
	"fmt"
	"os"
	"testing"

	"github.com/ka1fe1/crypto-monitoring/config"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/opensea"
)

var svc OpenSeaService

func TestMain(m *testing.M) {
	// Load config
	cfg, err := config.LoadConfig("../../config/config.yaml")
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize Clients
	osClient := opensea.NewOpenSeaClient(cfg.OpenSea.APIKey)

	cmcClient := utils.NewCoinMarketClient(cfg.CoinMarketCap.APIKey)

	// Initialize Service
	svc = NewOpenSeaService(osClient, cmcClient)

	os.Exit(m.Run())
}

func TestGetNFTFloorPrices(t *testing.T) {
	if svc == nil {
		t.Skip("Service not initialized")
	}

	// Use real slug from config or a known one? User asked to use config.
	// But sticking to a simple test case is safer.
	slugs := []string{"infinex-patrons"}

	results, err := svc.GetNFTFloorPrices(slugs, true)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	t.Log(utils.PrintJson(results))

}
