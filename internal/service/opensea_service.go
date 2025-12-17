package service

import (
	"fmt"
	"strings"

	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/opensea"
)

type NFTFloorPriceInfo struct {
	CollectionSlug   string  `json:"collection_slug"`
	FloorPrice       float64 `json:"floor_price"`
	FloorPriceSymbol string  `json:"floor_price_symbol"`
	FloorPriceUSD    float64 `json:"floor_price_usd"`
}

type OpenSeaService interface {
	GetNFTFloorPrices(slugs []string, convertToUsd bool) ([]NFTFloorPriceInfo, error)
}

type openSeaService struct {
	openSeaClient *opensea.OpenSeaClient
	cmcClient     *utils.CoinMarketClient
}

func NewOpenSeaService(openSeaClient *opensea.OpenSeaClient, cmcClient *utils.CoinMarketClient) OpenSeaService {
	return &openSeaService{
		openSeaClient: openSeaClient,
		cmcClient:     cmcClient,
	}
}

func (s *openSeaService) GetNFTFloorPrices(slugs []string, convertToUsd bool) ([]NFTFloorPriceInfo, error) {
	var results []NFTFloorPriceInfo

	for _, slug := range slugs {
		stats, err := s.openSeaClient.GetCollectionStats(slug)
		if err != nil {
			return nil, fmt.Errorf("failed to get stats for %s: %w", slug, err)
		}

		info := NFTFloorPriceInfo{
			CollectionSlug:   slug,
			FloorPrice:       stats.FloorPrice,
			FloorPriceSymbol: stats.FloorPriceSymbol,
		}

		if convertToUsd {
			symbol := strings.ToUpper(stats.FloorPriceSymbol)
			if symbol == "USD" || symbol == "USDT" || symbol == "USDC" {
				info.FloorPriceUSD = stats.FloorPrice
			} else if symbol != "" {
				// Fetch price from CMC
				prices, err := s.cmcClient.GetPriceBySymbol([]string{symbol})
				if err == nil {
					// Check for exact symbol match
					if tokenInfo, ok := prices[symbol]; ok {
						info.FloorPriceUSD = stats.FloorPrice * tokenInfo.Price
					}
				} else {
					// Just log error internally or ignore, we can't convert so leave as 0
					fmt.Printf("failed to fetch price for symbol %s: %v\n", symbol, err)
				}
			}
		}

		results = append(results, info)
	}

	return results, nil
}
