package tasks

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ka1fe1/crypto-monitoring/config"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/polymarket"
)

var (
	testCfg    *config.Config
	polyClient *polymarket.Client
)

func TestPolymarketDailyReportTask_Run(t *testing.T) {
	testCfg = cfg
	polyClient = polymarket.NewClient(cfg.Polymarket.APIKey)

	// Ensure config is set
	if testCfg.PolymarketReport.AddressListFile == "" {
		t.Fatal("AddressListFile not configured in config.yaml")
	}
	if testCfg.PolymarketReport.OutputDir == "" {
		t.Fatal("OutputDir not configured in config.yaml")
	}

	var qh utils.QuietHoursParams
	task := NewPolymarketDailyReportTask(testCfg, polyClient, 86400, qh)
	task.Run()

	// Verify a report file was generated
	outputDir := testCfg.PolymarketReport.OutputDir
	files, err := os.ReadDir(outputDir)
	if err != nil {
		t.Fatalf("Failed to read output dir: %v", err)
	}

	if len(files) == 0 {
		t.Fatal("No report file was generated")
	}

	// Read and log the latest generated report
	latestFile := files[len(files)-1]
	reportPath := filepath.Join(outputDir, latestFile.Name())
	content, err := os.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("Failed to read report file: %v", err)
	}

	t.Logf("Generated report file: %s", latestFile.Name())
	t.Logf("Report content:\n%s", string(content))
}
