# Detect-Executor-Go 检测执行模块

## 项目简介

Detect-Executor-Go 是一个基于 Go 语言开发的高性能视频设备检测执行模块，专门用于大规模视频监控设备的质量检测。该项目从原有的 Java 检测系统中抽离出核心检测功能，采用 Go 语言重新设计实现，旨在提供更高的并发性能和更简洁的架构设计。

### 核心特性

- 🚀 **高并发处理**：基于 Go 协程，支持数千设备同时检测
- 🎯 **简洁架构**：微服务设计，职责单一，易于维护
- 📊 **实时监控**：完整的指标监控和健康检查
- 🔄 **智能重试**：指数退避重试机制，提高检测成功率
- 💾 **高效缓存**：Redis 缓存热点数据，减少数据库压力
- 🐳 **容器化部署**：Docker 支持，便于部署和扩展

## 技术架构

### 整体架构图

```
┌─────────────────┐    HTTP     ┌─────────────────┐    HTTP     ┌─────────────────┐
│   Java 主服务    │ ────────► │ Go 检测执行模块  │ ────────► │  C++ 检测引擎   │
│  (原有系统)     │            │                │            │                │
└─────────────────┘            └─────────────────┘            └─────────────────┘
                                        │                              
                                        ▼                              
                               ┌─────────────────┐                     
                               │     Redis       │                     
                               │   (缓存层)      │                     
                               └─────────────────┘                     
                                        │                              
                                        ▼                              
                               ┌─────────────────┐                     
                               │     MySQL       │                     
                               │   (结果存储)    │                     
                               └─────────────────┘                     
```

### 技术栈

| 组件 | 技术选型 | 说明 |
|------|----------|------|
| Web框架 | Gin | 高性能HTTP框架 |
| 数据库ORM | GORM | 功能完善的Go ORM |
| 缓存 | go-redis | Redis官方Go客户端 |
| HTTP客户端 | resty | 简洁易用的HTTP客户端 |
| 配置管理 | viper | 灵活的配置管理库 |
| 日志 | logrus | 结构化日志库 |
| 监控 | prometheus | 指标监控 |
| 容器化 | Docker | 容器化部署 |

## 功能实现

### 核心功能模块

#### 1. 检测任务调度
- **批量任务处理**：支持批量设备检测请求
- **并发控制**：基于 Worker Pool 模式的并发控制
- **任务队列**：使用 Channel 实现高效任务队列
- **进度跟踪**：实时更新检测进度和状态

#### 2. 设备检测引擎
- **视频流获取**：调用设备接口获取 RTMP 视频流
- **质量检测**：调用 C++ 检测引擎进行视频质量分析
- **多维度检测**：支持模糊、对比度、遮挡等多种检测类型
- **结果解析**：解析检测引擎返回结果并标准化

#### 3. 缓存优化
- **设备信息缓存**：缓存设备基础信息，减少查询
- **检测配置缓存**：缓存检测参数配置
- **结果缓存**：临时缓存检测结果，支持批量写入

#### 4. 数据持久化
- **结果存储**：检测结果写入 MySQL 数据库
- **批量操作**：支持批量插入和更新操作
- **事务管理**：确保数据一致性

### 项目亮点

#### 1. 高并发架构设计
```go
// Worker Pool 模式实现高并发
type DetectEngine struct {
    maxWorkers   int
    taskQueue    chan DetectTask
    resultQueue  chan DetectResult
    workerPool   sync.WaitGroup
}

// 支持数千并发检测
func (e *DetectEngine) Start() {
    for i := 0; i < e.maxWorkers; i++ {
        go e.worker()
    }
}
```

#### 2. 智能重试机制
```go
// 指数退避重试策略
type RetryConfig struct {
    MaxRetries int           `yaml:"max_retries"`
    BaseDelay  time.Duration `yaml:"base_delay"`
    MaxDelay   time.Duration `yaml:"max_delay"`
}
```

#### 3. 完整的可观测性
- Prometheus 指标监控
- 结构化日志记录
- 健康检查接口
- 性能指标统计

#### 4. 云原生支持
- Docker 容器化
- Kubernetes 部署配置
- 配置热更新
- 优雅关闭

## 项目目录结构

```
detect-executor-go/
├── cmd/                          # 应用程序入口
│   └── server/
│       └── main.go              # 主程序入口，初始化服务
├── internal/                     # 内部包，不对外暴露
│   ├── handler/                 # HTTP 处理器层
│   │   ├── base.go             # 基础处理器，通用响应方法
│   │   ├── detect_handler.go   # 检测相关接口处理
│   │   ├── health_handerl.go   # 健康检查接口
│   │   ├── result_handler.go   # 结果查询接口处理
│   │   ├── task_handler.go     # 任务管理接口处理
│   │   ├── router.go           # 路由配置和中间件
│   │   └── handle_test.go      # 处理器测试
│   ├── service/                 # 业务逻辑层
│   │   ├── base.go             # 服务基础结构
│   │   ├── detect_service.go   # 检测业务逻辑
│   │   ├── task_service.go     # 任务管理服务
│   │   ├── result_service.go   # 结果查询服务
│   │   ├── manager.go          # 服务管理器
│   │   └── service_test.go     # 服务层测试
│   ├── repository/              # 数据访问层
│   │   ├── base.go             # 仓库基础接口
│   │   ├── device_repository.go # 设备信息数据访问
│   │   ├── task_repository.go  # 任务数据访问
│   │   ├── result_repository.go # 检测结果数据访问
│   │   ├── work_repository.go  # 工作节点数据访问
│   │   ├── manager.go          # 仓库管理器
│   │   └── repository_test.go  # 仓库层测试
│   └── model/                   # 数据模型
│       ├── base.go             # 基础模型结构
│       ├── detect_task.go      # 检测任务模型
│       ├── detect_result.go    # 检测结果模型
│       ├── device.go           # 设备模型
│       ├── worker.go           # 工作节点模型
│       ├── dto.go              # 数据传输对象
│       └── model_test.go       # 模型测试
├── pkg/                         # 公共包，可对外暴露
│   ├── config/                 # 配置管理
│   │   ├── config.go           # 配置结构定义
│   │   ├── loader.go           # 配置加载器
│   │   ├── config.yaml         # 默认配置文件
│   │   ├── config.dev.yaml     # 开发环境配置
│   │   ├── config.prod.yaml    # 生产环境配置
│   │   └── config_test.go      # 配置测试
│   ├── logger/                 # 日志管理
│   │   ├── logger.go           # 日志初始化
│   │   ├── middleware.go       # 日志中间件
│   │   └── logger_test.go      # 日志测试
│   ├── database/               # 数据库连接
│   │   ├── mysql.go            # MySQL 连接池
│   │   ├── redis.go            # Redis 连接池
│   │   ├── manager.go          # 数据库管理器
│   │   └── database_test.go    # 数据库测试
│   └── client/                 # 外部服务客户端
│       ├── base.go             # 客户端基础结构
│       ├── detect_engine_client.go # 检测引擎客户端
│       ├── http_client.go      # HTTP 客户端封装
│       ├── message_queue_client.go # 消息队列客户端
│       ├── storage_client.go   # 存储服务客户端
│       ├── manager.go          # 客户端管理器
│       └── client_test.go      # 客户端测试
├── configs/                     # 配置文件目录
│   ├── config.yaml             # 主配置文件
│   ├── config.dev.yaml         # 开发环境配置
│   └── config.prod.yaml        # 生产环境配置
├── scripts/                     # 脚本文件
│   ├── dev.sh                  # 开发环境启动脚本
│   ├── prod.sh                 # 生产环境启动脚本
│   ├── start.sh                # 通用启动脚本
│   ├── download-images.sh      # 镜像下载脚本
│   └── load-images.sh          # 镜像加载脚本
├── docs/                        # 项目文档
│   ├── detect-executor-go-doc-1-3.md  # 核心实现文档
│   ├── detect-executor-go-doc-4.md    # API接口文档
│   ├── detect-executor-go-doc-5.md    # 数据库设计文档
│   ├── detect-executor-go-doc-6.md    # 配置管理文档
│   ├── detect-executor-go-doc-7.md    # 部署运维文档
│   ├── detect-executor-go-doc-8.md    # 性能优化文档
│   ├── detect-executor-go-doc-9.md    # 监控告警文档
│   ├── detect-executor-go-doc-10.md   # 故障排查文档
│   ├── detect-executor-go-doc-11.md   # 开发指南文档
│   └── mock-creation-guide.md         # Mock测试指南
├── Dockerfile                   # Docker 镜像构建文件
├── docker-compose.yml          # Docker Compose 配置
├── Makefile                     # 构建配置
├── go.mod                       # Go 模块定义
├── go.sum                       # Go 模块校验和
└── README.md                    # 项目说明文档
```

### 核心文件功能说明

#### 应用入口层
- **`cmd/server/main.go`**: 程序主入口，负责初始化配置、数据库连接、启动HTTP服务器

#### HTTP处理器层
- **`internal/handler/base.go`**: 基础处理器，提供通用响应方法和参数验证
- **`internal/handler/task_handler.go`**: 任务管理相关API接口处理
- **`internal/handler/detect_handler.go`**: 检测控制相关API接口处理
- **`internal/handler/result_handler.go`**: 结果查询相关API接口处理
- **`internal/handler/health_handerl.go`**: 健康检查和监控接口处理
- **`internal/handler/router.go`**: 路由配置、中间件注册和API分组

#### 业务逻辑层
- **`internal/service/task_service.go`**: 任务管理业务逻辑，包括任务创建、调度和状态管理
- **`internal/service/detect_service.go`**: 检测执行业务逻辑，协调检测引擎和结果处理
- **`internal/service/result_service.go`**: 结果查询和统计业务逻辑
- **`internal/service/manager.go`**: 服务管理器，统一管理各业务服务的生命周期

#### 数据访问层
- **`internal/repository/task_repository.go`**: 检测任务的CRUD操作和查询
- **`internal/repository/result_repository.go`**: 检测结果的存储和查询操作
- **`internal/repository/device_repository.go`**: 设备信息的数据访问操作
- **`internal/repository/work_repository.go`**: 工作节点信息的数据访问
- **`internal/repository/manager.go`**: 仓库管理器，统一管理数据访问层

#### 数据模型层
- **`internal/model/detect_task.go`**: 检测任务数据模型和业务实体
- **`internal/model/detect_result.go`**: 检测结果数据模型和统计结构
- **`internal/model/device.go`**: 设备信息模型和设备状态定义
- **`internal/model/worker.go`**: 工作节点模型和节点状态管理
- **`internal/model/dto.go`**: 数据传输对象，用于API请求响应

#### 外部客户端层
- **`pkg/client/detect_engine_client.go`**: 检测引擎HTTP客户端封装
- **`pkg/client/http_client.go`**: 通用HTTP客户端，支持重试和超时配置
- **`pkg/client/message_queue_client.go`**: 消息队列客户端，用于异步任务通信
- **`pkg/client/storage_client.go`**: 存储服务客户端，用于文件和媒体存储
- **`pkg/client/manager.go`**: 客户端管理器，统一管理外部服务连接

#### 基础设施层
- **`pkg/config/config.go`**: 配置结构定义和环境变量映射
- **`pkg/database/mysql.go`**: MySQL连接池和事务管理
- **`pkg/database/redis.go`**: Redis连接池和缓存操作
- **`pkg/logger/logger.go`**: 结构化日志初始化和配置
- **`pkg/logger/middleware.go`**: HTTP请求日志中间件和错误恢复

## 测试流程

### 单元测试
```bash
# 运行所有单元测试
make test

# 运行特定包的测试
go test ./internal/service/...

# 生成测试覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 集成测试
```bash
# 启动测试环境（Docker Compose）
docker-compose -f docker-compose.test.yml up -d

# 运行集成测试
make integration-test

# 清理测试环境
docker-compose -f docker-compose.test.yml down
```

### 性能测试
```bash
# 使用 Apache Bench 进行压力测试
ab -n 1000 -c 100 http://localhost:8080/api/v1/detect/batch

# 使用 wrk 进行性能测试
wrk -t12 -c400 -d30s http://localhost:8080/api/v1/health
```

## 部署流程

### 本地开发部署
```bash
# 1. 克隆项目
git clone https://github.com/your-username/detect-executor-go.git
cd detect-executor-go

# 2. 安装依赖
go mod download

# 3. 启动依赖服务（MySQL, Redis）
docker-compose up -d mysql redis

# 4. 运行应用
make run
```

### Docker 部署
```bash
# 1. 构建镜像
make docker-build

# 2. 运行容器
docker run -d \
  --name detect-executor \
  -p 8080:8080 \
  -e CONFIG_PATH=/app/configs/config.prod.yaml \
  detect-executor-go:latest
```

### Kubernetes 部署
```bash
# 1. 应用配置
kubectl apply -f deployments/k8s/configmap.yaml

# 2. 部署应用
kubectl apply -f deployments/k8s/deployment.yaml

# 3. 创建服务
kubectl apply -f deployments/k8s/service.yaml

# 4. 检查部署状态
kubectl get pods -l app=detect-executor
```

### CI/CD 流程
```yaml
# .github/workflows/ci.yml
name: CI/CD Pipeline
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
      - run: make test
  
  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - run: make docker-build
      - run: docker push your-registry/detect-executor-go
```

## 最佳编程建议

### 开发顺序建议

#### 第一阶段：基础框架搭建 📖 [详细文档](docs/detect-executor-go-doc-1-3.md)
1. **从 `go.mod` 开始**
   ```bash
   go mod init detect-executor-go
   ```

2. **创建配置管理** (`pkg/config/config.go`)
   - 定义配置结构体
   - 实现配置加载逻辑

3. **实现日志系统** (`pkg/logger/logger.go`)
   - 配置结构化日志
   - 实现日志中间件

4. **搭建HTTP服务器** (`cmd/server/main.go`) 📖 [文档](docs/detect-executor-go-doc-4.md)
   - 初始化Gin框架
   - 配置基础路由

#### 第二阶段：数据层实现 
5. **数据库连接** (`pkg/database/mysql.go`, `pkg/database/redis.go`) 📖 [文档](docs/detect-executor-go-doc-5.md)
   - 实现连接池管理
   - 配置连接参数

6. **数据模型定义** (`internal/model/`) 📖 [文档](docs/detect-executor-go-doc-6.md)
   - 定义检测相关数据结构
   - 实现数据验证

7. **数据访问层** (`internal/repository/`) 📖 [文档](docs/detect-executor-go-doc-7.md)
   - 实现基础CRUD操作
   - 添加缓存逻辑

#### 第三阶段：业务逻辑实现
8. **外部客户端** (`pkg/client/`) 📖 [文档](docs/detect-executor-go-doc-8.md)
   - 实现检测引擎客户端
   - 添加重试机制

9. **核心业务服务** (`internal/service/`) 📖 [文档](docs/detect-executor-go-doc-9.md)
   - 实现检测引擎
   - 添加并发控制

10. **HTTP处理器** (`internal/handler/`) 📖 [文档](docs/detect-executor-go-doc-10.md)
    - 实现API接口
    - 添加参数验证

#### 第四阶段：完善功能
11. **监控和指标** (`pkg/metrics/`) 📖 [文档](docs/detect-executor-go-doc-11.md)
    - 添加Prometheus指标
    - 实现健康检查

12. **测试用例** 📖 [Mock测试指南](docs/mock-creation-guide.md)
    - 编写单元测试
    - 实现集成测试

### 编程最佳实践

#### 1. 错误处理
```go
// 使用包装错误，提供上下文信息
func (s *DetectService) ProcessDevice(deviceID string) error {
    device, err := s.repo.GetDevice(deviceID)
    if err != nil {
        return fmt.Errorf("failed to get device %s: %w", deviceID, err)
    }
    // 业务逻辑...
    return nil
}
```

#### 2. 并发安全
```go
// 使用sync.Pool 复用对象
var detectTaskPool = sync.Pool{
    New: func() interface{} {
        return &DetectTask{}
    },
}

func getDetectTask() *DetectTask {
    return detectTaskPool.Get().(*DetectTask)
}
```

#### 3. 配置管理
```go
// 使用环境变量覆盖配置
type Config struct {
    Server   ServerConfig   `yaml:"server"`
    Database DatabaseConfig `yaml:"database"`
    Redis    RedisConfig    `yaml:"redis"`
}

// 支持配置热更新
func (c *Config) Watch() {
    viper.WatchConfig()
    viper.OnConfigChange(func(e fsnotify.Event) {
        // 重新加载配置
    })
}
```

#### 4. 优雅关闭
```go
// 实现优雅关闭
func (s *Server) Shutdown(ctx context.Context) error {
    // 停止接收新请求
    if err := s.httpServer.Shutdown(ctx); err != nil {
        return err
    }
    
    // 等待现有任务完成
    s.detectEngine.Stop()
    
    // 关闭数据库连接
    return s.db.Close()
}
```

### 性能优化建议

1. **使用对象池减少GC压力**
2. **合理设置GOMAXPROCS**
3. **使用pprof进行性能分析**
4. **实现连接池复用**
5. **批量操作减少数据库访问**

### 代码质量保证

1. **使用golangci-lint进行代码检查**
2. **保持测试覆盖率在80%以上**
3. **编写清晰的API文档**
4. **使用Git hooks进行代码检查**
5. **定期进行代码审查**

通过以上的架构设计和实现建议，你可以构建一个高性能、可扩展的检测执行模块，既满足了技术需求，又能作为优秀的作品集项目展示你的Go语言开发能力和系统架构设计思维。
