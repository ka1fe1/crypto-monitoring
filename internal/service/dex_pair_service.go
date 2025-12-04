package service

import (
	"fmt"

	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
)

type DexPairInfo struct {
	Name             string  `json:"name"`
	Price            float64 `json:"price"`
	PercentChange1h  float64 `json:"percent_change_price_1h"`
	PercentChange24h float64 `json:"percent_change_price_24h"`
	DexSlug          string  `json:"dex_slug"`
	Liquidity        float64 `json:"liquidity"`
}

type DexPairService interface {
	GetDexPairInfo(contractAddress, networkSlug string) (*DexPairInfo, error)
}

type dexPairService struct {
	client *utils.CoinMarketClient
}

func NewDexPairService(client *utils.CoinMarketClient) DexPairService {
	return &dexPairService{
		client: client,
	}
}

func (s *dexPairService) GetDexPairInfo(contractAddress, networkSlug string) (*DexPairInfo, error) {
	pairs, err := s.client.GetDexPairQuotes([]string{contractAddress}, networkSlug)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch dex pair quotes: %w", err)
	}

	pair, ok := pairs[contractAddress]
	if !ok {
		return nil, fmt.Errorf("pair not found for contract address: %s", contractAddress)
	}

	info := &DexPairInfo{
		Name:    pair.Name,
		DexSlug: pair.DexSlug,
	}

	if len(pair.Quote) > 0 {
		quote := pair.Quote[0]
		info.Price = quote.Price
		info.PercentChange1h = quote.PercentChange1h
		info.PercentChange24h = quote.PercentChange24h
		info.Liquidity = quote.Liquidity
	}

	return info, nil
}
