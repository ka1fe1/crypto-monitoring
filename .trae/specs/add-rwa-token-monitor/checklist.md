# Checklist

## 配置相关
- [x] TokenPriceMonitorConfig 包含 RwaTokenIds 字段（YAML 字符串）
- [x] TokenPriceMonitorConfig 包含 RwaTokenIDs 字段（解析后的切片）
- [x] TokenPriceMonitorConfig 包含 RwaTokenNames 字段（map[string]string）
- [x] LoadConfig 函数正确解析 rwa_token_ids 配置
- [x] LoadConfig 函数正确解析 rwa_token_names 配置

## TokenPriceMonitorTask 相关
- [x] TokenPriceMonitorTask 结构体包含 rwaTokenIds 字段
- [x] TokenPriceMonitorTask 结构体包含 rwaTokenNames 字段
- [x] NewTokenPriceMonitorTask 接收 RWA 相关参数
- [x] run 方法同时获取 crypto 和 RWA 资产价格
- [x] 钉钉通知区分显示 "### Crypto Assets" 和 "### RWA Assets"
- [x] RWA 资产显示格式为 "SYMBOL (中文名称)"

## GeneralMonitorTask 相关
- [x] GeneralMonitorTask 结构体包含 rwaTokenIds 字段
- [x] GeneralMonitorTask 结构体包含 rwaTokenNames 字段
- [x] NewGeneralMonitorTask 接收 RWA 相关参数
- [x] getTokenPriceContent 同时获取 crypto 和 RWA 资产价格
- [x] formatTokenPricesSimple 区分显示 crypto 和 RWA 资产

## 初始化相关
- [x] init.go 中 TokenPriceMonitorTask 初始化传入 RWA 配置
- [x] init.go 中 GeneralMonitorTask 初始化传入 RWA 配置

## 配置模板
- [x] config.yaml.temp 包含 rwa_token_ids 示例配置
- [x] config.yaml.temp 包含 rwa_token_names 示例配置

## 测试
- [x] token_price_monitor_task_test.go 测试通过
- [x] general_monitor_task_test.go 测试通过
