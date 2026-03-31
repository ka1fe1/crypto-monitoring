## 1. 接入外部 API 配置

- [x] 1.1 在 `config/config.go` 增加用于请求 BP 的链上 API 的选项（支持 Coinglass 等供应商），或者为其设定默认 URL 逻辑。
- [x] 1.2 在相关测试文件和程序入口 (`init.go`) 中处理对新增 Provider 配置与服务注入。

## 2. BP Provider 接口及客户端实现

- [x] 2.1 增加 `BalancePriceProvider` 的规范接口（如提取至 `internal/service/btc_dashboard_service.go` 的内部接口列表），声明 `GetBalancedPrice() (float64, error)`。
- [x] 2.2 在 `pkg/utils/` 下创建一个实际调用此外部 HTTP API 以提取 BP 数据的 Client SDK 逻辑（例如 `pkg/utils/coinglass/`）。
- [x] 2.3 在同样的包目录下补全针对此实现方法的基础单元测试。

## 3. 业务层计算结合

- [x] 3.1 修改 `BtcDashboardService.FetchAndCalculateMetrics` 以在 goroutine 协程组（或其他异步或顺延流程）中同时拉取最新现价与 BP 估值数据。
- [x] 3.2 实现新的结构体字段（如 `BalancedPrice`），利用现价 / BP 产出倍率，添加类似判断 `1 异常低估 / 1-1.x 正常 / 高出 2 等` 的估值域阈值判断。
- [x] 3.3 修改 `GenerateMarkdownReport` 生成的内容模板，注入 `"Balanced Price (均衡价格)"` 及 `当前价格 / BP = xxx` 带 Emoji 指示的状态行。

## 4. 健壮性与单测更新

- [x] 4.1 在 `internal/service/btc_dashboard_service_test.go` 中新增对应 `MockBalancePriceProvider` 并适配进各种服务测试实例内。
- [x] 4.2 完善并运行 `go test ./...` 保证新加的容错模式——即当 BP 服务故障超时不阻碍原有 `Ahr999` 推送的流程有效运转无异常。
