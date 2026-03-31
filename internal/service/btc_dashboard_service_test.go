package service

import (
	"testing"
	"time"

	"github.com/ka1fe1/crypto-monitoring/pkg/utils/alternative"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/binance"
)

type MockBinanceProvider struct{}

func (m *MockBinanceProvider) GetKlines(symbol, interval string, limit int) ([]binance.Kline, error) {
	klines := make([]binance.Kline, limit)
	for i := 0; i < limit; i++ {
		klines[i] = binance.Kline{
			OpenTime: time.Now().Add(-time.Duration(limit-i) * 24 * time.Hour).UnixMilli(),
			Close:    60000.0,
		}
	}
	return klines, nil
}

type MockMempoolProvider struct{}

func (m *MockMempoolProvider) GetTipHeight() (int64, error) {
	return 840000, nil
}

type MockAlternativeProvider struct{}

func (m *MockAlternativeProvider) GetFng(limit int) (*alternative.FngResponse, error) {
	return &alternative.FngResponse{
		Data: []alternative.FngData{
			{Value: "70", ValueClassification: "Greed"},
		},
	}, nil
}

type MockBalancedPriceProvider struct{}

func (m *MockBalancedPriceProvider) GetBalancedPrice() (float64, error) {
	return 40000.0, nil
}

func TestFetchAndCalculateMetrics_Mock(t *testing.T) {
	mockSvc := NewBtcDashboardService(
		&MockBinanceProvider{},
		&MockMempoolProvider{},
		&MockAlternativeProvider{},
		&MockBalancedPriceProvider{},
	)

	metrics, err := mockSvc.FetchAndCalculateMetrics()
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

	report := mockSvc.GenerateMarkdownReport(metrics)
	if report == "" {
		t.Error("expected non-empty report")
	}
	t.Logf("\n[Mock Report]\n%s", report)
}

func TestFetchAndCalculateMetrics_Real(t *testing.T) {
	if btcDashboardSvc == nil {
		t.Skip("real service not initialized")
	}

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
	t.Logf("\n[Real Report]\n%s", report)
}
