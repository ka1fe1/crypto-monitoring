# RWA Token 24h 涨跌幅显示 Spec

## Why
RWA（Real World Asset）代币化资产的价格波动相对较小，1 小时涨跌幅变化不明显，使用 24 小时涨跌幅更能反映资产的实际价格走势。

## What Changes
- 修改 `TokenPriceMonitorTask` 中 RWA 资产的涨跌幅显示，从 1 小时改为 24 小时
- 修改 `GeneralMonitorTask` 中 RWA 资产的涨跌幅显示，从 1 小时改为 24 小时

## Impact
- Affected specs: RWA Token Monitor
- Affected code:
  - `internal/tasks/token_price_monitor_task.go`
  - `internal/tasks/general_monitor_task.go`

## MODIFIED Requirements

### Requirement: RWA Token Price Display
RWA 资产价格显示应使用 24 小时涨跌幅。

#### Scenario: RWA 资产显示 24h 涨跌幅
- **WHEN** 系统显示 RWA 资产价格
- **THEN** 涨跌幅应显示 24 小时变化百分比 (`PercentChange24h`)

#### Scenario: Crypto 资产保持 1h 涨跌幅
- **WHEN** 系统显示 Crypto 资产价格
- **THEN** 涨跌幅应保持显示 1 小时变化百分比 (`PercentChange1h`)
