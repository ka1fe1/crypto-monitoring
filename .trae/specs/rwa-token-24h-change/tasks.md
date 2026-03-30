# Tasks

- [x] Task 1: 修改 TokenPriceMonitorTask 中 RWA 资产涨跌幅显示
  - [x] SubTask 1.1: 修改 `formatTokenPricesDetailed` 方法，RWA 资产使用 `PercentChange24h` 替代 `PercentChange1h`

- [x] Task 2: 修改 GeneralMonitorTask 中 RWA 资产涨跌幅显示
  - [x] SubTask 2.1: 修改 `formatTokenPricesSimple` 方法，RWA 资产使用 `PercentChange24h` 替代 `PercentChange1h`

- [x] Task 3: 更新单元测试
  - [x] SubTask 3.1: 更新测试用例验证 RWA 资产使用 24h 涨跌幅

# Task Dependencies
- Task 3 依赖 Task 1 和 Task 2
