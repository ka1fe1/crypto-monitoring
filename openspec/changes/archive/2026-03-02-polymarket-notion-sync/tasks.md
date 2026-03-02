## 1. 配置设置 (Configuration Setup)

- [x] 1.1 在 `config/config.go` 中添加 `PolymarketReportConfig` 结构体（`AddressListFile` + `OutputDir`）
- [x] 1.2 在 `LoadConfig` 中添加相对路径自动解析为项目根目录绝对路径的逻辑
- [x] 1.3 更新 `config.yaml`，加入 `polymarket_report` 配置块

## 2. Markdown 文件读写工具 (Markdown I/O Utilities)

- [x] 2.1 创建 `pkg/utils/markdown` 包
- [x] 2.2 定义 `WalletEntry` 结构体（Name + Address）
- [x] 2.3 实现 `ParseAddressList` 方法，解析 Markdown 表格（`wallet_name | wallet_addr`）为 `[]WalletEntry`
- [x] 2.4 定义 `TraderReportData` 结构体（WalletName, Address, ProxyAddr, Volume, Rank, Pnl, Value, CurrentPositions）
- [x] 2.5 实现 `WriteReportTable` 方法，生成 8 列 Markdown 表格并按日期时间命名保存
- [x] 2.6 编写 `ParseAddressList` 和 `WriteReportTable` 的单元测试

## 3. Polymarket API 扩展 (Polymarket API Extension)

- [x] 3.1 定义 VO 结构体：`LeaderboardResponse`（Rank, ProxyWallet, Vol, Pnl）、`TotalValueResponse`（Value）、`Position`（Title, Outcome, AvgPrice, CurPrice, InitialValue, CurrentValue, Size, CashPnl, Redeemable）、`PublicProfileResponse`（ProxyWallet）
- [x] 3.2 实现 `ResolveProxyWallet` 方法（gamma-api `/public-profile`，EOA → proxyWallet）
- [x] 3.3 实现 `GetTraderLeaderboardRankings` 方法（data-api `/v1/leaderboard`，参数 `timePeriod=ALL&orderBy=VOL`）
- [x] 3.4 实现 `GetTotalValueOfUserPositions` 方法（data-api `/value`）
- [x] 3.5 实现 `GetCurrentPositionsForUser` 方法（data-api `/positions`）
- [x] 3.6 编写以上 4 个方法的单元测试

## 4. 任务编排 (Task Orchestration)

- [x] 4.1 创建 `internal/tasks/polymarket_daily_report_task.go`
- [x] 4.2 实现 `PolymarketDailyReportTask` 结构体与 `Run()` 核心逻辑：读取地址表 → ResolveProxyWallet → 遍历调用 3 个 data-api → 格式化 positions 为无序列表 → 汇总写入 Markdown
- [x] 4.3 在 `internal/tasks/init.go` 中注册为每天执行一次的定时任务
- [x] 4.4 编写 `PolymarketDailyReportTask` 单元测试

## 5. 数据文件 (Data Files)

- [x] 5.1 创建 `data/address_list.md` 示例输入文件（Markdown 表格格式）
- [x] 5.2 确保 `data/reports/` 输出目录在任务运行时自动创建
