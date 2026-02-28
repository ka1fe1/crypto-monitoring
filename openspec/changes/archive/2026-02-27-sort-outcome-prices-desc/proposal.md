## Why

`formatMarkets` 方法中遍历 `map[string]float64` 类型的 `OutcomePrices` 时，Go map 的遍历顺序是随机的，导致每次发送的消息中 Outcome Prices 的显示顺序不一致，影响可读性和前后对比。需要按 name 降序排列，使输出稳定且可预测。

## What Changes

- 修改 `polymarket_monitor_task.go` 中 `formatMarkets` 方法，将 `OutcomePrices` map 的 key 提取并按降序排序后，再依次插入 `prices` 数组。
- 同步修改 `general_monitor_task.go` 中相同逻辑的代码（如适用）。

## Capabilities

### New Capabilities

- `sorted-outcome-display`: 对 OutcomePrices 按 name 降序排序后输出，保证消息展示顺序稳定。

### Modified Capabilities

_(无)_

## Impact

- 影响文件：`internal/tasks/polymarket_monitor_task.go`，可能还有 `internal/tasks/general_monitor_task.go`
- 无 API 变更、无依赖变更
- 需要引入 `sort` 标准库包
