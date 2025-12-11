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
	client := utils.NewCoinMarketClient(cfg.CoinMarketCap.APIKey)

	// Initialize Services
	dexService := service.NewDexPairService(client)

	// Initialize Handlers
	_ = handlers.NewDexPairHandler(dexService)

	// Initialize DingTalk Bot
	dingBot := dingding.NewDingBot(cfg.DingTalk.AccessToken, cfg.DingTalk.Secret, cfg.DingTalk.Keyword)

	// Initialize Binance Service
	// binanceAnnouncementService := service.NewBinanceAnnouncementService(cfg.BinanceCex.APIKey, cfg.BinanceCex.SecretKey, cfg.BinanceCex.ProxyURL, dingBot)
	// go binanceAnnouncementService.Start(context.Background())

	// Initialize Tasks
	priceAlertTask := tasks.NewDexPairAlterTask(dexService, dingBot, cfg.DexPairAlter.ContractAddrInfo, cfg.DexPairAlter.IntervalSeconds)
	priceAlertTask.Start()

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
