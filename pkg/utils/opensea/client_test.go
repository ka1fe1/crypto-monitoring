package opensea

import (
	"testing"

	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
)

func TestGetCollectionStats(t *testing.T) {
	// Create client
	client := NewOpenSeaClient("d8e8122681144d1ca37d53a5787bcbd5")

	// Test GetCollectionStats
	stats, err := client.GetCollectionStats("infinex-patrons")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	t.Log(utils.PrintJson(stats))
}
