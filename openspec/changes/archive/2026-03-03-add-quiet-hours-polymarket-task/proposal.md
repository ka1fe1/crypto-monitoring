## 为什么

当前 `PolymarketDailyReportTask` 在执行时没有防打扰静默时间（quiet hours）的控制。这可能会导致在非工作时间或夜间生成报告及警报，从而打扰用户。为了和其他任务（例如 NFT 地板价监控任务）的逻辑保持一致，并防止不必要的非工作时间执行，我们需要添加相应的控制逻辑。

## 更改内容

- 在 `PolymarketDailyReportTask` 中添加 `quietHoursParams` 字段，并更新其构造函数 `NewPolymarketDailyReportTask`。
- 添加执行时间间隔控制，通过 `utils.ShouldExecTask` 判断是否在静默时间内。
- 更新相关的 `Run()` 方法，当处于静默时间内时跳过执行。
- 更新调用 `NewPolymarketDailyReportTask` 进行初始化的代码，将相关的静默时间配置传递给任务。

## 功能 (Capabilities)

### 新增功能
- `polymarket-quiet-hours`: 为 Polymarket 每日报告任务添加静默时间（防打扰机制）支持。

### 修改后的功能
（无）

## 影响范围

- 修改 `internal/tasks/polymarket_daily_report_task.go` 文件。
- 修改任务初始化逻辑中调用 `NewPolymarketDailyReportTask` 的地方以传递新配置。
