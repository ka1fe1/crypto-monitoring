# Crypto Monitoring Project Analysis

**项目名称**: Crypto Monitoring Service
**技术栈**: Go 1.25+, Gin Framework

整个项目是一个综合性的加密资产与社交媒体监控机器人，旨在实时监测多种数据源，并在触发条件时通过**钉钉机器人**发送通知。

## 1. 核心架构与请求链路 (Architecture & Flow)

整个系统采用“双引擎”结构：REST API 服务 + 定时监控后台任务

```text
┌─────────────────────────┐               ┌─────────────────────────────────┐
│       HTTP API Server   │               │      Background Task Engine     │
│       (Gin Framework)   │               │       (Time-based Polling)      │
├─────────────────────────┤               ├─────────────────────────────────┤
│ • GET /api/v1/token/..  │    Reads      │ • TokenPriceMonitorTask         │
│ • GET /api/v1/dex/pair  │ ◄───────────► │ • DexPairAlterTask              │
│ • GET /ping             │               │ • PolymarketMonitorTask         │
└────────────┬────────────┘   Config      │ • TwitterMonitorTask           │
             │                Driven      │ • NFTFloorPriceMonitorTask      │
             │                            └───────────────┬─────────────────┘
             │                                            │ (Pulls Data)
             ▼                                            ▼
   ═════════════════════════════════════════════════════════════════════════
      external Integrations: CoinMarketCap, OpenSea, Twitter, Polymarket
   ═════════════════════════════════════════════════════════════════════════
                                                          │ (Pushes Alerts)
                                                          ▼
                                            ┌───────────────────────────────┐
                                            │      DingTalk Bot Clients     │
                                            │     (Supports Quiet Hours)    │
                                            └───────────────────────────────┘
```

## 2. 主要功能 (Main Features)

- **市场数据监控**：
  - DEX（去中心化交易所）上的代币价格与流动性变化
  - CEX 上的 Token 价格
  - 重点 NFT 项目（如 OpenSea平台）的地板价
- **特定生态监控**：
  - 集成 Polymarket 预测平台，监控指定事件的概率变化或交易活动
  - 集成 Coinglass 相关指数（如 AHR999、恐慌与贪婪指数）
- **社交媒体监控**：
  - 监控特定核心 Twitter 用户的最新推文
  - 支持特定关键词 (Keywords) 过滤推送
- **数据查询接口**：
  - 提供一组外部可访问的 HTTP API，用于查询 Token 和 DEX 交易对状态

## 3. 输入输出 (Inputs & Outputs)

### 输入 (Inputs)
1. **静态配置 (`config.yaml`)**：定义了所有监控规则、轮询频率、静默时间 (Quiet Hours)、各大平台的 API Key，以及钉钉群机器人的 Webhook 参数。
2. **三方数据源 API**：
   - *CoinMarketCap* / *DexScreener* (价格、流动性)
   - *OpenSea* (NFT 报价)
   - *Twitter Developer API* (社交推文)
   - *Polymarket API* (预测市场指标)

### 输出 (Outputs)
1. **钉钉群通知**：基于 `pkg/utils/alter/dingding` 模块，将经过清洗和格式化的数据推送到钉钉。支持配置免打扰模式。
2. **REST JSON 响应**：终端或前端应用调用由 `gin` 驱动的 API，返回标准化 JSON。

## 4. 现有代码实现方案 (Technical Implementation)

- **服务框架**：
  - Web 框架：使用 `gin-gonic/gin`
  - 路由管理：结构化的 `handlers` 与 `service` 分层
  - 自动生成了 Swagger API 文档
- **代码结构**：
  - 严格遵循 `cmd/` (应用入口), `internal/` (业务逻辑核心代码), `pkg/` (公共依赖与工具类) 的经典 Go 目录结构。
  - **依赖注入**：主函数 `main.go` 负责初始化各大 API 客户端（如 `CoinMarketClient`, `OpenSeaClient`, `TwitterClient`, `DingBot` 等），并显式注入到 `service` 和后台任务初始化流程中。
- **工程化实践**：
  - 提供了 `Makefile` 封装常用命令 (`make deps`, `make swagger`, `make run`)。
  - 提供了完整的 `Dockerfile` 方案以支持容器化部署。
