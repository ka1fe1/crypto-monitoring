package service

import (
	"testing"

	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
)

func TestTwitterMonitorService_FetchNewTweets(t *testing.T) {
	if twitterSvc == nil {
		t.Skip("Service not initialized")
	}

	tweets, newestID, err := twitterSvc.FetchNewTweets("cz_binance", "2011840985085722926")
	if err != nil {
		t.Fatalf("FetchNewTweets failed: %v", err)
	}

	t.Log(newestID)
	t.Log(utils.PrintJson(tweets))
}
