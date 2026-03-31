package service

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/ka1fe1/crypto-monitoring/pkg/logger"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/alternative"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/binance"
)

const (
	ahr999CoefficientA = 5.84
	ahr999CoefficientB = 17.01

	wmaThresholdExtremeBear = 1.0
	wmaThresholdBottom      = 1.5
	wmaThresholdOverheated  = 3.0

	ahr999ThresholdBottom     = 0.45
	ahr999ThresholdInvest     = 1.2
	ahr999ThresholdWaitAndSee = 5.0

	bpRatioThresholdUndervalued = 1.0
	bpRatioThresholdOvervalued  = 2.0
)

type BtcDashboardMetrics struct {
	CurrentPrice   float64
	WMA200         float64
	WMARatio       float64
	Ahr999         float64
	FGIValue       int
	FGIClass       string
	HalvingDays    int
	HalvingBlock   int64
	HalvingEstDate time.Time
	BalancedPrice  float64
	BPRatio        float64
}

type BtcDashboardService interface {
	FetchAndCalculateMetrics() (*BtcDashboardMetrics, error)
	GenerateMarkdownReport(metrics *BtcDashboardMetrics) string
}

type BinanceProvider interface {
	GetKlines(symbol, interval string, limit int) ([]binance.Kline, error)
}

type MempoolProvider interface {
	GetTipHeight() (int64, error)
}

type AlternativeProvider interface {
	GetFng(limit int) (*alternative.FngResponse, error)
}

type BalancedPriceProvider interface {
	GetBalancedPrice() (float64, error)
}

type btcDashboardService struct {
	binanceClient     BinanceProvider
	mempoolClient     MempoolProvider
	alternativeClient AlternativeProvider
	bpClient          BalancedPriceProvider
}

func NewBtcDashboardService(b BinanceProvider, m MempoolProvider, a AlternativeProvider, bp BalancedPriceProvider) BtcDashboardService {
	return &btcDashboardService{
		binanceClient:     b,
		mempoolClient:     m,
		alternativeClient: a,
		bpClient:          bp,
	}
}

func (s *btcDashboardService) FetchAndCalculateMetrics() (*BtcDashboardMetrics, error) {
	metrics := &BtcDashboardMetrics{}

	// 1. Fetch 200 Weekly Klines for WMA
	weeklyKlines, err := s.binanceClient.GetKlines("BTCUSDT", "1w", 200)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weekly klines: %w", err)
	}
	if len(weeklyKlines) > 0 {
		var sum float64
		for _, k := range weeklyKlines {
			sum += k.Close
		}
		metrics.WMA200 = sum / float64(len(weeklyKlines))
	}

	// 2. Fetch 200 Daily Klines for ahr999
	dailyKlines, err := s.binanceClient.GetKlines("BTCUSDT", "1d", 200)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch daily klines: %w", err)
	}
	if len(dailyKlines) > 0 {
		var sum float64
		for _, k := range dailyKlines {
			sum += k.Close
		}
		dma200 := sum / float64(len(dailyKlines))
		currentPrice := dailyKlines[len(dailyKlines)-1].Close
		metrics.CurrentPrice = currentPrice

		if metrics.WMA200 > 0 {
			metrics.WMARatio = currentPrice / metrics.WMA200
		}

		// Calculate ahr999
		genesisTs := time.Date(2009, 1, 3, 0, 0, 0, 0, time.UTC).UnixMilli()
		nowTs := time.Now().UnixMilli()

		lastKLineTs := dailyKlines[len(dailyKlines)-1].OpenTime
		if lastKLineTs > 0 {
			// 使用最新日线的开盘时间计算天数，确保与价格数据的取样时间点对齐
			nowTs = lastKLineTs
		}

		coinDays := float64(nowTs-genesisTs) / 86400000.0
		expPrice := math.Pow(10, ahr999CoefficientA*math.Log10(coinDays)-ahr999CoefficientB)

		if dma200 > 0 && expPrice > 0 {
			metrics.Ahr999 = (currentPrice / dma200) * (currentPrice / expPrice)
		}
	}

	// 3. Fetch Mempool Tip Height for Halving
	height, err := s.mempoolClient.GetTipHeight()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tip height: %w", err)
	}
	if height > 0 {
		halvingInterval := int64(210000)
		nextHalving := ((height / halvingInterval) + 1) * halvingInterval
		remaining := nextHalving - height
		// 注意：这里的 10 分钟仅为比特币协议设计的预期平均出块时间，实际会根据全网算力波动而变化，仅作为概略估算
		daysLeft := float64(remaining*10) / (60.0 * 24.0)

		metrics.HalvingBlock = nextHalving
		metrics.HalvingDays = int(math.Round(daysLeft))
		metrics.HalvingEstDate = time.Now().Add(time.Duration(daysLeft*24) * time.Hour)
	}

	// 4. Fetch FGI
	fng, err := s.alternativeClient.GetFng(1)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch fng: %w", err)
	}
	if len(fng.Data) > 0 {
		v, _ := strconv.Atoi(fng.Data[0].Value)
		metrics.FGIValue = v
		metrics.FGIClass = fng.Data[0].ValueClassification
	}

	// 5. Fetch Balanced Price
	if s.bpClient != nil {
		bp, err := s.bpClient.GetBalancedPrice()
		if err != nil {
			logger.Warn("failed to fetch balanced price: %v", err)
		} else if bp > 0 && metrics.CurrentPrice > 0 {
			metrics.BalancedPrice = bp
			metrics.BPRatio = metrics.CurrentPrice / bp
		}
	}

	return metrics, nil
}

func (s *btcDashboardService) GenerateMarkdownReport(metrics *BtcDashboardMetrics) string {
	report := "### 📉 BTC 宏观周期指标监控\n\n"
	report += fmt.Sprintf("- **当前价格**: $%.2f\n", metrics.CurrentPrice)

	// WMA200
	wmaStatus := "正常牛市区间"
	if metrics.WMARatio < wmaThresholdExtremeBear {
		wmaStatus = "极端熊市信号"
	} else if metrics.WMARatio < wmaThresholdBottom {
		wmaStatus = "历史底部区间"
	} else if metrics.WMARatio >= wmaThresholdOverheated {
		wmaStatus = "过热信号"
	}
	report += fmt.Sprintf("- **200 周均线 (200WMA)**: $%.2f\n  - 偏离度: %.2fx (状态: %s)\n",
		metrics.WMA200, metrics.WMARatio, wmaStatus)

	// Balanced Price
	if metrics.BalancedPrice > 0 {
		bpStatus := "过高区间 🚨"
		if metrics.BPRatio <= bpRatioThresholdUndervalued {
			bpStatus = "基于成本线大底 📉"
		} else if metrics.BPRatio <= bpRatioThresholdOvervalued {
			bpStatus = "正常估值区间 📈"
		}

		report += fmt.Sprintf("- **均衡价格 (BP)**: $%.2f\n  - 当前售价 / BP: %.2fx (状态: %s)\n",
			metrics.BalancedPrice, metrics.BPRatio, bpStatus)
	}

	// Ahr999
	ahrStatus := "泡沫区间"
	if metrics.Ahr999 < ahr999ThresholdBottom {
		ahrStatus = "抄底区间"
	} else if metrics.Ahr999 < ahr999ThresholdInvest {
		ahrStatus = "定投区间"
	} else if metrics.Ahr999 < ahr999ThresholdWaitAndSee {
		ahrStatus = "观望区间"
	}
	report += fmt.Sprintf("- **ahr999 定投指数**: %.3f (状态: %s)\n", metrics.Ahr999, ahrStatus)

	// FGI
	report += fmt.Sprintf("- **恐惧贪婪指数**: %d (%s)\n", metrics.FGIValue, metrics.FGIClass)

	// Halving
	report += fmt.Sprintf("- **减半倒计时**: 还有约 %d 天\n  - 目标区块: %d\n  - 预计时间: %s\n",
		metrics.HalvingDays, metrics.HalvingBlock, metrics.HalvingEstDate.Format("2006-01-02"))

	return report
}
