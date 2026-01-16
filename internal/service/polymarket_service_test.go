package service

import (
	"testing"

	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
)

func TestPolymarketMonitorService_GetMarketDetails(t *testing.T) {
	if polySvc == nil {
		t.Skip("Service not initialized")
	}

	markets, err := polySvc.GetMarketDetails([]string{"983678", "763535"})
	if err != nil {
		t.Fatalf("GetMarketDetails failed: %v", err)
	}
	t.Log(utils.PrintJson(markets))
}
