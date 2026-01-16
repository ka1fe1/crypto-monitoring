package service

import (
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/polymarket"
)

type PolymarketMonitorService interface {
	GetMarketDetails(ids []string) ([]polymarket.MarketDetail, error)
}

type polymarketMonitorService struct {
	client *polymarket.Client
}

func NewPolymarketMonitorService(client *polymarket.Client) PolymarketMonitorService {
	return &polymarketMonitorService{
		client: client,
	}
}

func (s *polymarketMonitorService) GetMarketDetails(ids []string) ([]polymarket.MarketDetail, error) {
	var markets []polymarket.MarketDetail

	for _, id := range ids {
		market, err := s.client.GetMarketDetail(id)
		if err != nil {
			continue
		}
		markets = append(markets, *market)
	}
	return markets, nil
}
