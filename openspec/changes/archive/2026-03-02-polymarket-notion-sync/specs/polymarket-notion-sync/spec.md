## 功能规格: Polymarket 交易数据每日报告 (Polymarket Trader Data Daily Report)

### 概述

系统每天定时从本地 Markdown 表格文件中读取交易者钱包名称与地址列表，通过 Polymarket gamma-api 将 EOA 地址解析为 proxyWallet，再通过 data-api 获取排行榜交易量/PnL、持仓总价值和当前持仓明细，最终将汇总数据以 Markdown 表格形式保存到按日期命名的本地文件中。

### 输入格式

输入文件为 Markdown 表格，包含 `wallet_name` 和 `wallet_addr` 两列：

```markdown
| wallet_name | wallet_addr   |
| ----------- | ------------- |
| alice       | 0x1234...abcd |
| bob         | 0x5678...efgh |
```

### 输出格式

输出文件名格式为 `polymarket_volume_YYYYMMDD_HHMM.md`，表格列为：

`| wallet_addr | wallet_name | proxy_addr | total_volume | vol_rank | total_pnl | value | current_position |`

其中 `current_position` 列以无序列表格式展示每个持仓：
`- {title} \| {outcome} \| init: {avgPrice}({initialValue}) \| current: {curPrice}({currentValue}) \| cash pnl: {cashPnl} \| to win: {size} \| redeemable: {bool}`

### 涉及的 Polymarket API

| 用途                      | 端点                                          | 参数                                             |
| ------------------------- | --------------------------------------------- | ------------------------------------------------ |
| EOA → proxyWallet 解析    | `GET gamma-api.polymarket.com/public-profile` | `?address={eoa}`                                 |
| 排行榜（交易量/排名/PnL） | `GET data-api.polymarket.com/v1/leaderboard`  | `?user={proxyWallet}&timePeriod=ALL&orderBy=VOL` |
| 持仓总价值                | `GET data-api.polymarket.com/value`           | `?user={proxyWallet}`                            |
| 当前持仓明细              | `GET data-api.polymarket.com/positions`       | `?user={proxyWallet}`                            |

### 场景

#### 场景: 成功生成每日交易报告
- **当** 每天定时的 `PolymarketDailyReportTask` 执行时
- **那么** 从配置的本地 Markdown 表格文件解析出所有 `WalletEntry`（名称 + 地址）
- **那么** 对每个地址先调用 `ResolveProxyWallet` 将 EOA 解析为 proxyWallet
- **那么** 用 proxyWallet 依次请求 leaderboard、value、positions 三个 data-api 端点
- **那么** 汇总数据写入新的 Markdown 文件

#### 场景: 缺少配置
- **当** 系统初始化任务时，若 `config.yaml` 中缺少 `address_list_file` 或 `output_dir`
- **那么** 记录错误日志并中止

#### 场景: 地址解析失败
- **当** `ResolveProxyWallet` 无法解析某地址时
- **那么** 记录错误日志，跳过该地址，继续处理下一个

#### 场景: API 调用失败
- **当** 某个 data-api 请求返回非 200 或超时
- **那么** 记录该地址的错误日志，跳过该地址，继续处理下一个

### 配置

```yaml
polymarket_report:
  address_list_file: "./data/address_list.md"
  output_dir: "./data/reports"
```

相对路径在 `LoadConfig` 中自动解析为相对于项目根目录的绝对路径。
