## 1. 配置文件与预备工作

- [x] 1.1 在 `config/config.yaml` 和 `config.go` 中补充所需的外部 API 地址信息（虽然为硬编码亦可，但按惯例预留）。

## 2. 工具类 (Utils) 及单元测试开发

- [x] 2.1 创建 `pkg/utils/binance/binance_vo.go`，定义 K 线的响应结构体。
- [x] 2.2 创建 `pkg/utils/binance/binance.go`，实现获取指定时间周期的 K 线逻辑，并在同级目录添加 `binance_test.go`。
- [x] 2.3 创建 `pkg/utils/mempool/mempool.go` 及其独立子文件夹结构，完成获取最新爆块高度接口的方法与对应的单元测试。
- [x] 2.4 创建 `pkg/utils/alternative/alternative.go` 及相关结构体，实现抓取贪婪恐惧指数的逻辑与单测。

## 3. 核心计算服务 (Service)

- [x] 3.1 创建业务逻辑层 `btc_dashboard_service` 文件夹及同名文件。
- [x] 3.2 在服务中集成计算 200WMA 均线、ahr999 乘数等公式逻辑。
- [x] 3.3 在服务中聚合工具类返回的多项指标拼装成报告数据结构，并对计算与判断代码添加完整的单元测试 `btc_dashboard_service_test.go`。

## 4. 调度层 (Monitor Task) 和群组推送

- [x] 4.1 创建 `btc_dashboard_monitor_task.go` 以及定时任务启动入口逻辑。
- [x] 4.2 整合 `DingBot` 统一发送逻辑，将 Service 计算的报告结果 Markdown 格式化并投递至钉钉预警群组。
- [x] 4.3 提供 Monitor Task 相关的单元测试或入口模拟。
