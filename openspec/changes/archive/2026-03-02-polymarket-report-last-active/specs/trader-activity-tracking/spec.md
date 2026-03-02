## 新增需求 (ADDED Requirements)

### 需求: 获取最新活动 (Fetch Latest Activity)
系统必须能够获取给定 Polymarket 交易者地址的最近活动。

#### 场景: 成功获取活动
- **当 (WHEN)** 提供有效的代理钱包地址给 `GetUserActivity` 时。
- **那么 (THEN)** API 返回活动列表（JSON 数组），系统提取第一条（最近）记录的 `timestamp`（`int64` 类型，Unix 秒级时间戳）。

#### 场景: 未找到活动
- **当 (WHEN)** 用户在 API 中没有任何活动记录时。
- **那么 (THEN)** 系统在报告中显示 "N/A"。

### 需求: 格式化最后活跃时间 (Format Last Active Time)
提取的 Unix 时间戳需使用 `utils.FormatRelativeTime` 进行格式化，同时显示 UTC+8 绝对时间和相对时间。

#### 场景: 时间格式化
- **当 (WHEN)** 获取到 Unix 时间戳（如 `1709389200`）时。
- **那么 (THEN)** 系统格式化为 `2024-03-02 18:00:00 (N days ago)` 或 `(N hours ago)` 等相对时间形式。

### 需求: 报告列名 (Column Name)
Markdown 表格和 HTML 前端中的列名统一使用 `last_active`。

### 需求: 前端展示 (Frontend Display)
- `last_active` 列使用 `.time` 样式类渲染，支持点击排序。
- `current_position` 列使用与 `wallet_info` 相同的可折叠卡片样式：
  - 折叠时显示 `▶ N positions` 按钮。
  - 展开后每个仓位以卡片形式展示：头部为 `{title} | {outcome}`，下方按 key-value 逐行显示详情（init、current、cash pnl、to win、redeemable）。
  - 去除 Markdown 转义符 `\\`。
