package handlers

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ka1fe1/crypto-monitoring/config"
)

// PolymarketReportHandler handles requests for Polymarket daily report data.
type PolymarketReportHandler struct {
	cfg *config.Config
}

// NewPolymarketReportHandler creates a new handler instance.
func NewPolymarketReportHandler(cfg *config.Config) *PolymarketReportHandler {
	return &PolymarketReportHandler{cfg: cfg}
}

// ReportResponse is the JSON response for the report API.
type ReportResponse struct {
	Filename    string     `json:"filename"`
	GeneratedAt string     `json:"generatedAt"`
	Headers     []string   `json:"headers"`
	Rows        [][]string `json:"rows"`
}

// GetLatestReport reads the latest report file and returns parsed table data as JSON.
func (h *PolymarketReportHandler) GetLatestReport(c *gin.Context) {
	outputDir := h.cfg.PolymarketReport.OutputDir
	if outputDir == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OutputDir not configured"})
		return
	}

	// Find the latest .md file in the output directory
	filename, err := findLatestReportFile(outputDir)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No report files found"})
		return
	}

	filePath := filepath.Join(outputDir, filename)
	headers, rows, generatedAt, err := parseMarkdownTable(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to parse report: %v", err)})
		return
	}

	// Identify if we have the specific wallet columns to merge
	var mergedHeaders []string
	hasWalletCols := len(headers) >= 3 && headers[0] == "wallet_addr" && headers[1] == "wallet_name" && headers[2] == "proxy_addr"

	if hasWalletCols {
		mergedHeaders = append([]string{"wallet_info"}, headers[3:]...)
	} else {
		mergedHeaders = headers
	}

	var mergedRows [][]string
	for _, row := range rows {
		if hasWalletCols && len(row) >= 3 {
			// Merge first 3 columns: addr|name|proxy
			walletInfo := fmt.Sprintf("%s|%s|%s", row[0], row[1], row[2])
			newRow := append([]string{walletInfo}, row[3:]...)
			mergedRows = append(mergedRows, newRow)
		} else {
			mergedRows = append(mergedRows, row)
		}
	}

	c.JSON(http.StatusOK, ReportResponse{
		Filename:    filename,
		GeneratedAt: generatedAt,
		Headers:     mergedHeaders,
		Rows:        mergedRows,
	})
}

// findLatestReportFile returns the filename of the most recent report in the directory.
func findLatestReportFile(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	var mdFiles []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
			mdFiles = append(mdFiles, e.Name())
		}
	}

	if len(mdFiles) == 0 {
		return "", fmt.Errorf("no .md files found")
	}

	// Sort by name (lexicographic = chronological due to YYYYMMDD_HHMM format)
	sort.Strings(mdFiles)
	return mdFiles[len(mdFiles)-1], nil
}

// parseMarkdownTable reads a markdown file and extracts the table headers and rows.
// It also extracts the generated timestamp from the H1 header line.
func parseMarkdownTable(filePath string) ([]string, [][]string, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, "", err
	}
	defer file.Close()

	var headers []string
	var rows [][]string
	generatedAt := ""

	scanner := bufio.NewScanner(file)
	// Increase buffer size for long lines (positions can be very long)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Extract generated timestamp from header like: # Polymarket Trader Daily Report - 2026-03-02 16:27:00
		if strings.HasPrefix(line, "# ") && strings.Contains(line, " - ") {
			parts := strings.SplitN(line, " - ", 2)
			if len(parts) == 2 {
				generatedAt = strings.TrimSpace(parts[1])
			}
			continue
		}

		// Skip non-table lines
		if !strings.HasPrefix(line, "|") {
			continue
		}

		// Skip separator lines (|---|---|...)
		if strings.Contains(line, "---") {
			continue
		}

		// Parse table cells
		cells := parsePipeLine(line)

		if len(headers) == 0 {
			headers = cells
		} else {
			rows = append(rows, cells)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, "", err
	}

	return headers, rows, generatedAt, nil
}

// parsePipeLine splits a markdown table row by | and trims each cell.
// It respects the 9 columns and handles escaped pipes in the last column.
func parsePipeLine(line string) []string {
	// Remove leading and trailing pipes
	trimmed := strings.Trim(line, "|")
	parts := strings.Split(trimmed, "|")

	var cells []string
	for i := 0; i < len(parts); i++ {
		p := strings.TrimSpace(parts[i])
		if i < 8 {
			cells = append(cells, p)
		} else {
			// For the 9th column (current_position), join back any remaining parts
			// because they might have been split by escaped pipes \|
			remaining := parts[i:]
			lastCell := strings.Join(remaining, " | ")
			// Clean up escaped pipes for display if needed, but we prefer keeping them for JS to handle
			cells = append(cells, lastCell)
			break
		}
	}
	return cells
}
