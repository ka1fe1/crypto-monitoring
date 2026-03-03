## 背景
`PolymarketDailyReportTask` 目前会在任意时间执行，无论当前是否是夜间，这与 `NFTFloorPriceMonitorTask` 遵循了 `quietHoursParams` 不同。为了避免在深夜通过报告打扰用户，需要将该任务更新为严格遵守配置的防打扰静默时间。

## 目标与非目标

**目标:**
- 当当前时间处于 `quietHoursParams` 规定的防打扰时间段内时，跳过 `PolymarketDailyReportTask` 的执行。

**非目标:**
- 不修改防打扰时间配置的数据结构。
- 不影响其他任务的现有功能。

## 设计决策
- 将 `quietHoursParams`（防打扰参数）作为参数传入 `NewPolymarketDailyReportTask`。
- 在 `PolymarketDailyReportTask` 结构体中添加 `lastRunTime` (类型为 time.Time) 和 `interval` (类型为 time.Duration) 字段，并利用 `utils.ShouldExecTask(t.quietHoursParams, t.lastRunTime, t.interval)` 方法来进行时间判断（这与 NFT 地板价监控任务的做法保持一致）。
- 在注册和初始化该任务的位置，提供所需的配置参数。

## 风险与权衡
- 如果配置的触发时间与防打扰时间经常重合，报告可能会被延迟。但由于这是每日报告，用户通常会将其安排在清醒的特定时间，因此这不会造成实际问题。
