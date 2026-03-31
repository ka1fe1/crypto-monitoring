## Why

添加比特币的 "均衡价格（Balanced Price，简称 BP）" 指标，这是一个由 David Puell 创建的核心链上分析指标（等于已实现价格减去转让价格）。
该指标能够有效地识别比特币在熊市中的潜在底部区间。当当前市场价格触及或低于均衡价格时，通常被视为极佳的建仓机会（正常估值区间甚至是低估区间）。引入这一指标到我们现有的监控面板可以显著提升我们对比特币长期估值周期的研判能力，为交易决策提供更为专业和理性的数据支持。

## What Changes

- 在现有的 BTC 监控体系中新加入 `Balanced Price (BP)` 数据的自动获取。
- 因为常见的免费 API 不提供直接的链上数据，我们将集成专门提供链上数据的服务（如 Coinglass API 或其他链上衍生数据来源）来获取每天实时更新的 BP 数据。
- 增加当前价格与 BP 的比率（当前价格 / BP）的计算和评估区间判断（例如：> 1.0 时为正常区间，低于 1.0 为低估区间等），并将这些信息最终组装到 Markdown 报告中并通过钉钉推送给用户。

## Capabilities

### New Capabilities
- `fetch-bp-data`: 从外部数据源（如 Coinglass API）获取比特币的 Balanced Price 数据。
- `bp-valuation-logic`: 基于当前价格 / BP 的数值进行估值判断逻辑计算（判断状态是低估、正常等）。

### Modified Capabilities
- `generate-metric-report`: 更新生成 Markdown 报告的能力，将 BP 及其比值、估值区间加入推送内容的内容流。

## Impact

- **系统层面**: 将引入一个新的外部服务 API（或更新现有的 Provider）来专门获取 BP 数据。需要新增 API 的请求配额监控、容错和缓存。
- **配置层面**: 会在 `config.yaml` 或其对应的配置数据结构中新增有关新数据源（如果是 Coinglass 则需对应 API Key）的设置。
- **业务层面**: 每日生成的报告将变得更为丰富，依赖的第三方源增加，核心抓取逻辑 `btc_dashboard_service.go` 需要拓展对并发生态的处理。
