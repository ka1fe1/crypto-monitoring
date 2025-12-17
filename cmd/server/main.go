package main

import (
	"log"
	"os"
	"strings"

	"github.com/ka1fe1/crypto-monitoring/config"
	_ "github.com/ka1fe1/crypto-monitoring/docs"
	"github.com/ka1fe1/crypto-monitoring/internal/api/handlers"
	"github.com/ka1fe1/crypto-monitoring/internal/api/routers"
	"github.com/ka1fe1/crypto-monitoring/internal/service"
	"github.com/ka1fe1/crypto-monitoring/internal/tasks"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/alter/dingding"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/opensea"
)

// @title           Crypto Monitoring API
// @version         1.0
// @description     This is a crypto monitoring server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath  /

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	// Load Configuration
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/config.yaml"
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize CoinMarketCap Client
	cmcClient := utils.NewCoinMarketClient(cfg.CoinMarketCap.APIKey)

	// Initialize Services
	dexService := service.NewDexPairService(cmcClient)

	// Initialize Handlers
	_ = handlers.NewDexPairHandler(dexService)

	// Initialize DingTalk Bots
	dingBots := make(map[string]*dingding.DingBot)
	for name, botCfg := range cfg.DingTalk {
		dingBots[name] = dingding.NewDingBot(botCfg.AccessToken, botCfg.Secret, botCfg.Keyword)
	}

	// Initialize Binance Service
	// binanceAnnouncementService := service.NewBinanceAnnouncementService(cfg.BinanceCex.APIKey, cfg.BinanceCex.SecretKey, cfg.BinanceCex.ProxyURL, dingBot)
	// go binanceAnnouncementService.Start(context.Background())

	// Initialize Tasks
	dexPairAlterBot := dingBots[cfg.DexPairAlter.BotName]
	if dexPairAlterBot != nil {
		priceAlertTask := tasks.NewDexPairAlterTask(dexService, dexPairAlterBot, cfg.DexPairAlter.ContractAddrInfo, cfg.DexPairAlter.IntervalSeconds)
		priceAlertTask.Start()
	} else {
		log.Printf("Warning: Bot %s not found for DexPairAlterTask", cfg.DexPairAlter.BotName)
	}

	tokenService := service.NewTokenService(cmcClient)
	tokenPriceMonitorBot := dingBots[cfg.TokenPriceMonitor.BotName]
	if tokenPriceMonitorBot != nil {
		tokenPriceMonitorTask := tasks.NewTokenPriceMonitorTask(tokenService, tokenPriceMonitorBot, cfg.TokenPriceMonitor.TokenIds, cfg.TokenPriceMonitor.IntervalSeconds)
		tokenPriceMonitorTask.Start()
	} else {
		log.Printf("Warning: Bot %s not found for TokenPriceMonitorTask", cfg.TokenPriceMonitor.BotName)
	}

	// Initialize OpenSea Service and NFT Monitor Task
	openSeaClient := opensea.NewOpenSeaClient(cfg.OpenSea.APIKey)
	openSeaService := service.NewOpenSeaService(openSeaClient, cmcClient)

	nftMonitorBot := dingBots[cfg.NFTFloorPriceMonitor.BotName]
	if nftMonitorBot != nil {
		nftMonitorTask := tasks.NewNFTFloorPriceMonitorTask(openSeaService, nftMonitorBot, cfg.NFTFloorPriceMonitor.NFTCollections, cfg.NFTFloorPriceMonitor.IntervalSeconds)
		nftMonitorTask.Start()
	} else {
		log.Printf("Warning: Bot %s not found for NFTFloorPriceMonitorTask", cfg.NFTFloorPriceMonitor.BotName)
	}

	// SetupRouter
	r := routers.SetupRouter(cfg)

	// Start Server
	addr := cfg.Server.Port
	if !strings.HasPrefix(addr, ":") {
		addr = ":" + addr
	}
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
