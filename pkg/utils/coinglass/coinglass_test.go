package coinglass

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
)

func TestGetAHR999(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))
	defer server.Close()

	// Override BaseURL for testing
	// Note: in a real scenario, we might want to make BaseURL injectable
	// For this test, we'll manually construct the client and adjust the request if needed,
	// or use a more flexible design. Since BaseURL is a constant, I'll update the code to allow setting it or just mock the request logic.

	// Let's modify coinglass.go slightly to allow base URL override if we want to be clean,
	// or just test the Unmarshal logic separately.

	client := NewClient("87068487ab364af689cc6b958ce1ca0e")

	if res, err := client.GetAHR999(); err != nil {
		t.Fatalf("Failed to get AHR999: %v", err)
	} else {
		t.Logf("AHR999: %v", utils.PrintJson(res))
	}
}

func TestUnmarshalFearGreed(t *testing.T) {
	data := []byte(`{
		"code": "0",
		"msg": "success",
		"data": [
			{
				"values": [45.0],
				"price": [42000.0],
				"time_list": [1704067200]
			}
		]
	}`)

	var res Response[[]FearGreedData]
	err := json.Unmarshal(data, &res)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if res.Code != "0" {
		t.Errorf("Expected code 0, got %s", res.Code)
	}

	if len(res.Data) != 1 {
		t.Errorf("Expected 1 data entry, got %d", len(res.Data))
	}

	if res.Data[0].Values[0] != 45.0 {
		t.Errorf("Expected 45.0, got %f", res.Data[0].Values[0])
	}
}
