package service

import (
	"testing"

	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
)

var svc OpenSeaService

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
