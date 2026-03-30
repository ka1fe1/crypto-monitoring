# RWA Token Monitor Spec

## Why
当前项目中的 token 价格监控仅支持加密货币资产（crypto assets），但用户需要同时监控 RWA（Real World Asset）代币化资产（如美股代币化）。由于 RWA 资产的 symbol 不太熟悉，需要在通知中显示中文名称以便识别。

## What Changes
- 在 `TokenPriceMonitorConfig` 中新增 `rwa_token_ids` 配置字段
- 新增 `rwa_token_names` 配置字段，用于存储 RWA token ID 到中文名称的映射
- 修改 `TokenPriceMonitorTask`，支持同时监控 crypto assets 和 RWA assets
- 修改钉钉通知格式，区分两类资产显示，RWA 资产显示中文名称

## Impact
- Affected specs: Token Price Monitor
- Affected code:
  - `config/config.go`
  - `internal/tasks/token_price_monitor_task.go`
  - `internal/tasks/general_monitor_task.go`
  - `config/config.yaml.temp`

## ADDED Requirements

### Requirement: RWA Token Configuration
系统应支持配置 RWA（Real World Asset）代币化资产的监控。

#### Scenario: 配置 RWA token IDs
- **WHEN** 用户在配置文件中设置 `rwa_token_ids` 字段
- **THEN** 系统应解析并存储 RWA token IDs 列表

#### Scenario: 配置 RWA token 中文名称映射
- **WHEN** 用户在配置文件中设置 `rwa_token_names` 字段
- **THEN** 系统应建立 token ID 到中文名称的映射关系

### Requirement: RWA Token Price Monitoring
系统应复用现有的 CMC API 方法和 token monitor task 来监控 RWA 资产价格。

#### Scenario: 获取 RWA 资产价格
- **WHEN** 系统执行 token 价格监控任务
- **THEN** 系统应同时获取 crypto assets 和 RWA assets 的价格

### Requirement: DingTalk Notification with Asset Type Distinction
系统应在钉钉通知中区分显示 crypto assets 和 RWA assets。

#### Scenario: 显示分类资产通知
- **WHEN** 系统发送钉钉价格通知
- **THEN** 通知应分为两个部分：
  - "### Crypto Assets" 部分：显示加密货币资产
  - "### RWA Assets" 部分：显示 RWA 资产，symbol 后附带中文名称

#### Scenario: RWA 资产显示格式
- **WHEN** 显示 RWA 资产价格
- **THEN** 格式应为：`- **SYMBOL (中文名称)**: ***$价格*** (涨跌幅%)`

### Requirement: General Monitor Task Support
GeneralMonitorTask 应支持同时显示 crypto assets 和 RWA assets。

#### Scenario: General Monitor 显示 RWA
- **WHEN** GeneralMonitorTask 执行 token_price 模块
- **THEN** 应同时显示 crypto assets 和 RWA assets 的价格信息
