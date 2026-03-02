package markdown

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// WalletEntry holds a wallet name and address parsed from the input markdown table.
type WalletEntry struct {
	Name    string
	Address string
}

// TraderReportData holds the merged data for a single trader.
type TraderReportData struct {
	WalletName       string
	Address          string
	ProxyAddr        string
	Volume           float64
	Rank             string
	Pnl              float64
	PositionValue    float64
	LastActiveTime   string
	CurrentPositions string
}

// ParseAddressList reads a markdown table file with columns: wallet_name | wallet_addr
// and returns a list of WalletEntry.
func ParseAddressList(filePath string) ([]WalletEntry, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open address list file: %w", err)
	}
	defer file.Close()

	var entries []WalletEntry
	addressRegex := regexp.MustCompile(`(?i)(0x[a-f0-9]{40})`)

	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		// Skip non-table lines and header/separator rows
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "|") {
			continue
		}
		// Skip separator line (e.g. |---|---| )
		if strings.Contains(trimmed, "---") {
			continue
		}

		// Split by pipe and trim
		parts := strings.Split(trimmed, "|")
		// A valid row like "| name | addr |" splits into ["", " name ", " addr ", ""]
		var cells []string
		for _, p := range parts {
			c := strings.TrimSpace(p)
			if c != "" {
				cells = append(cells, c)
			}
		}

		if len(cells) < 2 {
			continue
		}

		name := cells[0]
		addrCell := cells[1]

		// Skip header row (contains "wallet_name" or "wallet_addr")
		if strings.EqualFold(name, "wallet_name") || strings.EqualFold(addrCell, "wallet_addr") {
			continue
		}

		// Extract address from cell
		matches := addressRegex.FindAllString(addrCell, -1)
		if len(matches) == 0 {
			continue
		}

		addr := matches[0]
		// Deduplicate
		duplicate := false
		for _, e := range entries {
			if strings.EqualFold(e.Address, addr) {
				duplicate = true
				break
			}
		}
		if !duplicate {
			entries = append(entries, WalletEntry{Name: name, Address: addr})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading address list file: %w", err)
	}

	return entries, nil
}

// WriteReportTable generates a markdown table and writes it to a new file in the output directory.
func WriteReportTable(outputDir string, data []TraderReportData) error {
	// Create output dir if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	timestamp := time.Now().Format("20060102_1504")
	filename := fmt.Sprintf("polymarket_volume_%s.md", timestamp)
	filePath := filepath.Join(outputDir, filename)

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create report file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write Report Header
	_, err = writer.WriteString(fmt.Sprintf("# Polymarket Trader Daily Report - %s\n\n", time.Now().Format("2006-01-02 15:04:05")))
	if err != nil {
		return err
	}

	// Write Table Headers
	_, err = writer.WriteString("| wallet_addr | wallet_name | proxy_addr | total_volume | vol_rank | total_pnl | position_value | last_active | current_position |\n")
	if err != nil {
		return err
	}
	_, err = writer.WriteString("|---|---|---|---|---|---|---|---|---|\n")
	if err != nil {
		return err
	}

	// Write Table Rows
	for _, row := range data {
		line := fmt.Sprintf("| `%s` | %s | `%s` | $%.2f | %s | $%.2f | $%.2f | %s | %s |\n",
			row.Address,
			row.WalletName,
			row.ProxyAddr,
			row.Volume,
			row.Rank,
			row.Pnl,
			row.PositionValue,
			row.LastActiveTime,
			escapeMarkdown(row.CurrentPositions),
		)
		_, err = writer.WriteString(line)
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

// Helper functions

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, item) {
			return true
		}
	}
	return false
}

func escapeMarkdown(text string) string {
	// Simple escape for pipe characters in table cells
	return strings.ReplaceAll(text, "|", "\\|")
}
