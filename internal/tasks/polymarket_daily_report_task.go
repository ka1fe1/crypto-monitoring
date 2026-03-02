package tasks

import (
	"fmt"
	"strings"
	"time"

	"github.com/ka1fe1/crypto-monitoring/config"
	"github.com/ka1fe1/crypto-monitoring/pkg/logger"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/markdown"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/polymarket"
)

type PolymarketDailyReportTask struct {
	cfg    *config.PolymarketReportConfig
	client *polymarket.Client
}

func NewPolymarketDailyReportTask(cfg *config.Config, pmClient *polymarket.Client) *PolymarketDailyReportTask {
	return &PolymarketDailyReportTask{
		cfg:    &cfg.PolymarketReport,
		client: pmClient,
	}
}

func (t *PolymarketDailyReportTask) Name() string {
	return "PolymarketDailyReportTask"
}

func (t *PolymarketDailyReportTask) Run() {
	logger.Info("Starting PolymarketDailyReportTask")

	// 1. Read Addresses
	if t.cfg.AddressListFile == "" || t.cfg.OutputDir == "" {
		logger.Error("PolymarketDailyReportTask configuration missing: AddressListFile or OutputDir")
		return
	}

	entries, err := markdown.ParseAddressList(t.cfg.AddressListFile)
	if err != nil {
		logger.Error("Failed to read address list: %v", err)
		return
	}

	if len(entries) == 0 {
		logger.Warn("PolymarketDailyReportTask: No entries found in file: %s", t.cfg.AddressListFile)
		return
	}

	var reportData []markdown.TraderReportData

	// 2. Iterate and Fetch Data
	for _, entry := range entries {
		addr := entry.Address
		// Resolve EOA address to proxyWallet (data-api requires proxyWallet)
		proxyWallet, err := t.client.ResolveProxyWallet(addr)
		if err != nil {
			logger.Error("Failed to resolve proxyWallet for address %s: %v", addr, err)
			continue
		}
		logger.Info("Resolved address %s -> proxyWallet %s", addr, proxyWallet)

		// Fetch Leaderboard for Volume (using proxyWallet)
		lb, err := t.client.GetTraderLeaderboardRankings(proxyWallet)
		if err != nil {
			logger.Error("Failed to get leaderboard for address %s: %v", addr, err)
			continue
		}

		// Fetch Total Value (using proxyWallet)
		tv, err := t.client.GetTotalValueOfUserPositions(proxyWallet)
		if err != nil {
			logger.Error("Failed to get total value for address %s: %v", addr, err)
			continue
		}

		// Fetch Current Positions (using proxyWallet)
		cp, err := t.client.GetCurrentPositionsForUser(proxyWallet)
		var positionsStr string
		if err != nil {
			logger.Error("Failed to get current positions for address %s: %v", addr, err)
			positionsStr = "Error fetching data"
		} else {
			if cp != nil && len(*cp) > 0 {
				var lines []string
				for _, p := range *cp {
					line := fmt.Sprintf("- %s \\| %s \\| init: %.4f(%.2f) \\| current: %.4f(%.2f) \\| cash pnl: %.2f \\| to win: %.2f \\| redeemable: %v",
						p.Title, p.Outcome, p.AvgPrice, p.InitialValue, p.CurPrice, p.CurrentValue, p.CashPnl, p.Size, p.Redeemable)
					lines = append(lines, line)
				}
				positionsStr = strings.Join(lines, " <br> ")
			} else {
				positionsStr = "[]"
			}
		}

		totalVal := 0.0
		if tv != nil {
			totalVal = tv.Value
		}

		vol := 0.0
		pnl := 0.0
		rank := "0"
		if lb != nil {
			vol = lb.Vol
			pnl = lb.Pnl
			rank = lb.Rank
		}

		reportData = append(reportData, markdown.TraderReportData{
			WalletName:       entry.Name,
			Address:          addr,
			ProxyAddr:        proxyWallet,
			Volume:           vol,
			Rank:             rank,
			Pnl:              pnl,
			Value:            totalVal,
			CurrentPositions: positionsStr,
		})

		// Optional: avoid rate limits
		time.Sleep(500 * time.Millisecond)
	}

	// 3. Write Output
	if err := markdown.WriteReportTable(t.cfg.OutputDir, reportData); err != nil {
		logger.Error("Failed to write daily report: %v", err)
		return
	}

	logger.Info("PolymarketDailyReportTask completed successfully. Processed %d addresses.", len(reportData))
}
