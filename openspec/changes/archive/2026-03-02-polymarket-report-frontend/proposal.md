## Why

`PolymarketDailyReportTask` 每天生成的 Markdown 报告目前只能通过文件系统查看，不方便远程访问和团队共享。需要一个前端页面，通过 URL 直接访问最新报告内容，并以美观的表格形式展示。

## What Changes

- 新增一个 HTTP API 端点，读取 `data/reports/` 目录下最新的 Markdown 报告文件，解析并返回 JSON 数据
- 新增一个前端 HTML 页面，通过该 API 获取数据并渲染为美观的表格页面
- 页面底部标注数据来源文件名
- 复用现有的 gin 路由框架，将新端点注册到 `SetupRouter`

## Capabilities

### New Capabilities
- `polymarket-report-viewer`: 通过 Web 页面展示最新的 Polymarket 每日报告，包括 API 端点和前端页面

### Modified Capabilities
（无）

## Impact

- `internal/api/routers/router.go`: 注册新路由
- `internal/api/handlers/`: 新增 report handler
- 前端静态文件: 新增 HTML/CSS/JS 页面
- 依赖: 无新外部依赖，复用 gin 框架
