package tasks

import (
	"log"

	"github.com/ka1fe1/crypto-monitoring/config"
	"github.com/ka1fe1/crypto-monitoring/internal/service"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/alter/dingding"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/polymarket"
)

func InitTasks(
	cfg *config.Config,
	dingBots map[string]*dingding.DingBot,
	dexService service.DexPairService,
	tokenService service.TokenService,
	openSeaService service.OpenSeaService,
	polyClient *polymarket.Client,
) {
	// 1. DexPairAlterTask
	dexBot := dingBots[cfg.DexPairAlter.BotName]
	if dexBot != nil {
		NewDexPairAlterTask(dexService, dexBot, cfg.DexPairAlter.ContractAddrInfo, cfg.DexPairAlter.IntervalSeconds).Start()
	} else {
		log.Printf("Warning: Bot %s not found for DexPairAlterTask", cfg.DexPairAlter.BotName)
	}

	// 2. TokenPriceMonitorTask
	tokenBot := dingBots[cfg.TokenPriceMonitor.BotName]
	if tokenBot != nil {
		NewTokenPriceMonitorTask(tokenService, tokenBot, cfg.TokenPriceMonitor.TokenIds, cfg.TokenPriceMonitor.IntervalSeconds).Start()
	} else {
		log.Printf("Warning: Bot %s not found for TokenPriceMonitorTask", cfg.TokenPriceMonitor.BotName)
	}

	// 3. NFTFloorPriceMonitorTask
	nftBot := dingBots[cfg.NFTFloorPriceMonitor.BotName]
	if nftBot != nil {
		NewNFTFloorPriceMonitorTask(openSeaService, nftBot, cfg.NFTFloorPriceMonitor.NFTCollections, cfg.NFTFloorPriceMonitor.IntervalSeconds).Start()
	} else {
		log.Printf("Warning: Bot %s not found for NFTFloorPriceMonitorTask", cfg.NFTFloorPriceMonitor.BotName)
	}

	// 4. PolymarketMonitorTask
	polyBot := dingBots[cfg.PolymarketMonitor.BotName]
	if polyBot != nil {
		NewPolymarketMonitorTask(polyClient, polyBot, cfg.PolymarketMonitor.MarketIDs, cfg.PolymarketMonitor.IntervalSeconds).Start()
	} else {
		log.Printf("Warning: Bot %s not found for PolymarketMonitorTask", cfg.PolymarketMonitor.BotName)
	}
}
