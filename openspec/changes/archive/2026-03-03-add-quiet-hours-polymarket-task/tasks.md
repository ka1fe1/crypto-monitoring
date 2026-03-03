## 1. 修改 PolymarketDailyReportTask
- [x] 1.1 在 `PolymarketDailyReportTask` 结构体中添加 `quietHoursParams` (类型为 `utils.QuietHoursParams`)、`interval` (类型为 `time.Duration`) 和 `lastRunTime` (类型为 `time.Time`) 字段。
- [x] 1.2 更新构造函数 `NewPolymarketDailyReportTask` 以接收并初始化这些新字段。
- [x] 1.3 更新 `Run()` 方法：添加条件判断，如果 `!utils.ShouldExecTask(t.quietHoursParams, t.lastRunTime, t.interval)` 为真则跳过执行；同时在任务实际运行时更新 `lastRunTime`。

## 2. 更新任务初始化逻辑
- [x] 2.1 找到调用 `NewPolymarketDailyReportTask` 的位置，并修改函数实参，从应用程序配置中传入正确的静默时间参数（`quietHoursParams`）和间隔时间（`interval`）。
