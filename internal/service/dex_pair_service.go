package service

import (
	"fmt"
	"time"

	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
)

type DexPairInfo struct {
	Name            string  `json:"name"`
	Price           float64 `json:"price"`
	PercentChange1h float64 `json:"percent_change_price_1h"`
	DexSlug         string  `json:"dex_slug"`
	NetworkSlug     string  `json:"network_slug"`
	Liquidity       float64 `json:"liquidity"`
	LastUpdated     string  `json:"last_updated"`
}

type DexPairService interface {
	GetDexPairInfo(contractAddresses []string, networkSlug, networkId string) ([]*DexPairInfo, error)
}

type dexPairService struct {
	client *utils.CoinMarketClient
}

func NewDexPairService(client *utils.CoinMarketClient) DexPairService {
	return &dexPairService{
		client: client,
	}
}

func (s *dexPairService) GetDexPairInfo(contractAddresses []string, networkSlug, networkId string) ([]*DexPairInfo, error) {
	pairs, err := s.client.GetDexPairQuotes(contractAddresses, networkSlug, networkId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch dex pair quotes: %w", err)
	}

	var results []*DexPairInfo
	for _, address := range contractAddresses {
		pair, ok := pairs[address]
		if !ok {
			// If not found in map, it might happen if API didn't return it.
			continue
		}

		info := &DexPairInfo{
			Name:        pair.Name,
			DexSlug:     pair.DexSlug,
			NetworkSlug: pair.NetworkSlug,
		}

		if len(pair.Quote) > 0 {
			quote := pair.Quote[0]
			info.Price = quote.Price
			info.PercentChange1h = quote.PercentChange1h
			info.Liquidity = quote.Liquidity

			// Format LastUpdated
			parsedTime, err := time.Parse(time.RFC3339, quote.LastUpdated)
			if err == nil {
				info.LastUpdated = utils.FormatBJTime(parsedTime)
			} else {
				info.LastUpdated = quote.LastUpdated
			}
		}
		results = append(results, info)
	}

	return results, nil
}
