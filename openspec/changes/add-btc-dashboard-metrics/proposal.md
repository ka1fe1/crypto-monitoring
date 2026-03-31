## Why
目前我们缺少对 BTC 相关的宏观周期及链上数据指标的监控。如 https://btc.x.fish/ 看板展示的一样，200WMA（200周均线）、恐惧贪婪指数、ahr999 定投指数、减半倒计时、均衡价格（Balanced Price）及 MVRV 等宏观指标对把握市场底顶非常关键。将其集成到我们的监控模块中能够丰富我们的市场异动预警能力。

## What Changes
- 增加对接外部行情数据源（Binance/OKX）的 K 线接口，获取底层资产的数据来计算 200 周和 200 日均线系列。
- 增加对接 Alternative.me 的公开 API，获取比特币的恐惧与贪婪指数。
- 增加对接 mempool.space 的 API，拉取当前区块高度，从而计算减半倒计时。
- 实现相关的统计算法（ahr999指数计算、200WMA 偏离度等）。
- 创建对应的 `btc_dashboard_monitor_task` 与 `btc_dashboard_service`，将逻辑规范化并接入调度。
- 支持 Dingding Bot 报警推送等配置关联。

## Capabilities

### New Capabilities
- `btc-dashboard-metrics`: 获取和计算比特币相关的核心宏观周期和链上数据指标。

### Modified Capabilities

## Impact
- 新增了数个外部公共 API 的网络依赖（无需 Token 的公共 API）。
- 增加服务层（Service）和定时监控层（Task）的资源消耗（网络请求与少量公式运算）。
