# Crypto Monitoring Service

这个项目是一个基于 Go (Gin) 的加密货币监控服务，主要包含以下功能点：

## 1. API 服务 (RESTful API)
提供了一系列 HTTP 接口供外部调用，并集成了 **Swagger** 文档。

*   **Token 价格查询**: `GET /api/v1/token/price`
    *   **功能**：根据 ID 查询加密货币的实时价格和 Symbol。
    *   **实现**：`TokenService` -> `CoinMarketCap Client`。
*   **DEX 交易对信息**: `GET /api/v1/dex/pair`
    *   **功能**：根据合约地址和网络查询 DEX 交易对的详细信息（价格、流动性、涨跌幅等）。
    *   **实现**：`DexPairService` -> `CoinMarketCap Client`。
*   **健康检查**: `GET /ping`
    *   **功能**：检查服务是否存活。

## 2. 后台任务 (Background Tasks)
*   **价格预警任务**: `DexPairAlterTask`
    *   **功能**：每分钟定时轮询指定 DEX 交易对的价格。
    *   **通知**：通过 **钉钉机器人 (DingTalk Bot)** 发送 Markdown 格式的价格预警消息，包含价格、涨跌幅、流动性等信息。

## 3. 核心组件与集成
*   **CoinMarketCap Client**: 封装了 CoinMarketCap 的 API 调用（支持 v2 和 v4 接口），用于获取核心数据。
*   **DingTalk Bot**: 封装了钉钉机器人的消息发送功能。

## 4. 基础设施
*   **Swagger 文档**: 自动生成的 API 文档，可在线调试 (`/swagger/index.html`)。
*   **配置管理**: 使用 `config.yaml` 管理 API Key、预警配置等。
*   **Makefile**: 提供了 `build`, `run`, `swagger` 等常用命令，方便开发和部署。

## 快速开始

### 依赖
*   Go 1.25+
*   CoinMarketCap API Key

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
docker build -t crypto-monitoring .
```

#### 2. 运行容器 (挂载配置文件)
使用 `-v` 参数将本地的 `config.yaml` 挂载到容器内的 `/app/config/config.yaml`。

```bash
docker run -d \
  --name crypto-monitoring \
  -v $(pwd)/docker-app/crypto-monitor/config.yaml:/app/config/config.yaml \
  -p 8080:8080 \
  crypto-monitoring
```

#### 3. 常用命令

*   **tag**
    ```bash
    docker tag crypto-monitoring:25121001 tataka1takes2/crypto-monitoring:25121001
    ```

*   **推送**
    ```bash
    docker push tataka1takes2/crypto-monitoring:25121001
    ```

*   **拉取**
    ```bash
    docker pull tataka1takes2/crypto-monitoring:25121001
    ```

*   **查看日志**
    ```bash
    docker logs -f crypto-monitoring
    ```

*   **重启容器**
    ```bash
    docker restart crypto-monitoring
    ```

*   **停止容器**
    ```bash
    docker stop crypto-monitoring
    ```