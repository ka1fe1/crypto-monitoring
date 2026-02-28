## ADDED Requirements

### Requirement: OutcomePrices 按 name 降序排序显示

系统在格式化市场数据时，SHALL 将 `OutcomePrices` 的 key 按字母降序排列，然后依次插入 `prices` 数组。

#### Scenario: 多个 outcome 按降序排列
- **WHEN** `OutcomePrices` 包含 keys `["Yes", "No"]`
- **THEN** `prices` 数组中先展示 `Yes` 的价格，再展示 `No` 的价格

#### Scenario: 单个 outcome 保持正常
- **WHEN** `OutcomePrices` 只包含一个 key
- **THEN** `prices` 数组仅包含该 key 对应的价格信息

#### Scenario: 空 OutcomePrices
- **WHEN** `OutcomePrices` 为空 map
- **THEN** `prices` 数组为空
