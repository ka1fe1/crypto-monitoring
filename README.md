# Crypto Monitoring Service

这个项目是一个基于 Go (Gin) 的多功能加密资产与社交媒体监控服务，旨在通过钉钉机器人实时推送各类监控提醒。

## 核心功能

### 1. 监控任务 (Background Tasks)
支持通过 `config.yaml` 配置多种类型的监控任务，并支持 **Quiet Hours** (免打扰) 机制。

*   **DEX 价格预警** (`DexPairAlterTask`)
    *   监控指定链上 DEX 交易对的价格变化、流动性变动。
    *   数据源：CoinMarketCap / DexScreener。
*   **Token 价格监控** (`TokenPriceMonitorTask`)
    *   监控 CEX/DEX 代币价格。
*   **NFT 地板价监控** (`NFTFloorPriceMonitorTask`)
    *   监控 OpenSea 等平台的 NFT Collection 地板价。
*   **Polymarket 预测市场监控** (`PolymarketMonitorTask`)
    *   监控 Polymarket 特定市场的概率变化与交易活动。
*   **Twitter (X) 监控** (`TwitterMonitorTask`)
    *   监控指定 Twitter 用户的最新推文。
    *   支持解析 Snowflake ID 获取发推时间，提供更友好的日志与通知展示。
    *   **关键词过滤支持**: 支持针对特定用户配置关键词 (Keywords) 过滤，仅推送包含特定关键词的推文。
*   **Coinglass 数据集成** (`CoinGlass`)
    *   集成 Coinglass API，支持获取 AHR999 指数（囤币指标）及加密货币恐慌与贪婪指数 (Fear & Greed Index)。

### 2. API 服务 (RESTful API)
提供 HTTP 接口供外部系统集成，并集成了 **Swagger** 文档。
*   `GET /api/v1/token/price`: 查询 Token 实时价格。
*   `GET /api/v1/dex/pair`: 查询 DEX 交易对详情。
*   `GET /ping`: 健康检查。

### 3. 特性与组件
*   **多平台集成**: CoinMarketCap, OpenSea, Polymarket, Twitter API, DingTalk Bot。
*   **按需配置**: 支持针对每个 Task 独立配置轮询间隔、机器人 Token、监控目标。
*   **免打扰模式 (Quiet Hours)**: 支持配置特定时间段（如 00:00-08:00）暂停或降低推送频率。
*   **Docker 化**: 提供完整的 Docker 构建与部署支持。

## 快速开始

### 依赖
*   Go 1.25+
*   相关 API Key (CoinMarketCap, OpenSea, Twitter, etc.)

### 运行
```bash
# 安装依赖
make deps

# 生成文档
make swagger

# 运行服务
make run
```

### 访问
*   Swagger UI: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

### Docker 运行

#### 1. 构建镜像
```bash
docker build -t crypto-monitoring:25121702 .
```

#### 2. 运行容器 (挂载配置文件)
使用 `-v` 参数将本地的 `config.yaml` 挂载到容器内的 `/app/config/config.yaml`。

```bash
docker run -d \
  --name crypto-monitoring \
  -v $(pwd)/docker-app/crypto-monitor/config.yaml:/app/config/config.yaml \
  -p 8080:8080 \
  tataka1takes2/crypto-monitoring:2601161628
```

#### 3. 常用命令

*   **tag**
    ```bash
    docker tag crypto-monitoring:2601161604 tataka1takes2/crypto-monitoring:2601161604
    ```

*   **推送**
    ```bash
    docker push tataka1takes2/crypto-monitoring:2601161604
    ```

*   **拉取**
    ```bash
    docker pull tataka1takes2/crypto-monitoring:2601161628
    ```

*   **删除容器**
    ```bash
    docker container rm -f crypto-monitoring
    ```
*   **删除镜像**
    ```bash
    docker image rm crypto-monitoring
    ```
*   **重启容器**
    ```bash
    docker restart crypto-monitoring
    ```

*   **停止容器**
    ```bash
    docker stop crypto-monitoring
    ```
*   **查看日志**
    ```bash
    docker logs -f crypto-monitoring
    ```



