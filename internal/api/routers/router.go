package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/ka1fe1/crypto-monitoring/config"
	_ "github.com/ka1fe1/crypto-monitoring/docs"
	"github.com/ka1fe1/crypto-monitoring/internal/api/handlers"
	"github.com/ka1fe1/crypto-monitoring/internal/service"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/opensea"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter initializes the Gin engine and defines the routes
func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Initialize services and handlers
	cmcClient := utils.NewCoinMarketClient(cfg.CoinMarketCap.APIKey)

	dexPairService := service.NewDexPairService(cmcClient)
	dexPairHandler := handlers.NewDexPairHandler(dexPairService)

	tokenService := service.NewTokenService(cmcClient)
	tokenHandler := handlers.NewTokenHandler(tokenService)

	openSeaClient := opensea.NewOpenSeaClient(cfg.OpenSea.APIKey)
	openSeaService := service.NewOpenSeaService(openSeaClient, cmcClient)
	openSeaHandler := handlers.NewOpenSeaHandler(openSeaService, cfg.NFTFloorPriceMonitor.NFTCollections)

	// Register routes
	r.GET("/ping", handlers.PingHandler)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
		api.GET("/dex/pair", dexPairHandler.GetDexPair)
		api.GET("/token/price", tokenHandler.GetTokenPrice)
		api.GET("/nft/floor_price", openSeaHandler.GetNFTFloorPrice)
	}

	return r
}
