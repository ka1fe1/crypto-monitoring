## 1. 核心实现

- [x] 1.1 修改 `internal/tasks/polymarket_monitor_task.go` 的 `formatMarkets` 方法，提取 `OutcomePrices` 的 key 并按降序排序后遍历
- [x] 1.2 在文件头部添加 `sort` 包导入

## 2. 同步修改

- [x] 2.1 检查 `internal/tasks/general_monitor_task.go` 中相同逻辑，同步应用降序排序

## 3. 验证

- [x] 3.1 运行 `go build ./...` 确认编译通过
- [x] 3.2 运行相关单元测试确认无回归
