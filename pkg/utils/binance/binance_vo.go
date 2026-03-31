package binance

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Kline represents a single candlestick data point from Binance.
type Kline struct {
	OpenTime                 int64   `json:"open_time"`
	Open                     float64 `json:"open"`
	High                     float64 `json:"high"`
	Low                      float64 `json:"low"`
	Close                    float64 `json:"close"`
	Volume                   float64 `json:"volume"`
	CloseTime                int64   `json:"close_time"`
	QuoteAssetVolume         float64 `json:"quote_asset_volume"`
	NumberOfTrades           int64   `json:"number_of_trades"`
	TakerBuyBaseAssetVolume  float64 `json:"taker_buy_base_asset_volume"`
	TakerBuyQuoteAssetVolume float64 `json:"taker_buy_quote_asset_volume"`
}

// UnmarshalJSON custom unmarshaler to parse Binance array format into struct
func (k *Kline) UnmarshalJSON(buf []byte) error {
	var tmp []interface{}
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if len(tmp) < 11 {
		return fmt.Errorf("invalid kline data length: %d", len(tmp))
	}

	// Helper function to parse float from either string or float64 safely
	parseFloat := func(v interface{}) (float64, error) {
		switch t := v.(type) {
		case string:
			return strconv.ParseFloat(t, 64)
		case float64:
			return t, nil
		default:
			return 0, fmt.Errorf("unexpected type for float: %T", v)
		}
	}

	var err error
	ot, err := parseFloat(tmp[0])
	if err != nil {
		return fmt.Errorf("invalid open_time: %w", err)
	}
	k.OpenTime = int64(ot)
	if k.Open, err = parseFloat(tmp[1]); err != nil {
		return err
	}
	if k.High, err = parseFloat(tmp[2]); err != nil {
		return err
	}
	if k.Low, err = parseFloat(tmp[3]); err != nil {
		return err
	}
	if k.Close, err = parseFloat(tmp[4]); err != nil {
		return err
	}
	if k.Volume, err = parseFloat(tmp[5]); err != nil {
		return err
	}
	ct, err := parseFloat(tmp[6])
	if err != nil {
		return fmt.Errorf("invalid close_time: %w", err)
	}
	k.CloseTime = int64(ct)
	if k.QuoteAssetVolume, err = parseFloat(tmp[7]); err != nil {
		return err
	}
	numT, err := parseFloat(tmp[8])
	if err != nil {
		return fmt.Errorf("invalid number_of_trades: %w", err)
	}
	k.NumberOfTrades = int64(numT)
	if k.TakerBuyBaseAssetVolume, err = parseFloat(tmp[9]); err != nil {
		return err
	}
	if k.TakerBuyQuoteAssetVolume, err = parseFloat(tmp[10]); err != nil {
		return err
	}

	return nil
}
