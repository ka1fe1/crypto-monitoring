package service

import (
	"testing"

	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
)

func TestTwitterService_FetchNewTweets(t *testing.T) {
	if twitterSvc == nil {
		t.Skip("Service not initialized")
	}

	tweets, newestID, err := twitterSvc.FetchNewTweets("bwenews", "", "",
		[]string{
			"Binance Alpha", "Binance", "UPBIT LISTING",
			"elonmusk", "trump"})
	if err != nil {
		t.Fatalf("FetchNewTweets failed: %v", err)
	}

	t.Log(newestID)
	t.Log(utils.PrintJson(tweets))
}
