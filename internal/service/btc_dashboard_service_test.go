package service

import (
	"testing"
)

func TestFetchAndCalculateMetrics(t *testing.T) {
	metrics, err := btcDashboardSvc.FetchAndCalculateMetrics()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if metrics.CurrentPrice == 0 {
		t.Error("expected non-zero current price")
	}
	if metrics.WMA200 == 0 {
		t.Error("expected non-zero wma200")
	}
	if metrics.Ahr999 == 0 {
		t.Error("expected non-zero ahr999")
	}

	report := btcDashboardSvc.GenerateMarkdownReport(metrics)
	if report == "" {
		t.Error("expected non-empty report")
	}
	t.Logf("\n%s", report)
}
