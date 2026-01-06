package utils

import (
	"encoding/json"
	"fmt"
	"time"
)

func PrintJson(obj interface{}) string {
	b, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return fmt.Sprintf("marshal error: %v", err.Error())
	}
	return string(b)
}

func FormatBJTime(t time.Time) string {
	location := time.FixedZone("CST", 8*3600)
	return t.In(location).Format("2006-01-02 15:04:05")
}

func FormatPrice(price float64) string {
	if price >= 1 {
		return fmt.Sprintf("%.2f", price)
	}
	return fmt.Sprintf("%.4f", price)
}
