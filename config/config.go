package config

import (
	"fmt"
	"os"

	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server               ServerConfig               `yaml:"server"`
	CoinMarketCap        CoinMarketCapConfig        `yaml:"coinmarketcap"`
	DingTalk             map[string]DingTalkConfig  `yaml:"dingtalk"`
	DexPairAlter         DexPairAlterConfig         `yaml:"dex_pair_alter"`
	TokenPriceMonitor    TokenPriceMonitorConfig    `yaml:"token_price_monitor"`
	BinanceCex           BinanceCexConfig           `yaml:"binance-cex"`
	OpenSea              OpenSeaConfig              `yaml:"opensea"`
	NFTFloorPriceMonitor NFTFloorPriceMonitorConfig `yaml:"nft_floor_price_monitor"`
	CoinGlass            CoinGlassConfig            `yaml:"coinglass"`
	Polymarket           PolymarketConfig           `yaml:"polymarket"`
	PolymarketMonitor    PolymarketMonitorConfig    `yaml:"polymarket_monitor"`
	Twitter              TwitterConfig              `yaml:"twitter"`
	TwitterMonitor       TwitterMonitorConfig       `yaml:"twitter_monitor"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type CoinMarketCapConfig struct {
	APIKey string `yaml:"api_key"`
}

type OpenSeaConfig struct {
	APIKey string `yaml:"api_key"`
}

type DingTalkConfig struct {
	AccessToken string `yaml:"access_token"`
	Secret      string `yaml:"secret"`
	Keyword     string
}

type DexPairAlterConfig struct {
	ContractAddrs   []string `yaml:"contract_addrs"`
	IntervalSeconds int      `yaml:"interval_seconds"`
	BotName         string   `yaml:"bot_name"`
	// key: networkId, value: contractAddrs
	ContractAddrInfo map[string][]string
}

type TokenPriceMonitorConfig struct {
	TokenIds        string `yaml:"token_ids"`
	IntervalSeconds int    `yaml:"interval_seconds"`
	BotName         string `yaml:"bot_name"`
}

type NFTFloorPriceMonitorConfig struct {
	IntervalSeconds   int      `yaml:"interval_seconds"`
	BotName           string   `yaml:"bot_name"`
	NFTCollectionsStr string   `yaml:"nft_collections"`
	NFTCollections    []string `yaml:"-"`
}

type BinanceCexConfig struct {
	APIKey    string `yaml:"api_key"`
	SecretKey string `yaml:"secret_key"`
	ProxyURL  string `yaml:"proxy_url"`
}

type CoinGlassConfig struct {
	APIKey string `yaml:"api_key"`
}

type PolymarketConfig struct {
	APIKey string `yaml:"api_key"`
}

type TwitterConfig struct {
	APIKey string `yaml:"api_key"`
}

type PolymarketMonitorConfig struct {
	IntervalSeconds int      `yaml:"interval_seconds"`
	BotName         string   `yaml:"bot_name"`
	MarketIDsStr    string   `yaml:"market_ids"`
	MarketIDs       []string `yaml:"-"`
}

type TwitterMonitorConfig struct {
	IntervalSeconds int      `yaml:"interval_seconds"`
	BotName         string   `yaml:"bot_name"`
	UsernamesStr    string   `yaml:"usernames"`
	Usernames       []string `yaml:"-"`
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

	// Parse NFTFloorPriceMonitor NFTCollections
	if cfg.NFTFloorPriceMonitor.NFTCollectionsStr != "" {
		parts := strings.Split(cfg.NFTFloorPriceMonitor.NFTCollectionsStr, ",")
		for _, p := range parts {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" {
				cfg.NFTFloorPriceMonitor.NFTCollections = append(cfg.NFTFloorPriceMonitor.NFTCollections, trimmed)
			}
		}
	}

	// Parse PolymarketMonitor MarketIDs
	if cfg.PolymarketMonitor.MarketIDsStr != "" {
		parts := strings.Split(cfg.PolymarketMonitor.MarketIDsStr, ",")
		for _, p := range parts {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" {
				cfg.PolymarketMonitor.MarketIDs = append(cfg.PolymarketMonitor.MarketIDs, trimmed)
			}
		}
	}

	// Parse TwitterMonitor Usernames
	if cfg.TwitterMonitor.UsernamesStr != "" {
		parts := strings.Split(cfg.TwitterMonitor.UsernamesStr, ",")
		for _, p := range parts {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" {
				cfg.TwitterMonitor.Usernames = append(cfg.TwitterMonitor.Usernames, trimmed)
			}
		}
	}

	// Parse ContractAddrInfo
	cfg.DexPairAlter.ContractAddrInfo = make(map[string][]string)
	for _, entry := range cfg.DexPairAlter.ContractAddrs {
		parts := strings.Split(entry, ":")
		if len(parts) == 2 {
			networkId := strings.TrimSpace(parts[0])
			addrsStr := parts[1]
			addrs := strings.Split(addrsStr, ",")
			var trimmedAddrs []string
			for _, addr := range addrs {
				trimmedAddrs = append(trimmedAddrs, strings.TrimSpace(addr))
			}
			cfg.DexPairAlter.ContractAddrInfo[networkId] = trimmedAddrs
		}
	}

	// set keyword equal to bot name
	for botName, v := range cfg.DingTalk {
		v.Keyword = botName
		cfg.DingTalk[botName] = v
	}
	return &cfg, nil
}
