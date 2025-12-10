package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server        ServerConfig        `yaml:"server"`
	CoinMarketCap CoinMarketCapConfig `yaml:"coinmarketcap"`
	DingTalk      DingTalkConfig      `yaml:"dingtalk"`
	DexPairAlter  DexPairAlterConfig  `yaml:"dex_pair_alter"`
	BinanceCex    BinanceCexConfig    `yaml:"binance-cex"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type CoinMarketCapConfig struct {
	APIKey string `yaml:"api_key"`
}

type DingTalkConfig struct {
	AccessToken string `yaml:"access_token"`
	Secret      string `yaml:"secret"`
	Keyword     string `yaml:"keyword"`
}

type DexPairAlterConfig struct {
	ContractAddress string `yaml:"contract_address"`
	NetworkSlug     string `yaml:"network_slug"`
	IntervalSeconds int    `yaml:"interval_seconds"`
}

type BinanceCexConfig struct {
	APIKey    string `yaml:"api_key"`
	SecretKey string `yaml:"secret_key"`
	ProxyURL  string `yaml:"proxy_url"`
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return &cfg, nil
}
