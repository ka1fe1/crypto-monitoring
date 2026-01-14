package twitter

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ka1fe1/crypto-monitoring/config"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
)

var (
	cfg    *config.Config
	client *TwitterClient
)

// loadTestConfig resolves the absolute path to config.yaml and loads it.
func loadTestConfig() (*config.Config, error) {
	// 1. Get the absolute path of the current file to determine the project root.
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("failed to get current file path")
	}

	// The current file is in <ProjectRoot>/pkg/utils/twitter/twitter_test.go
	// So we go up three levels to get to <ProjectRoot>
	// filename: .../pkg/utils/twitter/twitter_test.go
	// Dir1: .../pkg/utils/twitter
	// Dir2: .../pkg/utils
	// Dir3: .../pkg
	// Dir4: ... (Project Root)
	rootDir := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(filename))))

	// 2. Construct the absolute path to config.yaml
	configPath := filepath.Join(rootDir, "config", "config.yaml")

	// 3. Load the configuration
	return config.LoadConfig(configPath)
}

func TestMain(m *testing.M) {
	var err error
	cfg, err = loadTestConfig()
	if err != nil {
		log.Printf("Warning: Could not load config: %v", err)
	}

	apiKey := ""
	if cfg != nil {
		apiKey = cfg.Twitter.APIKey
	}

	client = NewTwitterClient(apiKey)
	os.Exit(m.Run())
}

func TestTwitterClient_Search(t *testing.T) {
	if client.APIKey == "" {
		t.Skip("Twitter API key not configured, skipping integration test")
	}

	req := AdvancedSearchRequest{
		Query: "from:joejoedefi since_id:2010999018004922548",
	}

	resp, err := client.Search(req)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if resp == nil {
		t.Fatal("Response is nil")
	}

	fmt.Printf("Found %d tweets\n", len(resp.Tweets))

	t.Log(utils.PrintJson(resp.Tweets))
}
