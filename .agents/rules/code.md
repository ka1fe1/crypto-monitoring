---
trigger: always_on
---

# 语言
所有的响应，请使用中文回答

# 工具类的写法要求

1. `/pkg/utils` 目录下写对应的工具类，其中需要新建文件夹，如需要新建一个 `polymarket` 的工具类，则需要在新建的 `/pkg/utils/polymarket` 文件夹下，新建 `ploymarket.go` 文件
2. 所有的工具类，都需要写对应的单元测试文件，与工具类放在同一文件夹下
3. 工具类若有需要外部传入的配置，如 `api_key` 等，则默认先在 `config/config.go` 中配置，其他调用工具类的方法，从配置中读取并传入
4. 工具类的方法如果需要定义其请求和响应的结构体，可以新建一个文件，如 `ploymarket_vo.go`，单独存储该工具类对应方法的请求与响应的结构体。这样工具类中就只有方法的定义了，不会因请求响应结构体太大，导致文件行数比较多。


# 单元测试的写法要求

1. 单元测试若需要从配置文件中读取一些配置，先在 TestMain 中读取配置文件，再在具体单元测试方法中引用即可，如:

```go
var (
	cfg *config.Config
	bot *DingBot
)

// loadTestConfig resolves the absolute path to config.yaml and loads it.
func loadTestConfig() (*config.Config, error) {
	// 1. Get the absolute path of the current file to determine the project root.
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("failed to get current file path")
	}

	// The current file is in <ProjectRoot>/pkg/utils/alter/dingding/bot_test.go
	// So we go up four levels to get to <ProjectRoot>
	rootDir := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(filename)))))

	// 2. Construct the absolute path to config.yaml
	configPath := filepath.Join(rootDir, "config", "config.yaml")

	// 3. Load the configuration
	return config.LoadConfig(configPath)
}

func TestMain(m *testing.M) {
	var err error
	cfg, err = loadTestConfig()
	if err != nil {
		log.Printf("Warning: Could not load config: %v", err)
	}

	var token, secret, keyword string
	if cfg != nil {
		token = cfg.DingTalk[constant.DEFAULT_BOT_NAME].AccessToken
		secret = cfg.DingTalk[constant.DEFAULT_BOT_NAME].Secret
		keyword = cfg.DingTalk[constant.DEFAULT_BOT_NAME].Keyword
	}

	bot = NewDingBot(token, secret, keyword)
	os.Exit(m.Run())
}
```

2. 单元测试中不应使用 `defaultMockTransport` 或其他形式的 HTTP Transport 模拟。应直接使用真实配置与外部服务交互。

# 监控 task 和 service 的写法要求

1. 文件名以 `xx_monitor_task` 或 `xx_service`，如 `polymarket_monitor_task`, `polymarket_service`
2. monitor task 和 service 所需要的配置，需从 `config/config.go` 中读取，即需要先配置相关的配置项
3. monitor task 和 service 也需要写对应的单元测试，单元测试遵循上面单元测试的写法要求