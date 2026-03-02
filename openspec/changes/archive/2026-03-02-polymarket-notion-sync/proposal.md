## 为什么 (Why)

用户希望将在 Polymarket 上的交易者统计数据（特别是指定地址在所有时间段和所有类别下的总交易量）每天定时获取并记录。该自动化任务将读取预设的地址列表，并将跑批的结果以 Markdown 表格的形式保存到每天新建的文档中，以时间和日期区分。这样无需依赖外部系统即可随时查看每日切片数据。

## 有哪些改动 (What Changes)

- 创建一个新的后台监控任务（如 `PolymarketDailyReportTask`），该任务每天触发一次。
- 从指定的本地 Markdown 静态文件中读取**交易者地址列表**。
- 遍历地址列表，集成三个 Polymarket API 接口以获取全方位数据：
  1. 获取排行榜上的总交易量 (`/api-reference/core/get-trader-leaderboard-rankings`)
  2. 获取用户当前持仓总价值 (`/api-reference/core/get-total-value-of-a-users-positions`)
  3. 获取用户当前的详细持仓明细 (`/api-reference/core/get-current-positions-for-a-user`)
- 将获取到的所有交易者的数据（地址、总交易量、总持仓价值、当前持仓明细）格式化拼接到一行的 Markdown 表格中，并写入到以当天日期命名的新 Markdown 文件中。
- 在 `config.yaml` 中添加源地址列表文件的路径、以及每日结果报告保存目录的配置。

## 能力 (Capabilities)

### 新增能力 (New Capabilities)
- `polymarket-volume-report`: 每天定时从 Polymarket 获取交易者的统计数据，并将总交易量输出为本地的每日 Markdown 报告。

### 修改的能力 (Modified Capabilities)
- 

## 影响面 (Impact)

- **配置 (Configuration)**: `config.yaml` 将需要新增报告输出目录及输入地址文件的路径配置。
- **任务 (Tasks)**: 在 `internal/tasks` 包中将新增一个基于 cron 调度的后台任务。
- **工具/客户端 (Utils/Clients)**: 现有的 Polymarket 客户端需要增加一个新方法来调用排行榜 API。需要新增读写本地 Markdown 文件的工具方法。
