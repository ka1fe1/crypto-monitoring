## 1. 工具层 (pkg/utils/polymarket)

- [x] 1.1 在 `polymarket_vo.go` 中添加 `Activity`（`Timestamp` 为 `int64`）及 `ActivityResponse` 结构体。
- [x] 1.2 在 `polymarket.go` 中实现 `GetUserActivity(address)`，请求 `data-api.polymarket.com/activity`。
- [x] 1.3 在 `polymarket_test.go` 中为 `GetUserActivity` 添加单元测试，验证时间戳解析和 UTC+8 转换。

## 2. 任务层 (internal/tasks)

- [x] 2.1 在 `markdown.go` 的 `TraderReportData` 结构体中增加 `LastActiveTime` 字段。
- [x] 2.2 在 `polymarket_daily_report_task.go` 中调用 `GetUserActivity` 并提取最新时间戳。
- [x] 2.3 使用 `utils.FormatRelativeTime(time.Unix(ts, 0))` 格式化为 UTC+8 绝对时间 + 相对时间。
- [x] 2.4 更新 Markdown 报告表头和行格式，增加 `last_active` 列。

## 3. 后端解析层 (internal/api/handlers)

- [x] 3.1 更新 `polymarket_report_handler.go` 中 `parsePipeLine` 的列分割索引（`i < 7` → `i < 8`），适配第 9 列 `current_position`。

## 4. 前端层 (web/static)

- [x] 4.1 在 `polymarket_report.html` 中增加 `.time` CSS 类，用于 `last_active` 列的渲染。
- [x] 4.2 在 `formatCell` 中处理 `last_active` 列的格式化。
- [x] 4.3 重构 `formatPositions` 和相关 CSS，使用与 `wallet_info` 一致的可折叠卡片样式。
- [x] 4.4 去除 `current_position` HTML 渲染中的 Markdown 转义符 `\\`。
- [x] 4.5 确保 `last_active` 列支持排序。

## 5. 通用工具 (pkg/utils)

- [x] 5.1 更新 `FormatRelativeTime` 中超过 24 小时的分支，显示为 `{bjTime} ({N} days ago)`。

## 6. 后续优化 (Follow-up)

- [x] 6.1 将 `TraderReportData.Value` 字段重命名为 `PositionValue`，Markdown 报告列名由 `value` 更新为 `position_value`，HTML `formatCell` 中对应判断条件同步更新。
- [x] 6.2 `wallet_info` 折叠按钮改为展示 `wallet_name`（金色字体），与展开后的 Wallet Name 标签保持一致，提升可读性。
