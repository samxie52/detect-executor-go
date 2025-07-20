# Detect-Executor-Go 项目实现文档

## 项目概述
本文档记录了 Detect-Executor-Go 项目的完整实现过程，按照最佳编程建议的开发顺序逐步实现。

## 实现进度跟踪

### 开发阶段规划

#### 第一阶段：基础框架搭建 ✅
- [x] 1. 初始化 Go 模块 (`go.mod`) ✅
- [x] 2. 创建配置管理 (`pkg/config/`) ✅
- [x] 3. 实现日志系统 (`pkg/logger/`) ✅
- [x] 4. 搭建HTTP服务器 (`cmd/server/main.go`) ✅

#### 第二阶段：数据层实现 ✅
- [x] 5. 数据库连接 (`pkg/database/`) ✅
- [x] 6. 数据模型定义 (`internal/model/`) ✅
- [x] 7. 数据访问层 (`internal/repository/`) ✅

#### 第三阶段：业务逻辑实现 ✅
- [x] 8. 外部客户端 (`pkg/client/`) ✅
- [x] 9. 核心业务服务 (`internal/service/`) ✅
- [x] 10. HTTP处理器 (`internal/handler/`) ✅

#### 第四阶段：完善功能 🚧
- [ ] 11. 监控和指标 (`pkg/metrics/`)
- [x] 12. 测试用例 (已有基础测试)

---

## 步骤 1: 初始化 Go 模块

### 目标
创建 Go 项目的基础结构，初始化模块管理。

### 操作步骤

#### 1.1 创建项目目录
```bash
# 创建项目根目录
mkdir detect-executor-go
cd detect-executor-go
```

#### 1.2 初始化 Go 模块
```bash
# 初始化 Go 模块
go mod init detect-executor-go
```

#### 1.3 创建基础目录结构
```bash
# 创建主要目录结构
mkdir -p cmd/server
mkdir -p internal/{handler,service,repository,model}
mkdir -p pkg/{config,logger,database,client}
mkdir -p configs
mkdir -p deployments
mkdir -p scripts
mkdir -p docs
```

**实际项目结构**：
```
detect-executor-go/
├── cmd/
│   └── server/              # 服务器启动入口
├── internal/                # 内部包（不对外暴露）
│   ├── handler/            # HTTP处理器层
│   ├── service/            # 业务逻辑层
│   ├── repository/         # 数据访问层
│   └── model/              # 数据模型
├── pkg/                    # 可复用的包
│   ├── config/             # 配置管理
│   ├── logger/             # 日志系统
│   ├── database/           # 数据库连接
│   └── client/             # 外部客户端
├── configs/                # 配置文件
├── docs/                   # 文档
└── scripts/                # 脚本文件
```

#### 1.4 创建基础文件
```bash
# 创建基础配置文件
touch configs/config.yaml
touch configs/config.dev.yaml
touch configs/config.prod.yaml

# 创建构建相关文件
touch Makefile
touch .gitignore
touch README.md
```

#### 1.5 配置 .gitignore
```gitignore
# 二进制文件
/detect-executor-go
*.exe
*.exe~
*.dll
*.so
*.dylib

# 测试二进制文件
*.test

# 输出目录
/bin/
/dist/

# 依赖目录
/vendor/

# Go workspace file
go.work

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# 操作系统
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# 日志文件
*.log

# 配置文件（敏感信息）
configs/config.local.yaml

# 临时文件
/tmp/
```

#### 1.6 创建基础 Makefile
```makefile
# Makefile for detect-executor-go

# 变量定义
APP_NAME := detect-executor-go
VERSION := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date +%Y-%m-%d\ %H:%M:%S)
GO_VERSION := $(shell go version | awk '{print $$3}')

# 构建标志
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GoVersion=$(GO_VERSION)"

# 默认目标
.PHONY: all
all: build

# 构建
.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	go build $(LDFLAGS) -o bin/$(APP_NAME) cmd/server/main.go

# 运行
.PHONY: run
run:
	@echo "Running $(APP_NAME)..."
	go run cmd/server/main.go

# 测试
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

# 测试覆盖率
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# 代码检查
.PHONY: lint
lint:
	@echo "Running linter..."
	golangci-lint run

# 格式化代码
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# 整理依赖
.PHONY: tidy
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

# 清理
.PHONY: clean
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf dist/
	rm -f coverage.out coverage.html

# Docker 构建
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(APP_NAME):$(VERSION) -f deployments/docker/Dockerfile .

# 帮助
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  lint          - Run linter"
	@echo "  fmt           - Format code"
	@echo "  tidy          - Tidy dependencies"
	@echo "  clean         - Clean build artifacts"
	@echo "  docker-build  - Build Docker image"
	@echo "  help          - Show this help"
```

#### 1.7 创建基础 README.md
```markdown
# Detect-Executor-Go

高性能视频设备检测执行模块

## 快速开始

### 环境要求
- Go 1.19+
- MySQL 8.0+
- Redis 6.0+

### 安装依赖
```bash
go mod download
```

### 运行应用
```bash
make run
```

### 构建应用
```bash
make build
```

## 项目结构

```
detect-executor-go/
├── cmd/                    # 应用程序入口
├── internal/               # 内部包
├── pkg/                    # 公共包
├── configs/                # 配置文件
├── deployments/            # 部署配置
└── tests/                  # 测试文件
```

## 开发指南

### 代码规范
- 使用 `golangci-lint` 进行代码检查
- 使用 `go fmt` 格式化代码
- 保持测试覆盖率在 80% 以上

### 构建命令
```bash
make build      # 构建应用
make test       # 运行测试
make lint       # 代码检查
make clean      # 清理构建产物
```
```

### 验证步骤

#### 1.8 验证项目结构
```bash
# 检查目录结构
tree -L 3

# 验证 go.mod 文件
cat go.mod

# 验证 Makefile
make help
```

### 预期结果
- 项目目录结构创建完成
- `go.mod` 文件生成，模块名为 `detect-executor-go`
- 基础配置文件和构建文件创建完成
- Makefile 可以正常执行 help 命令

### 状态
- [x] 步骤 1 完成：项目基础结构初始化

---

---

## 步骤 2: 创建配置管理

### 目标
实现灵活的配置管理系统，支持多环境配置、环境变量覆盖和配置热更新。

### 操作步骤

#### 2.1 添加配置管理依赖
```bash
# 添加 viper 配置管理库
go get github.com/spf13/viper

# 添加 fsnotify 文件监控库（用于配置热更新）
go get github.com/fsnotify/fsnotify
```

#### 2.2 创建配置结构体 (`pkg/config/config.go`)
```go
package config

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Config 应用程序配置结构
type Config struct {
	Server   ServerConfig   `yaml:"server" mapstructure:"server"`
	Database DatabaseConfig `yaml:"database" mapstructure:"database"`
	Redis    RedisConfig    `yaml:"redis" mapstructure:"redis"`
	Detect   DetectConfig   `yaml:"detect" mapstructure:"detect"`
	Log      LogConfig      `yaml:"log" mapstructure:"log"`
	Metrics  MetricsConfig  `yaml:"metrics" mapstructure:"metrics"`
	Client   ClientConfig   `yaml:"client" mapstructure:"client"`
}

// ClientConfig 客户端配置
type ClientConfig struct {
	HTTP         HTTPClientConfig   `yaml:"http" mapstructure:"http"`
	DetectEngine DetectEngineConfig `yaml:"detect_engine" mapstructure:"detect_engine"`
	Storage      StorageConfig      `yaml:"storage" mapstructure:"storage"`
	MessageQueue MessageQueueConfig `yaml:"message_queue" mapstructure:"message_queue"`
}

// HTTPClientConfig HTTP 客户端配置
type HTTPClientConfig struct {
	Timeout         time.Duration `mapstructure:"timeout"`
	RetryCount      int           `mapstructure:"retry_count"`
	RetryInterval   time.Duration `mapstructure:"retry_interval"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxConnsPerHost int           `mapstructure:"max_conns_per_host"`
}

// ServerConfig HTTP服务器配置
type ServerConfig struct {
	Host         string        `yaml:"host" mapstructure:"host"`
	Port         int           `yaml:"port" mapstructure:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout" mapstructure:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout" mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" mapstructure:"idle_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string        `yaml:"host" mapstructure:"host"`
	Port            int           `yaml:"port" mapstructure:"port"`
	Username        string        `yaml:"username" mapstructure:"username"`
	Password        string        `yaml:"password" mapstructure:"password"`
	Database        string        `yaml:"database" mapstructure:"database"`
	Charset         string        `yaml:"charset" mapstructure:"charset"`
	MaxIdleConns    int           `yaml:"max_idle_conns" mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `yaml:"max_open_conns" mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" mapstructure:"conn_max_lifetime"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host         string        `yaml:"host" mapstructure:"host"`
	Port         int           `yaml:"port" mapstructure:"port"`
	Password     string        `yaml:"password" mapstructure:"password"`
	DB           int           `yaml:"db" mapstructure:"db"`
	PoolSize     int           `yaml:"pool_size" mapstructure:"pool_size"`
	MinIdleConns int           `yaml:"min_idle_conns" mapstructure:"min_idle_conns"`
	DialTimeout  time.Duration `yaml:"dial_timeout" mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout" mapstructure:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout" mapstructure:"write_timeout"`
}

// DetectConfig 检测相关配置
type DetectConfig struct {
	MaxWorkers      int           `yaml:"max_workers" mapstructure:"max_workers"`
	BatchSize       int           `yaml:"batch_size" mapstructure:"batch_size"`
	QueueSize       int           `yaml:"queue_size" mapstructure:"queue_size"`
	Timeout         time.Duration `yaml:"timeout" mapstructure:"timeout"`
	RetryCount      int           `yaml:"retry_count" mapstructure:"retry_count"`
	RetryDelay      time.Duration `yaml:"retry_delay" mapstructure:"retry_delay"`
	CppEngineURL    string        `yaml:"cpp_engine_url" mapstructure:"cpp_engine_url"`
	VideoStreamURL  string        `yaml:"video_stream_url" mapstructure:"video_stream_url"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `yaml:"level" mapstructure:"level"`
	Format     string `yaml:"format" mapstructure:"format"`
	Output     string `yaml:"output" mapstructure:"output"`
	Filename   string `yaml:"filename" mapstructure:"filename"`
	MaxSize    int    `yaml:"max_size" mapstructure:"max_size"`
	MaxAge     int    `yaml:"max_age" mapstructure:"max_age"`
	MaxBackups int    `yaml:"max_backups" mapstructure:"max_backups"`
	Compress   bool   `yaml:"compress" mapstructure:"compress"`
}

// MetricsConfig 监控配置
type MetricsConfig struct {
	Enabled bool   `yaml:"enabled" mapstructure:"enabled"`
	Path    string `yaml:"path" mapstructure:"path"`
	Port    int    `yaml:"port" mapstructure:"port"`
}

// GetAddr 获取服务器地址
func (s *ServerConfig) GetAddr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// GetDSN 获取数据库连接字符串
func (d *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		d.Username, d.Password, d.Host, d.Port, d.Database, d.Charset)
}

// GetAddr 获取Redis地址
func (r *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

// GetAddr 获取监控服务地址
func (m *MetricsConfig) GetAddr() string {
	return fmt.Sprintf(":%d", m.Port)
}
```

#### 2.3 创建配置加载器 (`pkg/config/loader.go`)
```go
package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Loader 配置加载器
type Loader struct {
	v      *viper.Viper
	config *Config
}

// NewLoader 创建配置加载器
func NewLoader() *Loader {
	v := viper.New()
	
	// 设置配置文件类型
	v.SetConfigType("yaml")
	
	// 设置配置文件搜索路径
	v.AddConfigPath(".")
	v.AddConfigPath("./configs")
	v.AddConfigPath("/etc/detect-executor-go")
	v.AddConfigPath("$HOME/.detect-executor-go")
	
	// 设置环境变量前缀
	v.SetEnvPrefix("DETECT")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	
	return &Loader{
		v:      v,
		config: &Config{},
	}
}

// Load 加载配置
func (l *Loader) Load(configFile string) (*Config, error) {
	// 设置默认值
	l.setDefaults()
	
	// 如果指定了配置文件，则使用指定的文件
	if configFile != "" {
		l.v.SetConfigFile(configFile)
	} else {
		// 根据环境变量确定配置文件名
		env := os.Getenv("GO_ENV")
		if env == "" {
			env = "dev"
		}
		l.v.SetConfigName(fmt.Sprintf("config.%s", env))
	}
	
	// 读取配置文件
	if err := l.v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到，使用默认配置
			fmt.Printf("Warning: Config file not found, using defaults\n")
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}
	
	// 解析配置到结构体
	if err := l.v.Unmarshal(l.config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	return l.config, nil
}

// Watch 监控配置文件变化
func (l *Loader) Watch(callback func(*Config)) error {
	l.v.WatchConfig()
	l.v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Config file changed: %s\n", e.Name)
		
		// 重新解析配置
		if err := l.v.Unmarshal(l.config); err != nil {
			fmt.Printf("Failed to reload config: %v\n", err)
			return
		}
		
		// 调用回调函数
		if callback != nil {
			callback(l.config)
		}
	})
	
	return nil
}

// setDefaults 设置默认配置值
func (l *Loader) setDefaults() {
	// 服务器默认配置
	l.v.SetDefault("server.host", "0.0.0.0")
	l.v.SetDefault("server.port", 8080)
	l.v.SetDefault("server.read_timeout", "30s")
	l.v.SetDefault("server.write_timeout", "30s")
	l.v.SetDefault("server.idle_timeout", "60s")
	
	// 数据库默认配置
	l.v.SetDefault("database.host", "localhost")
	l.v.SetDefault("database.port", 3306)
	l.v.SetDefault("database.username", "root")
	l.v.SetDefault("database.password", "")
	l.v.SetDefault("database.database", "detect_db")
	l.v.SetDefault("database.charset", "utf8mb4")
	l.v.SetDefault("database.max_idle_conns", 10)
	l.v.SetDefault("database.max_open_conns", 100)
	l.v.SetDefault("database.conn_max_lifetime", "1h")
	
	// Redis默认配置
	l.v.SetDefault("redis.host", "localhost")
	l.v.SetDefault("redis.port", 6379)
	l.v.SetDefault("redis.password", "")
	l.v.SetDefault("redis.db", 0)
	l.v.SetDefault("redis.pool_size", 10)
	l.v.SetDefault("redis.min_idle_conns", 5)
	l.v.SetDefault("redis.dial_timeout", "5s")
	l.v.SetDefault("redis.read_timeout", "3s")
	l.v.SetDefault("redis.write_timeout", "3s")
	
	// 检测默认配置
	l.v.SetDefault("detect.max_workers", 1000)
	l.v.SetDefault("detect.batch_size", 100)
	l.v.SetDefault("detect.queue_size", 10000)
	l.v.SetDefault("detect.timeout", "30s")
	l.v.SetDefault("detect.retry_count", 3)
	l.v.SetDefault("detect.retry_delay", "1s")
	l.v.SetDefault("detect.cpp_engine_url", "http://localhost:8081")
	l.v.SetDefault("detect.video_stream_url", "http://localhost:8082")
	
	// 日志默认配置
	l.v.SetDefault("log.level", "info")
	l.v.SetDefault("log.format", "json")
	l.v.SetDefault("log.output", "stdout")
	l.v.SetDefault("log.filename", "logs/app.log")
	l.v.SetDefault("log.max_size", 100)
	l.v.SetDefault("log.max_age", 7)
	l.v.SetDefault("log.max_backups", 10)
	l.v.SetDefault("log.compress", true)
	
	// 监控默认配置
	l.v.SetDefault("metrics.enabled", true)
	l.v.SetDefault("metrics.path", "/metrics")
	l.v.SetDefault("metrics.port", 9090)
}
```

#### 2.4 创建配置文件模板

**开发环境配置** (`configs/config.dev.yaml`)
```yaml
# 开发环境配置
server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: "30s"
  write_timeout: "30s"
  idle_timeout: "60s"

database:
  host: "localhost"
  port: 3306
  username: "root"
  password: "password"
  database: "detect_db_dev"
  charset: "utf8mb4"
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: "1h"

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
  pool_size: 10
  min_idle_conns: 5
  dial_timeout: "5s"
  read_timeout: "3s"
  write_timeout: "3s"

detect:
  max_workers: 100
  batch_size: 10
  queue_size: 1000
  timeout: "30s"
  retry_count: 3
  retry_delay: "1s"
  cpp_engine_url: "http://localhost:8081"
  video_stream_url: "http://localhost:8082"

log:
  level: "debug"
  format: "text"
  output: "stdout"
  filename: "logs/app.log"
  max_size: 100
  max_age: 7
  max_backups: 10
  compress: true

metrics:
  enabled: true
  path: "/metrics"
  port: 9090
```

**生产环境配置** (`configs/config.prod.yaml`)
```yaml
# 生产环境配置
server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: "30s"
  write_timeout: "30s"
  idle_timeout: "120s"

database:
  host: "${DETECT_DATABASE_HOST}"
  port: 3306
  username: "${DETECT_DATABASE_USERNAME}"
  password: "${DETECT_DATABASE_PASSWORD}"
  database: "${DETECT_DATABASE_DATABASE}"
  charset: "utf8mb4"
  max_idle_conns: 20
  max_open_conns: 200
  conn_max_lifetime: "1h"

redis:
  host: "${DETECT_REDIS_HOST}"
  port: 6379
  password: "${DETECT_REDIS_PASSWORD}"
  db: 0
  pool_size: 50
  min_idle_conns: 10
  dial_timeout: "5s"
  read_timeout: "3s"
  write_timeout: "3s"

detect:
  max_workers: 1000
  batch_size: 100
  queue_size: 10000
  timeout: "60s"
  retry_count: 3
  retry_delay: "2s"
  cpp_engine_url: "${DETECT_CPP_ENGINE_URL}"
  video_stream_url: "${DETECT_VIDEO_STREAM_URL}"

log:
  level: "info"
  format: "json"
  output: "file"
  filename: "/var/log/detect-executor-go/app.log"
  max_size: 500
  max_age: 30
  max_backups: 30
  compress: true

metrics:
  enabled: true
  path: "/metrics"
  port: 9090
```

#### 2.5 创建配置测试文件 (`pkg/config/config_test.go`)
```go
package config

import (
	"os"
	"testing"
	"time"
)

func TestConfigLoad(t *testing.T) {
	// 创建临时配置文件
	configContent := `
server:
  host: "127.0.0.1"
  port: 9999
  read_timeout: "15s"

database:
  host: "testdb"
  port: 3307
  username: "testuser"
  password: "testpass"
  database: "testdb"

detect:
  max_workers: 500
  batch_size: 50
`
	
	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "config_test_*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	
	// 写入配置内容
	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}
	tmpFile.Close()
	
	// 加载配置
	loader := NewLoader()
	config, err := loader.Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	// 验证配置
	if config.Server.Host != "127.0.0.1" {
		t.Errorf("Expected server host '127.0.0.1', got '%s'", config.Server.Host)
	}
	
	if config.Server.Port != 9999 {
		t.Errorf("Expected server port 9999, got %d", config.Server.Port)
	}
	
	if config.Server.ReadTimeout != 15*time.Second {
		t.Errorf("Expected read timeout 15s, got %v", config.Server.ReadTimeout)
	}
	
	if config.Database.Host != "testdb" {
		t.Errorf("Expected database host 'testdb', got '%s'", config.Database.Host)
	}
	
	if config.Detect.MaxWorkers != 500 {
		t.Errorf("Expected max workers 500, got %d", config.Detect.MaxWorkers)
	}
}

func TestConfigDefaults(t *testing.T) {
	// 加载不存在的配置文件，应该使用默认值
	loader := NewLoader()
	config, err := loader.Load("nonexistent.yaml")
	if err != nil {
		t.Fatalf("Failed to load config with defaults: %v", err)
	}
	
	// 验证默认值
	if config.Server.Port != 8080 {
		t.Errorf("Expected default server port 8080, got %d", config.Server.Port)
	}
	
	if config.Database.Port != 3306 {
		t.Errorf("Expected default database port 3306, got %d", config.Database.Port)
	}
	
	if config.Detect.MaxWorkers != 1000 {
		t.Errorf("Expected default max workers 1000, got %d", config.Detect.MaxWorkers)
	}
}

func TestConfigMethods(t *testing.T) {
	config := &Config{
		Server: ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
		Database: DatabaseConfig{
			Host:     "dbhost",
			Port:     3306,
			Username: "user",
			Password: "pass",
			Database: "testdb",
			Charset:  "utf8mb4",
		},
		Redis: RedisConfig{
			Host: "redishost",
			Port: 6379,
		},
		Metrics: MetricsConfig{
			Port: 9090,
		},
	}
	
	// 测试地址生成方法
	if addr := config.Server.GetAddr(); addr != "localhost:8080" {
		t.Errorf("Expected server addr 'localhost:8080', got '%s'", addr)
	}
	
	if dsn := config.Database.GetDSN(); dsn != "user:pass@tcp(dbhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local" {
		t.Errorf("Unexpected DSN: %s", dsn)
	}
	
	if addr := config.Redis.GetAddr(); addr != "redishost:6379" {
		t.Errorf("Expected redis addr 'redishost:6379', got '%s'", addr)
	}
	
	if addr := config.Metrics.GetAddr(); addr != ":9090" {
		t.Errorf("Expected metrics addr ':9090', got '%s'", addr)
	}
}
```

### 验证步骤

#### 2.6 验证配置管理功能
```bash
# 运行配置测试
go test ./pkg/config/

# 检查依赖是否正确添加
go mod tidy
cat go.mod

# 验证配置文件格式
cat configs/config.dev.yaml
cat configs/config.prod.yaml
```

### 预期结果
- 配置管理包创建完成，包含完整的配置结构体定义
- 支持多环境配置文件加载
- 支持环境变量覆盖配置
- 支持配置文件热更新
- 包含完整的单元测试
- 依赖库正确添加到 go.mod

### 状态
- [ ] 步骤 2 完成：配置管理系统实现

---

---

## 步骤 3: 实现日志系统

### 目标
实现功能完善的日志系统，支持结构化日志、多种输出格式、日志轮转和中间件集成。

### 操作步骤

#### 3.1 添加日志相关依赖
```bash
# 添加 logrus 结构化日志库
go get github.com/sirupsen/logrus

# 添加日志轮转库
go get gopkg.in/natefinch/lumberjack.v2

# 添加 Gin 日志中间件
go get github.com/gin-gonic/gin
```

#### 3.2 创建日志管理器 (`pkg/logger/logger.go`)
```go
package logger

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger 日志管理器
type Logger struct {
	*logrus.Logger
	config *Config
}

// Config 日志配置
type Config struct {
	Level      string `yaml:"level"`
	Format     string `yaml:"format"`
	Output     string `yaml:"output"`
	Filename   string `yaml:"filename"`
	MaxSize    int    `yaml:"max_size"`
	MaxAge     int    `yaml:"max_age"`
	MaxBackups int    `yaml:"max_backups"`
	Compress   bool   `yaml:"compress"`
}

// NewLogger 创建日志管理器
func NewLogger(config *Config) (*Logger, error) {
	logger := logrus.New()
	
	// 设置日志级别
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)
	
	// 设置日志格式
	switch config.Format {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	default:
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}
	
	// 设置输出目标
	if err := setOutput(logger, config); err != nil {
		return nil, err
	}
	
	return &Logger{
		Logger: logger,
		config: config,
	}, nil
}

// setOutput 设置日志输出
func setOutput(logger *logrus.Logger, config *Config) error {
	switch config.Output {
	case "stdout":
		logger.SetOutput(os.Stdout)
	case "stderr":
		logger.SetOutput(os.Stderr)
	case "file":
		if config.Filename == "" {
			config.Filename = "logs/app.log"
		}
		
		// 确保日志目录存在
		logDir := filepath.Dir(config.Filename)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return err
		}
		
		// 配置日志轮转
		lumberjackLogger := &lumberjack.Logger{
			Filename:   config.Filename,
			MaxSize:    config.MaxSize,
			MaxAge:     config.MaxAge,
			MaxBackups: config.MaxBackups,
			Compress:   config.Compress,
		}
		
		logger.SetOutput(lumberjackLogger)
	case "both":
		// 同时输出到控制台和文件
		if config.Filename == "" {
			config.Filename = "logs/app.log"
		}
		
		logDir := filepath.Dir(config.Filename)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return err
		}
		
		lumberjackLogger := &lumberjack.Logger{
			Filename:   config.Filename,
			MaxSize:    config.MaxSize,
			MaxAge:     config.MaxAge,
			MaxBackups: config.MaxBackups,
			Compress:   config.Compress,
		}
		
		// 创建多重写入器
		multiWriter := io.MultiWriter(os.Stdout, lumberjackLogger)
		logger.SetOutput(multiWriter)
	default:
		logger.SetOutput(os.Stdout)
	}
	
	return nil
}

// WithFields 添加字段
func (l *Logger) WithFields(fields map[string]interface{}) *logrus.Entry {
	return l.Logger.WithFields(fields)
}

// WithField 添加单个字段
func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.Logger.WithField(key, value)
}

// Close 关闭日志文件
func (l *Logger) Close() error {
	if closer, ok := l.Logger.Out.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
```

#### 3.3 创建日志中间件 (`pkg/logger/middleware.go`)
```go
package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GinLogger Gin日志中间件
func GinLogger(logger *Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 记录请求日志
		logger.WithFields(logrus.Fields{
			"timestamp":    param.TimeStamp.Format(time.RFC3339),
			"status":       param.StatusCode,
			"latency":      param.Latency,
			"client_ip":    param.ClientIP,
			"method":       param.Method,
			"path":         param.Path,
			"error":        param.ErrorMessage,
			"body_size":    param.BodySize,
			"user_agent":   param.Request.UserAgent(),
			"request_id":   param.Request.Header.Get("X-Request-ID"),
		}).Info("HTTP Request")
		
		return ""
	})
}

// GinRecovery Gin恢复中间件
func GinRecovery(logger *Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger.WithFields(logrus.Fields{
			"error":      recovered,
			"path":       c.Request.URL.Path,
			"method":     c.Request.Method,
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
			"request_id": c.Request.Header.Get("X-Request-ID"),
		}).Error("Panic recovered")
		
		c.AbortWithStatus(500)
	})
}

// RequestID 请求ID中间件
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
```

### 验证步骤

#### 3.4 创建日志测试文件 (`pkg/logger/logger_test.go`)
```go
package logger

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	config := &Config{
		Level:      "debug",
		Format:     "json",
		Output:     "stdout",
		Filename:   "",
		MaxSize:    100,
		MaxAge:     7,
		MaxBackups: 10,
		Compress:   true,
	}
	
	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	
	if logger == nil {
		t.Fatal("Logger is nil")
	}
	
	// 测试日志输出
	logger.Info("Test info message")
	logger.Debug("Test debug message")
	logger.Warn("Test warning message")
	logger.Error("Test error message")
}

func TestLoggerWithFile(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")
	
	config := &Config{
		Level:      "info",
		Format:     "text",
		Output:     "file",
		Filename:   logFile,
		MaxSize:    1,
		MaxAge:     1,
		MaxBackups: 3,
		Compress:   false,
	}
	
	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	
	// 写入日志
	logger.Info("Test file logging")
	logger.WithField("key", "value").Info("Test with field")
	
	// 关闭日志
	logger.Close()
	
	// 检查文件是否存在
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Errorf("Log file was not created: %s", logFile)
	}
}

func TestLoggerWithFields(t *testing.T) {
	config := &Config{
		Level:  "debug",
		Format: "json",
		Output: "stdout",
	}
	
	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	
	// 测试带字段的日志
	logger.WithFields(map[string]interface{}{
		"user_id":    123,
		"action":     "login",
		"ip_address": "192.168.1.1",
		"timestamp":  time.Now(),
	}).Info("User login")
	
	logger.WithField("request_id", "req-123").Error("Request failed")
}
```

#### 3.5 验证日志系统
```bash
# 运行日志测试
go test ./pkg/logger/

# 检查依赖
go mod tidy

# 创建日志目录
mkdir -p logs
```

### 预期结果
- 日志系统创建完成，支持多种输出格式
- 支持日志轮转和归档
- 包含完整的Gin中间件集成
- 支持结构化日志记录
- 包含完整的单元测试
- 依赖库正确添加到 go.mod

### 状态
- [ ] 步骤 3 完成：日志系统实现

---

## 下一步预告
步骤 4: 搭建HTTP服务器 (`cmd/server/main.go`)
- 初始化Gin框架
- 配置基础路由
- 集成日志和配置
- 优雅关闭机制
