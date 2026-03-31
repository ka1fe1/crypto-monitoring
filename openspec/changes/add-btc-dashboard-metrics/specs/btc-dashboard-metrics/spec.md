## ADDED Requirements

### Requirement: 获取及计算基础资产与链上衍生指标
系统 MUST 能够从可用的外部数据源（如 Binance, OKX, mempool.space, alternative.me）成功获取核心基础数据并计算出宏观判断所需的衍生指标（200WMA、ahr999、恐慌贪婪指数、减半倒计时信息）。

#### Scenario: 成功计算 ahr999 与 200WMA 
- **WHEN** 监控定时任务触发，且外部交易所 API 正常返回 K 线历史数据
- **THEN** 系统正确求得过去 200 周价格平均线和日线 ahr999 值，并评判出指标所在状态区间（例如 ahr999 < 1.2 为定投区间）。

#### Scenario: 成功获取外部实时状态面板数据
- **WHEN** 触发数据拉取时
- **THEN** 系统成功从 mempool.space 和 alternative.me 接口获取到最新出块高度和当前 FGI（Fear & Greed Index）值，并正确预估剩余减半天数。

### Requirement: 汇总指标并进行监控预警推送
系统 MUST 负责将每次计算得到的各个异动状态整合成完整的 Markdown 报告，定时向绑定的钉钉机器人（Dingding Bot）推送。

#### Scenario: 监控信息定时推送
- **WHEN** 定时任务到达且所有上述监控数据已组装完毕
- **THEN** 系统将拼接好的监控文本通过配置好的 DingTalk Secret/Token 发送至指定的运维或投资提示群组。
