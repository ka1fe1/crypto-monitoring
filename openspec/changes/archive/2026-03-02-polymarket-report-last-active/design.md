## 上下文 (Context)

目前的 `PolymarketDailyReportTask` 聚合了多个端点（排行榜、价值、仓位）的数据，但缺乏对用户最近活动时间的感知。Polymarket 提供了 `data-api.polymarket.com/activity` 端点，列出了用户的交易和其他操作及其时间戳。

## 目标 / 非目标 (Goals / Non-Goals)

**目标:**
- 在 Polymarket 工具类中实现 `GetUserActivity`。
- 将活动获取逻辑集成到每日报告任务中。
- 使用 `utils.FormatRelativeTime` 显示 UTC+8 时间及相对时间。
- 确保 HTML 前端可以按此时间进行排序。
- 重构 `current_position` 列，使用与 `wallet_info` 一致的卡片式折叠展示。

**非目标:**
- 跟踪除最近一次活动以外的历史记录。
- 实时活动监控（这是一个每日报告）。

## 决策 (Decisions)

- **工具类 API**: 在 `pkg/utils/polymarket/polymarket.go` 中添加 `GetUserActivity(address)`。请求 `https://data-api.polymarket.com/activity?user={address}`。
- **数据结构**: `Activity.Timestamp` 类型为 `int64`（Unix 秒级时间戳），非 ISO8601 字符串。
- **时间格式**: 使用 `utils.FormatRelativeTime(time.Unix(ts, 0))` 格式化，输出形如 `2026-02-06 23:09:49 (23 days ago)`。
- **Markdown 表格**: 在 `value` 和 `current_position` 之间增加名为 `last_active` 的新列。
- **后端解析**: `parsePipeLine` 中列分割索引从 7 增加到 8，确保第 9 列 `current_position` 正确处理转义管道符。
- **前端样式**: `current_position` 重构为可折叠卡片，每个仓位展示为 `{title} | {outcome}` 头部 + key-value 详情行。HTML 渲染时去除 Markdown 转义符 `\\`。
- **wallet_info 折叠按钮**: 折叠按钮默认展示 `wallet_name`（金色加粗字体），展开后可查看完整 Wallet Addr 与 Proxy Addr 详情，使界面更直观易读。
- **列名变更**: 报告中 `value` 列重命名为 `position_value`，对应 `TraderReportData.PositionValue` 字段，更准确地描述仓位市值语义。

## 风险 / 权衡 (Risks / Trade-offs)

- **API 速率限制**: 为每个地址获取活动会增加一个额外的 API 调用。鉴于地址列表目前较小，这是可以接受的。
- **无活动**: 某些用户可能在 API 中没有任何活动记录。在这种情况下显示 "N/A"。
