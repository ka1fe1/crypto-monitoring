## 背景 (Context)

用户需要每天定时抓取 Polymarket 上指定交易者的统计数据（交易量、排名、PnL、持仓总价值、当前持仓明细），并输出为本地 Markdown 报告文件。输入为一个包含钱包名称和地址的 Markdown 表格文件，输出为按日期命名的 Markdown 表格文档。

项目中已有 `PolymarketClient`（基于 gamma-api），但缺少 data-api 端点的调用能力以及 EOA 到 proxyWallet 的地址解析。

## 目标与非目标

**目标:**
- 实现从 Markdown 表格（`wallet_name | wallet_addr`）中解析交易者地址列表
- 新增 `ResolveProxyWallet` 方法，通过 gamma-api `/public-profile` 将 EOA 地址解析为 Polymarket proxyWallet
- 扩展 `pkg/utils/polymarket`，增加 3 个 data-api 端点的调用方法
- 创建 `PolymarketDailyReportTask` 定时任务（每天一次）
- 在 `config.go` 中自动将相对路径解析为项目根目录下的绝对路径

**非目标:**
- 同步完整历史交易流水
- 集成外部存储系统（如 Notion）

## 架构设计

### 数据流

```
address_list.md (Markdown Table)
    ↓ ParseAddressList() → []WalletEntry
    ↓
PolymarketDailyReportTask.Run()
    ↓ for each entry:
    │   1. ResolveProxyWallet(eoa) → proxyWallet  [gamma-api]
    │   2. GetTraderLeaderboardRankings(proxyWallet) → rank, vol, pnl  [data-api]
    │   3. GetTotalValueOfUserPositions(proxyWallet) → value  [data-api]
    │   4. GetCurrentPositionsForUser(proxyWallet) → []Position  [data-api]
    │   5. 格式化 positions 为无序列表字符串
    │   6. 汇总为 TraderReportData
    ↓
WriteReportTable() → polymarket_volume_YYYYMMDD_HHMM.md
```

### 关键设计决策

1. **双地址体系**: Polymarket 的 data-api 需要 proxyWallet 而非 EOA 地址。因此新增 `ResolveProxyWallet` 步骤，通过 gamma-api `/public-profile?address=` 获取映射关系。

2. **Markdown 表格输入**: 采用 `wallet_name | wallet_addr` 两列表格格式，支持为每个地址赋予可读名称，便于识别。

3. **相对路径解析**: 在 `LoadConfig` 中基于 config 文件路径推算项目根目录，自动将 `./data/...` 等相对路径转为绝对路径，确保 `go test` 和正式启动时路径一致。

4. **持仓明细格式**: Position 数据格式化为无序列表，多个 position 之间用 `<br>` 换行，`|` 字符用 `\|` 转义以避免破坏 Markdown 表格结构。

5. **API 速率控制**: 每个地址处理后 sleep 500ms，避免触发频率限制。

### 涉及文件

| 文件                                             | 作用                                                                                                               |
| ------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------ |
| `config/config.go`                               | 新增 `PolymarketReportConfig`，相对路径自动解析                                                                    |
| `config/config.yaml`                             | 新增 `polymarket_report` 配置块                                                                                    |
| `pkg/utils/markdown/markdown.go`                 | `WalletEntry`、`TraderReportData`、`ParseAddressList`、`WriteReportTable`                                          |
| `pkg/utils/polymarket/polymarket_vo.go`          | `LeaderboardResponse`、`TotalValueResponse`、`Position`、`PublicProfileResponse`                                   |
| `pkg/utils/polymarket/polymarket.go`             | `ResolveProxyWallet`、`GetTraderLeaderboardRankings`、`GetTotalValueOfUserPositions`、`GetCurrentPositionsForUser` |
| `internal/tasks/polymarket_daily_report_task.go` | `PolymarketDailyReportTask` 结构体与 `Run()` 方法                                                                  |
| `internal/tasks/init.go`                         | 注册任务到 cron 调度                                                                                               |

## 风险与权衡

- **风险**: Polymarket API 频率限制或接口变更
  - **缓解**: 每天仅执行一次，每地址间隔 500ms，健壮的错误处理与日志
- **风险**: 输入文件格式异常
  - **缓解**: 正则解析 + 跳过无效行，记录日志
- **风险**: proxyWallet 解析失败（地址未注册 Polymarket）
  - **缓解**: 跳过该地址并记录错误，不影响其他地址处理
