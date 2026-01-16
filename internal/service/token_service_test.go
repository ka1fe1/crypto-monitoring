package service

import (
	"testing"

	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
)

func TestTokenService_GetTokenPrice(t *testing.T) {
	if tokenSvc == nil {
		t.Skip("Service not initialized")
	}

	prices, err := tokenSvc.GetTokenPrice([]string{"1", "1027", "1839", "5426"})
	if err != nil {
		t.Fatalf("GetTokenPrice failed: %v", err)
	}

	t.Log(utils.PrintJson(prices))
}
