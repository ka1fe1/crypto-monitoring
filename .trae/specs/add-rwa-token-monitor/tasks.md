# Tasks

- [x] Task 1: 更新配置结构体，新增 RWA token 相关字段
  - [x] SubTask 1.1: 在 `TokenPriceMonitorConfig` 中添加 `RwaTokenIds` 字段（字符串类型，YAML 解析）
  - [x] SubTask 1.2: 在 `TokenPriceMonitorConfig` 中添加 `RwaTokenIDs` 字段（切片类型，解析后存储）
  - [x] SubTask 1.3: 在 `TokenPriceMonitorConfig` 中添加 `RwaTokenNames` 字段（map 类型，存储 token ID 到中文名称的映射）
  - [x] SubTask 1.4: 在 `LoadConfig` 函数中添加解析 `RwaTokenIds` 和 `RwaTokenNames` 的逻辑

- [x] Task 2: 修改 TokenPriceMonitorTask 支持 RWA 资产
  - [x] SubTask 2.1: 在 `TokenPriceMonitorTask` 结构体中添加 `rwaTokenIds` 和 `rwaTokenNames` 字段
  - [x] SubTask 2.2: 修改 `NewTokenPriceMonitorTask` 函数，接收 RWA 相关参数
  - [x] SubTask 2.3: 修改 `run` 方法，同时获取 crypto 和 RWA 资产价格
  - [x] SubTask 2.4: 修改 `formatTokenPricesDetailed` 方法，区分显示 crypto 和 RWA 资产，RWA 显示中文名称

- [x] Task 3: 修改 GeneralMonitorTask 支持 RWA 资产
  - [x] SubTask 3.1: 在 `GeneralMonitorTask` 结构体中添加 `rwaTokenIds` 和 `rwaTokenNames` 字段
  - [x] SubTask 3.2: 修改 `NewGeneralMonitorTask` 函数，接收 RWA 相关参数
  - [x] SubTask 3.3: 修改 `getTokenPriceContent` 方法，同时获取 crypto 和 RWA 资产价格
  - [x] SubTask 3.4: 修改 `formatTokenPricesSimple` 方法，区分显示 crypto 和 RWA 资产

- [x] Task 4: 更新 init.go 中的任务初始化逻辑
  - [x] SubTask 4.1: 修改 `TokenPriceMonitorTask` 初始化，传入 RWA 配置参数
  - [x] SubTask 4.2: 修改 `GeneralMonitorTask` 初始化，传入 RWA 配置参数

- [x] Task 5: 更新配置模板文件
  - [x] SubTask 5.1: 在 `config.yaml.temp` 中添加 `rwa_token_ids` 和 `rwa_token_names` 示例配置

- [x] Task 6: 编写单元测试
  - [x] SubTask 6.1: 更新 `token_price_monitor_task_test.go` 测试 RWA 功能
  - [x] SubTask 6.2: 更新 `general_monitor_task_test.go` 测试 RWA 功能

# Task Dependencies
- Task 2 依赖 Task 1
- Task 3 依赖 Task 1
- Task 4 依赖 Task 2 和 Task 3
- Task 6 依赖 Task 1-5
