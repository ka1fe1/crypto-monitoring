## Context

项目已有一个 gin Web 服务器（`internal/api/routers/router.go`），提供 REST API 端点。`PolymarketDailyReportTask` 每天生成 Markdown 报告文件到 `data/reports/` 目录，文件名格式为 `polymarket_volume_YYYYMMDD_HHMM.md`。

目前报告只能通过文件系统访问，需要通过 Web 页面直接查看最新报告。

## Goals / Non-Goals

**Goals:**
- 新增 API 端点 `GET /api/v1/polymarket/report`，读取最新报告文件并返回解析后的 JSON 数据
- 新增前端静态页面，通过 API 获取数据并渲染为美观的表格
- 页面标注数据来源文件名和生成时间
- 复用现有 gin 路由框架，无需引入新的前端框架

**Non-Goals:**
- 不实现历史报告列表或分页浏览
- 不实现报告编辑功能
- 不引入前端构建工具（直接使用原生 HTML/CSS/JS）

## Decisions

### 1. API 设计
- 新增 `GET /api/v1/polymarket/report` 端点
- Handler 读取 `config.PolymarketReport.OutputDir` 目录，按文件名排序取最新文件
- 解析 Markdown 表格为结构化 JSON 返回
- 响应包含 `filename`、`generatedAt`、`headers`、`rows` 字段

### 2. 前端页面
- 使用单个 HTML 文件（内嵌 CSS + JS），通过 gin 的 `Static` 或 `StaticFile` 提供服务
- 页面路由: `GET /polymarket/report`
- 使用原生 fetch API 调用后端接口
- 表格样式使用现代深色主题，position 列支持展开查看

### 3. 静态文件位置
- 前端文件放在 `web/static/` 目录下
- gin 配置 `r.Static("/static", "./web/static")` 和 `r.StaticFile("/polymarket/report", "./web/static/polymarket_report.html")`

### 4. Markdown 解析
- 在 handler 中解析 Markdown 表格：按 `|` 分割每行，提取表头和数据行
- 复用已有的 `pkg/utils/markdown` 包或在 handler 中直接做轻量解析

## Risks / Trade-offs

- **风险**: 报告文件不存在（首次部署或未执行过 task）
  - **缓解**: API 返回 404 + 友好提示，前端显示"暂无报告"
- **风险**: 报告文件过大（地址数量多，position 明细长）
  - **缓解**: current_position 列默认折叠，点击展开
- **风险**: 静态文件路径在 Docker 部署时可能不正确
  - **缓解**: 在 `LoadConfig` 中统一处理相对路径解析
