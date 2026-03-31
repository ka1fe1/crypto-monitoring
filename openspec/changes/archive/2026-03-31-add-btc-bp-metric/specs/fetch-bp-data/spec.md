## ADDED Requirements

### Requirement: 获取比特币均衡价格

必须从相应的第三方开放接口（例如 Coinglass 或同等可依赖链上数据商）中获取实时的 Balanced Price (均衡价格)。

#### Scenario: 当服务抓取常规指标时连带获取均衡价格 BP
- **WHEN** 指标监控定时任务触发，调用 `FetchAndCalculateMetrics` 聚合函数时，同时请求 BP 数据对应的接口。
- **THEN** 若请求成功且返回合理参数，需要解析该金额（如 `$40466`）记录进当期 Metrics 结果中以备组装报告。

#### Scenario: 服务抓取失败或响应超时
- **WHEN** 调用外置 BP 接口期间因网络问题、速率限制或其他 HTTP Error 导致获取失败。
- **THEN** 计算器不应该报抛异常阻挠其余流程，仅仅会将 BP 的数值标明为 0/缺失，在报告侧对应空载输出或者标明为“数据无法获取”，确保整个系统（Ahr999 及 FGI 等）照常运行。
