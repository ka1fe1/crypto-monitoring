package bgeometrics

type MetricData struct {
	Date          string  `json:"d"`
	UnixTs        int64   `json:"unixTs"`
	BalancedPrice float64 `json:"balancedPrice"`
}
