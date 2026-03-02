## ADDED Requirements

### Requirement: 最新报告 API 端点
系统应当 (SHALL) 提供 `GET /api/v1/polymarket/report` 端点，读取配置的输出目录中最新的 Markdown 报告文件，解析其中的表格数据，以 JSON 格式返回。

#### Scenario: 成功获取最新报告
- **WHEN** 客户端请求 `GET /api/v1/polymarket/report`
- **THEN** 系统读取 `output_dir` 中按文件名排序最新的 `.md` 文件
- **THEN** 解析 Markdown 表格为结构化数据
- **THEN** 返回 200 状态码及 JSON 响应，包含 `filename`、`generatedAt`、`headers`、`rows` 字段

#### Scenario: 无报告文件
- **WHEN** 客户端请求 `GET /api/v1/polymarket/report`，但输出目录中没有报告文件
- **THEN** 返回 404 状态码及 `{"error": "No report files found"}` 响应

#### Scenario: 文件读取失败
- **WHEN** 系统在读取或解析报告文件时发生错误
- **THEN** 返回 500 状态码及包含错误信息的 JSON 响应

### Requirement: 报告展示前端页面
系统应当 (SHALL) 提供 `GET /polymarket/report` 路由，返回一个美观的 HTML 页面，自动从 API 获取最新报告数据并渲染为表格。

#### Scenario: 正常展示报告
- **WHEN** 用户通过浏览器访问 `/polymarket/report`
- **THEN** 页面加载后自动调用 `/api/v1/polymarket/report` API
- **THEN** 将返回的表格数据渲染为美观的深色主题表格
- **THEN** 页面底部标注"数据来源: {filename}"及生成时间

#### Scenario: 暂无报告
- **WHEN** 用户访问页面但 API 返回 404
- **THEN** 页面显示友好提示"暂无报告数据，请等待每日定时任务执行"

#### Scenario: 持仓明细展开
- **WHEN** 报告中的 `current_position` 列包含多条持仓记录
- **THEN** 该列默认以折叠形式展示，点击可展开查看所有持仓明细
