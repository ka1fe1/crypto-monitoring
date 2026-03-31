## Context

当前 Bitcoin 的指标推送面板（`BtcDashboardService`）已经整合了现价、FGI（恐惧贪婪指数）、Ahr999 及 WMA200 等指标，但这些多偏向于技术形态或是市场情绪模型（比如 FGI）。
在评估真正的长线底部价值时，使用底层链上模型的“Balanced Price (均衡价格，BP)”能够带来结构性参考。近期提案 `add-btc-bp-metric` 要求将其接入报告流中。

## Goals / Non-Goals

**Goals:**
- 在不破坏现有指标服务稳定性的同时，集成获取 Bitcoin 的 Balanced Price 数据的逻辑。
- 设计健壮的计算模块以计算（当前价格 / BP）并进行状态分类（高估/正常/低估）。

**Non-Goals:**
- 不打算引入重型的加密节点来自行读取链上区块重算 BP。
- 此次改动不涉及前端展示或界面设计，仅限于后端指标报告（Markdown 报告）补充。

## Decisions

- **数据源获取方案**: 鉴于 Balanced Price 属于复杂的链上衍生指标（基于 UTXO 销毁币日计算的已实现均价 - 转移均价），业内**极少**提供完**全公开且匿名免费**的 API。我们面临如下两种核心决策并需要用户确认（见 Open Questions）：
  1. **申请免费层的 API Key（推荐）**：例如注册 `BGeometrics`、`CryptoQuant` 或是 `Checkonchain` 的开发者免费档位，将这把 API Key 注入到 `config.yaml` 中，以一个标准的 HTTP Provider 的规范稳定集成。
  2. **网页内部接口爬虫（兜底）**：无需 Key 的纯免费方式。针对目标网页（如 `buybitcoinworldwide.com` 或 `lookintobitcoin.com`）解析他们图表加载时的公开 JSON Endpoint。优点是0成本，但存在被反爬或数据接口突然变动的结构性风险。

- **依赖注入规范**: 无论选用何种数据来源，均会抽象为统一的 `BalancePriceProvider` 接口挂载进 `BtcDashboardService`，以利用已有的协程池并置与 Mock 假数据中，使其不受未来供应商变换的底层逻辑影响。
- **并发调用与指标化常量化**: BP 请求将与 Binance/Mempool 的基础查询同层并行。当前价格 / BP 估值的倍率阈值设定在服务头部作常量处理。

## Open Questions

- **[ACTION REQUIRED] 您倾向于使用哪种 BP 数据获取方案？**
  - Option A: 我可以去注册一个 BGeometrics（或类似平台）的 API Key，我们在代码中把它当标准的认证请求。
  - Option B: 先不要配置 Key，写个简单的轻量爬虫或拉取免授权网站的图表后端 Json（即便可能会随着时间不稳定）。

## Risks / Trade-offs

- **Risk: 数据源免费配额/权限受限** 
  - *Trade-off / Mitigation:* 第三方链上服务的 API 频繁度要求较高。我们的 `btc_dashboard_monitor_task` 设计了灵活的获取间距（例如几分钟一次）。可为其增加容错 fallback：若获取 BP 失败或超限，不应该导致整个报告失败，而是在 Markdown 报告中的相关栏位输出 `Unknown / 获取失败`。
- **Risk: 模型延迟性**
  - BP 数据经常有 24 小时的数据滞后更新期，不作为高频实时指标去参考。说明注释或系统提示内也应告知只适于中长线。
