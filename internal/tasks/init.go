package tasks

import (
	"github.com/ka1fe1/crypto-monitoring/config"
	"github.com/ka1fe1/crypto-monitoring/internal/service"
	"github.com/ka1fe1/crypto-monitoring/pkg/logger"
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
	// Create services
	twitterMonitorService := service.NewTwitterService(twitterClient)
	polymarketService := service.NewPolymarketMonitorService(polyClient)

	// 1. DexPairAlterTask
	if cfg.DexPairAlter.IntervalSeconds > 0 {
		dexBot := dingBots[cfg.DexPairAlter.BotName]
		if dexBot != nil {
			var qh utils.QuietHoursParams
			if cfg.DexPairAlter.QuietHours != nil {
				qh = utils.QuietHoursParams{
					Enabled:            cfg.DexPairAlter.QuietHours.Enabled,
					StartHour:          cfg.DexPairAlter.QuietHours.StartHour,
					EndHour:            cfg.DexPairAlter.QuietHours.EndHour,
					Behavior:           cfg.DexPairAlter.QuietHours.Behavior,
					ThrottleMultiplier: cfg.DexPairAlter.QuietHours.ThrottleMultiplier,
				}
			} else {
				// Default: Pause during 00:00-08:00
				qh = utils.QuietHoursParams{Enabled: true, StartHour: 0, EndHour: 8, Behavior: constant.QUIET_HOURS_BEHAVIOR_PAUSE}
			}
			NewDexPairAlterTask(dexService, dexBot, cfg.DexPairAlter.ContractAddrInfo, cfg.DexPairAlter.IntervalSeconds, qh).Start()
		} else {
			logger.Warn("Warning: Bot %s not found for DexPairAlterTask", cfg.DexPairAlter.BotName)
		}
	}

	// 2. TokenPriceMonitorTask
	if cfg.TokenPriceMonitor.IntervalSeconds > 0 {
		tokenBot := dingBots[cfg.TokenPriceMonitor.BotName]
		if tokenBot != nil {
			var qh utils.QuietHoursParams
			if cfg.TokenPriceMonitor.QuietHours != nil {
				qh = utils.QuietHoursParams{
					Enabled:            cfg.TokenPriceMonitor.QuietHours.Enabled,
					StartHour:          cfg.TokenPriceMonitor.QuietHours.StartHour,
					EndHour:            cfg.TokenPriceMonitor.QuietHours.EndHour,
					Behavior:           cfg.TokenPriceMonitor.QuietHours.Behavior,
					ThrottleMultiplier: cfg.TokenPriceMonitor.QuietHours.ThrottleMultiplier,
				}
			} else {
				// Default: Pause during 00:00-08:00, Throttle
				qh = utils.QuietHoursParams{Enabled: true, StartHour: 0, EndHour: 8, Behavior: constant.QUIET_HOURS_BEHAVIOR_THROTTLE, ThrottleMultiplier: 5}
			}
			NewTokenPriceMonitorTask(tokenService, tokenBot, cfg.TokenPriceMonitor.TokenIds, cfg.TokenPriceMonitor.IntervalSeconds, qh).Start()
		} else {
			logger.Warn("Warning: Bot %s not found for TokenPriceMonitorTask", cfg.TokenPriceMonitor.BotName)
		}
	}

	// 3. NFTFloorPriceMonitorTask
	if cfg.NFTFloorPriceMonitor.IntervalSeconds > 0 {
		nftBot := dingBots[cfg.NFTFloorPriceMonitor.BotName]
		if nftBot != nil {
			var qh utils.QuietHoursParams
			if cfg.NFTFloorPriceMonitor.QuietHours != nil {
				qh = utils.QuietHoursParams{
					Enabled:            cfg.NFTFloorPriceMonitor.QuietHours.Enabled,
					StartHour:          cfg.NFTFloorPriceMonitor.QuietHours.StartHour,
					EndHour:            cfg.NFTFloorPriceMonitor.QuietHours.EndHour,
					Behavior:           cfg.NFTFloorPriceMonitor.QuietHours.Behavior,
					ThrottleMultiplier: cfg.NFTFloorPriceMonitor.QuietHours.ThrottleMultiplier,
				}
			} else {
				// Default: Pause during 00:00-08:00
				qh = utils.QuietHoursParams{Enabled: true, StartHour: 0, EndHour: 8, Behavior: constant.QUIET_HOURS_BEHAVIOR_PAUSE}
			}
			NewNFTFloorPriceMonitorTask(openSeaService, nftBot, cfg.NFTFloorPriceMonitor.NFTCollections, cfg.NFTFloorPriceMonitor.IntervalSeconds, qh).Start()
		} else {
			logger.Warn("Warning: Bot %s not found for NFTFloorPriceMonitorTask", cfg.NFTFloorPriceMonitor.BotName)
		}
	}

	// 4. PolymarketMonitorTask
	if cfg.PolymarketMonitor.IntervalSeconds > 0 {
		polyBot := dingBots[cfg.PolymarketMonitor.BotName]
		if polyBot != nil {
			var qh utils.QuietHoursParams
			if cfg.PolymarketMonitor.QuietHours != nil {
				qh = utils.QuietHoursParams{
					Enabled:            cfg.PolymarketMonitor.QuietHours.Enabled,
					StartHour:          cfg.PolymarketMonitor.QuietHours.StartHour,
					EndHour:            cfg.PolymarketMonitor.QuietHours.EndHour,
					Behavior:           cfg.PolymarketMonitor.QuietHours.Behavior,
					ThrottleMultiplier: cfg.PolymarketMonitor.QuietHours.ThrottleMultiplier,
				}
			} else {
				// Default: Pause during 00:00-08:00
				qh = utils.QuietHoursParams{Enabled: true, StartHour: 0, EndHour: 8, Behavior: constant.QUIET_HOURS_BEHAVIOR_PAUSE}
			}
			NewPolymarketMonitorTask(polymarketService, polyBot, cfg.PolymarketMonitor.MarketIDs, cfg.PolymarketMonitor.IntervalSeconds, qh).Start()
		} else {
			logger.Warn("Warning: Bot %s not found for PolymarketMonitorTask", cfg.PolymarketMonitor.BotName)
		}
	}

	// 5. TwitterMonitorTask
	if cfg.TwitterMonitor.IntervalSeconds > 0 {
		twitterBot := dingBots[cfg.TwitterMonitor.BotName]
		if twitterBot != nil {
			var qh utils.QuietHoursParams
			if cfg.TwitterMonitor.QuietHours != nil {
				qh = utils.QuietHoursParams{
					Enabled:            cfg.TwitterMonitor.QuietHours.Enabled,
					StartHour:          cfg.TwitterMonitor.QuietHours.StartHour,
					EndHour:            cfg.TwitterMonitor.QuietHours.EndHour,
					Behavior:           cfg.TwitterMonitor.QuietHours.Behavior,
					ThrottleMultiplier: cfg.TwitterMonitor.QuietHours.ThrottleMultiplier,
				}
			} else {
				// Default: Pause during 00:00-07:00
				qh = utils.QuietHoursParams{Enabled: true, StartHour: 0, EndHour: 7, Behavior: constant.QUIET_HOURS_BEHAVIOR_PAUSE}
			}
			NewTwitterMonitorTask(twitterMonitorService, twitterBot, cfg.TwitterMonitor.Usernames, cfg.TwitterMonitor.Keywords, cfg.TwitterMonitor.IntervalSeconds, qh).Start()
		} else {
			logger.Warn("Warning: Bot %s not found for TwitterMonitorTask", cfg.TwitterMonitor.BotName)
		}
	}

	// 6. GeneralMonitorTask
	if cfg.GeneralMonitor.IntervalSeconds > 0 {
		generalBot := dingBots[cfg.GeneralMonitor.BotName]
		if generalBot != nil {
			var qh utils.QuietHoursParams
			if cfg.GeneralMonitor.QuietHours != nil {
				qh = utils.QuietHoursParams{
					Enabled:            cfg.GeneralMonitor.QuietHours.Enabled,
					StartHour:          cfg.GeneralMonitor.QuietHours.StartHour,
					EndHour:            cfg.GeneralMonitor.QuietHours.EndHour,
					Behavior:           cfg.GeneralMonitor.QuietHours.Behavior,
					ThrottleMultiplier: cfg.GeneralMonitor.QuietHours.ThrottleMultiplier,
				}
			} else {
				// Default: Pause during 00:00-07:00
				qh = utils.QuietHoursParams{Enabled: true, StartHour: 0, EndHour: 8, Behavior: constant.QUIET_HOURS_BEHAVIOR_THROTTLE, ThrottleMultiplier: 5}
			}

			NewGeneralMonitorTask(
				tokenService,
				polymarketService,
				generalBot,
				cfg.GeneralMonitor.Modules,
				cfg.TokenPriceMonitor.TokenIDs,
				cfg.PolymarketMonitor.MarketIDs,
				cfg.GeneralMonitor.IntervalSeconds,
				qh,
			).Start()
		} else {
			logger.Warn("Warning: Bot %s not found for GeneralMonitorTask", cfg.GeneralMonitor.BotName)
		}
	}
}
