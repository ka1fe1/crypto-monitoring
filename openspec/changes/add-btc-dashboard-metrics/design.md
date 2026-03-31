## Context
目前缺乏对比特币相关宏观和链上数据的监控（如 200WMA，ahr999指数，恐慌贪婪指数等）。我们需要实现一个定时任务（`btc_dashboard_monitor_task`）并通过 `btc_dashboard_service` 集中获取、计算这些指标，最后整合判断市场所处的阶段。此系统的开发需遵循已有规范：外部 API 的调用应封装为 `/pkg/utils/xxx` 下的独立工具类。

## Goals / Non-Goals

**Goals:**
- 实现获取 Binance 或 OKX 的长周期日线和周线数据，计算 200 均线（MA）。
- 实现获取 mempool.space 区块高度并推算减半周期。
- 实现获取 Alternative.me 贪婪与恐慌指数。
- 在 `btc_dashboard_service` 层实现复合指标（ahr999、Balanced Price、MVRV等，若缺乏公开链上接口，优先实现 ahr999、WMA200 这些可算出的纯历史数据衍生指标）。
- 整合报警触发逻辑并通过 Dingding bot 发出。

**Non-Goals:**
- 不实现高频即时行情的监控抓取（仅定位于 12~24 小时级别的宏观指标报警）。
- 不自建全节点的链上数据解析，纯依赖成熟外部轻量 API。

## Decisions
1. **封装独立工具类 (Utils):** 
   - `pkg/utils/binance`: 获取 `klines` (日线、周线) 用于计算 200DMA / 200WMA。
   - `pkg/utils/alternative`: 获取贪婪恐慌指数 `fng/?limit=1`。
   - `pkg/utils/mempool`: 获取 `blocks/tip/height` 顶层区块高度用于估算减半时间。
   所有工具类配备对应的单元测试 `xxx_test.go` 与结构体 `xxx_vo.go`。
2. **业务集成 (Service & Task):**
   - 创建 `btc_dashboard_service` 来调用各工具类并进行聚合计算（如 `ahr999` 指数公式计算），同时生成监控消息文本。
   - 创建 `btc_dashboard_monitor_task` 负责按照定时策略触发 Service 以及调用 DingBot 工具。

## Risks / Trade-offs

- **Risk: 免费公共 API 的限流被封问题** ->  Mitigation: 设定为长线宏观监控指标，每日或每 12 小时拉取一次。Binance 和 mempool 的公共限流远超我们的调用频次要求。
- **Risk: 复杂链上指标 (MVRV / 均衡价格) 免费 API 难寻** -> Mitigation: 优先对接 Glassnode 免费版或查阅是否有公开免费平替，若无，则第一期仅计算能用 K 线数据直接推算的（如 ahr999 及 200WMA），同时将依赖链上 UTXO 解析的 MVRV/BP 预留后期接口。
