package stock

import "time"

// StockInfo 存储单只股票（港股/A股）的实时行情信息
type StockInfo struct {
	Code          string    // 股票代码，如 hk00700, sh600519
	Name          string    // 股票名称，如 腾讯控股, 贵州茅台
	Price         float64   // 当前价格
	PercentChange float64   // 涨跌幅 (%)
	LastUpdated   time.Time // 数据更新时间
}
