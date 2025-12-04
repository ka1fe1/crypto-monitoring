package service

import (
	"fmt"

	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
)

type TokenService interface {
	GetTokenPrice(ids []string) (map[string]utils.TokenInfo, error)
}

type tokenService struct {
	client *utils.CoinMarketClient
}

func NewTokenService(client *utils.CoinMarketClient) TokenService {
	return &tokenService{
		client: client,
	}
}

func (s *tokenService) GetTokenPrice(ids []string) (map[string]utils.TokenInfo, error) {
	prices, err := s.client.GetPrice(ids)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch token prices: %w", err)
	}

	return prices, nil
}
