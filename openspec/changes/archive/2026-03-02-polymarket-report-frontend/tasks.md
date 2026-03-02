## 1. 后端 API (Report API)

- [x] 1.1 创建 `internal/api/handlers/polymarket_report_handler.go`，实现 `GetLatestReport` handler
  - 读取 `config.PolymarketReport.OutputDir` 目录
  - 按文件名排序取最新的 `.md` 文件
  - 解析 Markdown 表格（提取表头和每行数据）
  - 返回 JSON: `{ filename, generatedAt, headers, rows }`
  - 无文件时返回 404，解析错误返回 500
- [x] 1.2 在 `internal/api/routers/router.go` 中注册路由 `GET /api/v1/polymarket/report`
- [x] 1.3 将 `config` 传入 handler（用于读取 `OutputDir` 路径）

## 2. 前端页面 (Frontend Page)

- [x] 2.1 创建 `web/static/polymarket_report.html`，单文件包含 HTML + CSS + JS
  - 深色主题，现代表格样式
  - 自动调用 `/api/v1/polymarket/report` 获取数据
  - 渲染 8 列表格（wallet_addr, wallet_name, proxy_addr, total_volume, vol_rank, total_pnl, position_value, last_active, current_position）
  - current_position 列折叠/展开交互
  - 页面底部标注来源文件名和生成时间
  - 无数据时显示友好提示
- [x] 2.3 增加表头点击排序功能（支持数值、金额、文本）
- [x] 2.4 优化持仓数据展开后的内部布局与色彩一致性
- [x] 2.5 合并 wallet_addr, wallet_name, proxy_addr 为 wallet_info 列，支持展开/折叠显示详细信息；折叠按钮显示 wallet_name（金色字体）而非截断地址

## 3. 路径与配置 (Path & Config)

- [x] 3.1 确保 `web/static/` 目录路径在不同环境下（本地开发、Docker）均可正确访问
- [x] 3.2 如需要，在 `config.go` 中为静态文件目录添加路径解析

## 4. 验证 (Verification)

- [x] 4.1 启动服务后访问 `/polymarket/report` 确认页面正常渲染
- [x] 4.2 确认 API `/api/v1/polymarket/report` 返回正确的 JSON 数据
- [x] 4.3 确认无报告文件时页面和 API 均返回友好提示
