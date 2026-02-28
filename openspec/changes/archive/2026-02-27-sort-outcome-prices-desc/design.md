## Context

`polymarket_monitor_task.go` 的 `formatMarkets` 方法中，`market.OutcomePrices` 类型为 `map[string]float64`。Go 中 map 遍历顺序是随机的，导致每次生成的消息中价格展示顺序不确定。

## Goals / Non-Goals

**Goals:**
- 将 `OutcomePrices` 按 key（name）降序排列后输出到 `prices` 数组
- 保持输出格式不变，仅改变遍历顺序

**Non-Goals:**
- 不修改 `OutcomePrices` 的数据结构本身
- 不修改 `MarketDetail` 的定义

## Decisions

### 使用 `sort.Sort(sort.Reverse(sort.StringSlice(keys)))` 排序

提取 map 的所有 key 到 `[]string`，使用标准库 `sort` 包进行降序排序，然后按排序后的 key 顺序遍历 map 填充 `prices`。

**理由**: Go 标准库自带，无需外部依赖，实现简洁清晰。

**替代方案**: 使用 `slices.SortFunc`（Go 1.21+），但 `sort` 包兼容性更好且足够简洁。

## Risks / Trade-offs

- **性能**: map key 排序引入微小开销，但 OutcomePrices 通常只有 2-3 个 key，可忽略不计。
- **breaking change**: 无，仅改变展示顺序。
