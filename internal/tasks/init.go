package tasks

import (
	"log"

	"github.com/ka1fe1/crypto-monitoring/config"
	"github.com/ka1fe1/crypto-monitoring/internal/service"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/alter/dingding"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/constant"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/polymarket"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/twitter"
)

func InitTasks(
	cfg *config.Config,
	dingBots map[string]*dingding.DingBot,
	dexService service.DexPairService,
	tokenService service.TokenService,
	openSeaService service.OpenSeaService,
	polyClient *polymarket.Client,
	twitterClient *twitter.TwitterClient,
) {
	// 1. DexPairAlterTask
	if cfg.DexPairAlter.IntervalSeconds > 0 {
		dexBot := dingBots[cfg.DexPairAlter.BotName]
		if dexBot != nil {
			// DexPair: Pause during 00:00-07:00
			qh := utils.QuietHoursParams{Enabled: true, StartHour: 0, EndHour: 7, Behavior: constant.QUIET_HOURS_BEHAVIOR_PAUSE}
			NewDexPairAlterTask(dexService, dexBot, cfg.DexPairAlter.ContractAddrInfo, cfg.DexPairAlter.IntervalSeconds, qh).Start()
		} else {
			log.Printf("Warning: Bot %s not found for DexPairAlterTask", cfg.DexPairAlter.BotName)
		}
	}

	// 2. TokenPriceMonitorTask
	if cfg.TokenPriceMonitor.IntervalSeconds > 0 {
		tokenBot := dingBots[cfg.TokenPriceMonitor.BotName]
		if tokenBot != nil {
			// TokenPrice: Pause during 00:00-07:00
			qh := utils.QuietHoursParams{Enabled: true, StartHour: 0, EndHour: 7, Behavior: constant.QUIET_HOURS_BEHAVIOR_THROTTLE, ThrottleMultiplier: 5}
			NewTokenPriceMonitorTask(tokenService, tokenBot, cfg.TokenPriceMonitor.TokenIds, cfg.TokenPriceMonitor.IntervalSeconds, qh).Start()
		} else {
			log.Printf("Warning: Bot %s not found for TokenPriceMonitorTask", cfg.TokenPriceMonitor.BotName)
		}
	}

	// 3. NFTFloorPriceMonitorTask
	if cfg.NFTFloorPriceMonitor.IntervalSeconds > 0 {
		nftBot := dingBots[cfg.NFTFloorPriceMonitor.BotName]
		if nftBot != nil {
			// NFT: Pause during 00:00-07:00
			qh := utils.QuietHoursParams{Enabled: true, StartHour: 0, EndHour: 7, Behavior: constant.QUIET_HOURS_BEHAVIOR_PAUSE}
			NewNFTFloorPriceMonitorTask(openSeaService, nftBot, cfg.NFTFloorPriceMonitor.NFTCollections, cfg.NFTFloorPriceMonitor.IntervalSeconds, qh).Start()
		} else {
			log.Printf("Warning: Bot %s not found for NFTFloorPriceMonitorTask", cfg.NFTFloorPriceMonitor.BotName)
		}
	}

	// 4. PolymarketMonitorTask
	if cfg.PolymarketMonitor.IntervalSeconds > 0 {
		polyBot := dingBots[cfg.PolymarketMonitor.BotName]
		if polyBot != nil {
			// Polymarket: Pause during 00:00-07:00
			qh := utils.QuietHoursParams{Enabled: true, StartHour: 0, EndHour: 7, Behavior: constant.QUIET_HOURS_BEHAVIOR_PAUSE}
			NewPolymarketMonitorTask(polyClient, polyBot, cfg.PolymarketMonitor.MarketIDs, cfg.PolymarketMonitor.IntervalSeconds, qh).Start()
		} else {
			log.Printf("Warning: Bot %s not found for PolymarketMonitorTask", cfg.PolymarketMonitor.BotName)
		}
	}

	// 5. TwitterMonitorTask
	if cfg.TwitterMonitor.IntervalSeconds > 0 {
		twitterBot := dingBots[cfg.TwitterMonitor.BotName]
		if twitterBot != nil {
			// Twitter: Pause during 00:00-07:00
			// Example throttling: qh := utils.QuietHoursParams{Enabled: true, StartHour: 0, EndHour: 7, Behavior: constant.QUIET_HOURS_BEHAVIOR_THROTTLE, ThrottleMultiplier: 10}
			qh := utils.QuietHoursParams{Enabled: true, StartHour: 0, EndHour: 7, Behavior: constant.QUIET_HOURS_BEHAVIOR_PAUSE}
			NewTwitterMonitorTask(twitterClient, twitterBot, cfg.TwitterMonitor.Usernames, cfg.TwitterMonitor.IntervalSeconds, qh).Start()
		} else {
			log.Printf("Warning: Bot %s not found for TwitterMonitorTask", cfg.TwitterMonitor.BotName)
		}
	}
}
