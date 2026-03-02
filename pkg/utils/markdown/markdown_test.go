package markdown

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseAddressList(t *testing.T) {
	// Create a temporary file with markdown table format
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test_addresses.md")

	content := []byte(`| wallet_name | wallet_addr |
|---|---|
| alice | 0x71c7656ec7ab88b098defB751B7401B5f6d8976f |
| bob | 0x1234567890123456789012345678901234567890 |
| alice_dup | 0x71c7656ec7ab88b098defB751B7401B5f6d8976f |
`)

	err := os.WriteFile(tempFile, content, 0644)
	if err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	entries, err := ParseAddressList(tempFile)
	if err != nil {
		t.Fatalf("ParseAddressList returned error: %v", err)
	}

	if len(entries) != 2 {
		t.Errorf("Expected 2 unique entries, got %d", len(entries))
	}

	if entries[0].Name != "alice" {
		t.Errorf("Expected first name to be 'alice', got '%s'", entries[0].Name)
	}
	if entries[0].Address != "0x71c7656ec7ab88b098defB751B7401B5f6d8976f" {
		t.Errorf("Expected first address to be 0x71c7..., got %s", entries[0].Address)
	}
	if entries[1].Name != "bob" {
		t.Errorf("Expected second name to be 'bob', got '%s'", entries[1].Name)
	}
}

func TestWriteReportTable(t *testing.T) {
	tempDir := t.TempDir()

	data := []TraderReportData{
		{
			WalletName:       "alice",
			Address:          "0x1111111111111111111111111111111111111111",
			ProxyAddr:        "0xaaaa111111111111111111111111111111111111",
			Volume:           1000.50,
			Rank:             "1",
			Pnl:              250.00,
			PositionValue:    500.25,
			CurrentPositions: "Long BTC|Short ETH",
		},
		{
			WalletName:       "bob",
			Address:          "0x2222222222222222222222222222222222222222",
			ProxyAddr:        "0xbbbb222222222222222222222222222222222222",
			Volume:           0,
			Rank:             "2",
			Pnl:              0,
			PositionValue:    0,
			CurrentPositions: "None",
		},
	}

	err := WriteReportTable(tempDir, data)
	if err != nil {
		t.Fatalf("WriteReportTable returned error: %v", err)
	}

	files, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("Could not read temp dir: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("Expected 1 file to be created, found %d", len(files))
	}

	fileContent, err := os.ReadFile(filepath.Join(tempDir, files[0].Name()))
	if err != nil {
		t.Fatalf("Could not read generated file: %v", err)
	}

	contentStr := string(fileContent)

	if !strings.Contains(contentStr, "| alice | `0x1111111111111111111111111111111111111111` |") {
		t.Errorf("Generated file missing expected row 1 format\nGot: %s", contentStr)
	}

	if !strings.Contains(contentStr, "| bob | `0x2222222222222222222222222222222222222222` |") {
		t.Errorf("Generated file missing expected row 2 format\nGot: %s", contentStr)
	}
}
