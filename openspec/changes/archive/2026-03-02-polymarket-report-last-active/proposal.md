## 动机 (Why)

目前 Polymarket 报告显示了交易额、排名、PnL 和价值，但无法体现交易者的活跃程度。添加“最后活跃时间”可以帮助用户区分活跃交易者和静默交易者，为报告数据提供更丰富的背景。

## 改动内容 (What Changes)

- 更新 `PolymarketDailyReportTask`，为每个地址获取其最近的活动记录。
- 提取最新的活动时间戳并转换为 UTC+8 时区。
- 在生成的 Markdown 报告文字中增加“最后活跃时间”列。
- 更新 HTML 前端页面，支持该列的展示和排序。

## 能力 (Capabilities)

### 新增能力
- `trader-activity-tracking`: 从 Polymarket Data API 获取并处理交易者的最近活动记录。

### 修改能力
- `polymarket-daily-report`: 每日报告生成流程新增了一个维度（最后活跃时间）。

## 影响 (Impact)

- **后端**: 需要更新 `PolymarketReportHandler` 和 `PolymarketDailyReportTask`。
- **API**: `pkg/utils/polymarket` 需要新增获取用户活动的方法。
- **前端**: `polymarket_report.html` 需要渲染新列并支持排序。
- **数据**: 旧报告将不含此列，新生成的报告将包含该数据。
