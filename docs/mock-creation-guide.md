# Mock 创建流程详细文档

## 概述

根据 `internal/service/service_test.go` 文件中的注释，本文档详细说明如何为服务层测试创建 Mock 对象。该项目使用 testify/suite 测试框架，需要对 Repository 层和 Client 层进行 Mock。

## 1. 项目结构分析

### 1.1 服务层依赖关系
```
ServiceTestSuite
├── Manager (服务管理器)
│   ├── TaskService (任务服务)
│   ├── DetectService (检测服务)
│   └── ResultService (结果服务)
└── ServiceContext (服务上下文)
    ├── Logger (日志器)
    ├── Repository.Manager (仓储管理器)
    └── Client.Manager (客户端管理器)
```

### 1.2 需要 Mock 的组件

#### Repository 层组件
- `repository.Manager` - 仓储管理器
  - `TaskRepository` - 任务仓储
  - `ResultRepository` - 结果仓储  
  - `DeviceRepository` - 设备仓储
  - `WorkRepository` - 工作者仓储
  - `*gorm.DB` - 数据库连接
  - `*redis.Client` - Redis 客户端

#### Client 层组件
- `client.Manager` - 客户端管理器
  - `HttpClient` - HTTP 客户端
  - `DetectEngineClient` - 检测引擎客户端
  - `StorageClient` - 存储客户端
  - `MessageQueueClient` - 消息队列客户端

## 2. Mock 工具选择

推荐使用以下 Mock 工具：

### 2.1 GoMock (推荐)
```bash
# 安装 GoMock
go install github.com/golang/mock/mockgen@latest

# 或者使用 uber-go/mock (更活跃的分支)
go install go.uber.org/mock/mockgen@latest
```

### 2.2 Testify Mock (备选)
```bash
# 已经在项目中使用 testify/suite
go get github.com/stretchr/testify/mock
```

## 3. 详细 Mock 创建步骤

### 3.1 步骤一：为 Repository 接口生成 Mock

#### 3.1.1 创建 Mock 生成配置
在项目根目录创建 `tools.go` 文件：
```go
//go:build tools
// +build tools

package tools

import (
    _ "go.uber.org/mock/mockgen"
)
```

#### 3.1.2 生成 Repository Mock
```bash
# 生成仓储接口 Mock
mockgen -source=internal/repository/task_repository.go -destination=internal/repository/mocks/mock_task_repository.go -package=mocks
mockgen -source=internal/repository/result_repository.go -destination=internal/repository/mocks/mock_result_repository.go -package=mocks
mockgen -source=internal/repository/device_repository.go -destination=internal/repository/mocks/mock_device_repository.go -package=mocks
mockgen -source=internal/repository/work_repository.go -destination=internal/repository/mocks/mock_work_repository.go -package=mocks

# 生成仓储管理器 Mock (如果有接口定义)
mockgen -source=internal/repository/manager.go -destination=internal/repository/mocks/mock_manager.go -package=mocks
```

#### 3.1.3 手动创建 Repository Manager Mock
由于 `repository.Manager` 是结构体而非接口，需要创建接口或手动 Mock：

创建 `internal/repository/interfaces.go`：
```go
package repository

import (
    "github.com/redis/go-redis/v9"
    "gorm.io/gorm"
)

// ManagerInterface 仓储管理器接口
type ManagerInterface interface {
    GetDb() *gorm.DB
    GetRedis() *redis.Client
    Transaction(fn func(*Manager) error) error
    
    // 仓储访问器
    GetTaskRepository() TaskRepository
    GetResultRepository() ResultRepository
    GetDeviceRepository() DeviceRepository
    GetWorkerRepository() WorkRepository
}
```

然后生成 Mock：
```bash
mockgen -source=internal/repository/interfaces.go -destination=internal/repository/mocks/mock_interfaces.go -package=mocks
```

### 3.2 步骤二：为 Client 接口生成 Mock

#### 3.2.1 生成 Client Mock
```bash
# 生成客户端接口 Mock
mockgen -source=pkg/client/http_client.go -destination=pkg/client/mocks/mock_http_client.go -package=mocks
mockgen -source=pkg/client/detect_engine_client.go -destination=pkg/client/mocks/mock_detect_engine_client.go -package=mocks
mockgen -source=pkg/client/storage_client.go -destination=pkg/client/mocks/mock_storage_client.go -package=mocks
mockgen -source=pkg/client/message_queue_client.go -destination=pkg/client/mocks/mock_message_queue_client.go -package=mocks

# 生成客户端管理器 Mock
mockgen -source=pkg/client/manager.go -destination=pkg/client/mocks/mock_manager.go -package=mocks
```

#### 3.2.2 创建 Client Manager 接口
创建 `pkg/client/interfaces.go`：
```go
package client

import "context"

// ManagerInterface 客户端管理器接口
type ManagerInterface interface {
    Close() error
    HealthCheck(ctx context.Context) map[string]error
    
    // 客户端访问器
    GetHTTP() HttpClient
    GetDetectEngine() DetectEngineClient
    GetStorage() StorageClient
    GetMessageQueue() MessageQueueClient
}
```

### 3.3 步骤三：创建 Mock 辅助函数

创建 `internal/service/mocks/helpers.go`：
```go
package mocks

import (
    "detect-executor-go/internal/repository"
    "detect-executor-go/internal/repository/mocks"
    clientMocks "detect-executor-go/pkg/client/mocks"
    "detect-executor-go/pkg/client"
    "detect-executor-go/pkg/logger"
    
    "go.uber.org/mock/gomock"
)

// MockDependencies Mock 依赖集合
type MockDependencies struct {
    Ctrl *gomock.Controller
    
    // Repository Mocks
    RepoManager    *mocks.MockManagerInterface
    TaskRepo       *mocks.MockTaskRepository
    ResultRepo     *mocks.MockResultRepository
    DeviceRepo     *mocks.MockDeviceRepository
    WorkerRepo     *mocks.MockWorkRepository
    
    // Client Mocks
    ClientManager  *clientMocks.MockManagerInterface
    HTTPClient     *clientMocks.MockHttpClient
    DetectEngine   *clientMocks.MockDetectEngineClient
    Storage        *clientMocks.MockStorageClient
    MessageQueue   *clientMocks.MockMessageQueueClient
}

// NewMockDependencies 创建所有 Mock 依赖
func NewMockDependencies(ctrl *gomock.Controller) *MockDependencies {
    return &MockDependencies{
        Ctrl: ctrl,
        
        // Repository Mocks
        RepoManager: mocks.NewMockManagerInterface(ctrl),
        TaskRepo:    mocks.NewMockTaskRepository(ctrl),
        ResultRepo:  mocks.NewMockResultRepository(ctrl),
        DeviceRepo:  mocks.NewMockDeviceRepository(ctrl),
        WorkerRepo:  mocks.NewMockWorkRepository(ctrl),
        
        // Client Mocks
        ClientManager: clientMocks.NewMockManagerInterface(ctrl),
        HTTPClient:    clientMocks.NewMockHttpClient(ctrl),
        DetectEngine:  clientMocks.NewMockDetectEngineClient(ctrl),
        Storage:       clientMocks.NewMockStorageClient(ctrl),
        MessageQueue:  clientMocks.NewMockMessageQueueClient(ctrl),
    }
}

// SetupRepositoryExpectations 设置仓储层基础期望
func (m *MockDependencies) SetupRepositoryExpectations() {
    m.RepoManager.EXPECT().GetTaskRepository().Return(m.TaskRepo).AnyTimes()
    m.RepoManager.EXPECT().GetResultRepository().Return(m.ResultRepo).AnyTimes()
    m.RepoManager.EXPECT().GetDeviceRepository().Return(m.DeviceRepo).AnyTimes()
    m.RepoManager.EXPECT().GetWorkerRepository().Return(m.WorkerRepo).AnyTimes()
}

// SetupClientExpectations 设置客户端层基础期望
func (m *MockDependencies) SetupClientExpectations() {
    m.ClientManager.EXPECT().GetHTTP().Return(m.HTTPClient).AnyTimes()
    m.ClientManager.EXPECT().GetDetectEngine().Return(m.DetectEngine).AnyTimes()
    m.ClientManager.EXPECT().GetStorage().Return(m.Storage).AnyTimes()
    m.ClientManager.EXPECT().GetMessageQueue().Return(m.MessageQueue).AnyTimes()
}
```

### 3.4 步骤四：更新测试套件

更新 `internal/service/service_test.go`：
```go
package service

import (
    "context"
    "detect-executor-go/internal/model"
    "detect-executor-go/internal/service/mocks"
    "detect-executor-go/pkg/logger"
    "testing"

    "github.com/stretchr/testify/suite"
    "go.uber.org/mock/gomock"
)

type ServiceTestSuite struct {
    suite.Suite
    ctx     context.Context
    manager *Manager
    
    // Mock 依赖
    ctrl  *gomock.Controller
    mocks *mocks.MockDependencies
}

func (suite *ServiceTestSuite) SetupSuite() {
    suite.ctx = context.Background()
    
    // 创建 Mock 控制器
    suite.ctrl = gomock.NewController(suite.T())
    
    // 创建所有 Mock 依赖
    suite.mocks = mocks.NewMockDependencies(suite.ctrl)
    
    // 设置基础期望
    suite.mocks.SetupRepositoryExpectations()
    suite.mocks.SetupClientExpectations()
    
    // 创建服务上下文
    logger := logger.NewLogger()
    ctx := &ServiceContext{
        Logger:     logger,
        Repository: suite.mocks.RepoManager,
        Client:     suite.mocks.ClientManager,
    }
    
    // 创建服务管理器
    suite.manager = &Manager{
        ctx: ctx,
    }
    suite.manager.Task = NewTaskService(ctx)
    suite.manager.Detect = NewDetectService(ctx)
    suite.manager.Result = NewResultService(ctx)
}

func (suite *ServiceTestSuite) TearDownSuite() {
    if suite.ctrl != nil {
        suite.ctrl.Finish()
    }
}

// TestTaskService 测试任务服务
func (suite *ServiceTestSuite) TestTaskService() {
    // 设置 Mock 期望
    deviceID := "device-001"
    
    // Mock 设备存在且在线
    suite.mocks.DeviceRepo.EXPECT().
        GetByID(suite.ctx, deviceID).
        Return(&model.Device{
            ID:     deviceID,
            Status: model.DeviceStatusOnline,
        }, nil)
    
    // Mock 创建任务
    suite.mocks.TaskRepo.EXPECT().
        Create(suite.ctx, gomock.Any()).
        Return(&model.DetectTask{
            ID:       "task-001",
            DeviceID: deviceID,
            Type:     model.DetectTypeQuality,
            Status:   model.TaskStatusPending,
        }, nil)
    
    // 测试创建任务
    req := &CreateTaskRequest{
        DeviceID:    deviceID,
        Type:        model.DetectTypeQuality,
        VideoURL:    "rtsp://example.com/stream",
        Priority:    model.PriorityHigh,
        Description: "Test task",
    }

    task, err := suite.manager.Task.CreateTask(suite.ctx, req)
    suite.NoError(err)
    suite.NotNil(task)
    suite.Equal(req.DeviceID, task.DeviceID)
    suite.Equal(req.Type, task.Type)
}

// TestDetectService 测试检测服务
func (suite *ServiceTestSuite) TestDetectService() {
    taskID := "test-task-001"
    
    // Mock 任务存在
    suite.mocks.TaskRepo.EXPECT().
        GetByID(suite.ctx, taskID).
        Return(&model.DetectTask{
            ID:     taskID,
            Status: model.TaskStatusPending,
        }, nil)
    
    // Mock 检测引擎调用
    suite.mocks.DetectEngine.EXPECT().
        StartDetection(suite.ctx, gomock.Any()).
        Return(nil)
    
    // Mock 更新任务状态
    suite.mocks.TaskRepo.EXPECT().
        UpdateStatus(suite.ctx, taskID, model.TaskStatusRunning, gomock.Any()).
        Return(nil)

    // 测试开始检测
    err := suite.manager.Detect.StartDetection(suite.ctx, taskID)
    suite.NoError(err)

    // 测试获取检测状态
    status, err := suite.manager.Detect.GetDetectionStatus(suite.ctx, taskID)
    suite.NoError(err)
    suite.NotNil(status)
    suite.Equal(taskID, status.TaskID)
}

// TestResultService 测试结果服务
func (suite *ServiceTestSuite) TestResultService() {
    // Mock 数据
    expectedResults := []model.DetectResult{
        {
            ID:     "result-001",
            TaskID: "task-001",
            Score:  0.95,
        },
    }
    
    suite.mocks.ResultRepo.EXPECT().
        List(suite.ctx, gomock.Any()).
        Return(expectedResults, nil)

    // 测试列表结果
    req := &ListResultsRequest{
        Page:     1,
        PageSize: 10,
    }

    results, total, err := suite.manager.Result.ListResults(suite.ctx, req)
    suite.NoError(err)
    suite.NotNil(results)
    suite.Equal(len(expectedResults), len(results))
    suite.GreaterOrEqual(total, int64(0))
}

// TestServiceTestSuite 运行测试套件
func TestServiceTestSuite(t *testing.T) {
    suite.Run(t, new(ServiceTestSuite))
}
```

## 4. 高级 Mock 技巧

### 4.1 条件期望设置
```go
// 根据不同条件设置不同的 Mock 行为
suite.mocks.TaskRepo.EXPECT().
    GetByID(suite.ctx, "existing-task").
    Return(&model.DetectTask{ID: "existing-task"}, nil)

suite.mocks.TaskRepo.EXPECT().
    GetByID(suite.ctx, "non-existing-task").
    Return(nil, errors.New("task not found"))
```

### 4.2 参数匹配器
```go
// 使用 gomock 匹配器
suite.mocks.TaskRepo.EXPECT().
    Create(suite.ctx, gomock.Any()).
    DoAndReturn(func(ctx context.Context, task *model.DetectTask) (*model.DetectTask, error) {
        // 自定义验证逻辑
        suite.NotEmpty(task.DeviceID)
        task.ID = "generated-id"
        return task, nil
    })
```

### 4.3 调用次数控制
```go
// 精确控制调用次数
suite.mocks.DetectEngine.EXPECT().
    StartDetection(suite.ctx, gomock.Any()).
    Times(1)

// 允许多次调用
suite.mocks.DeviceRepo.EXPECT().
    GetByID(suite.ctx, gomock.Any()).
    Return(&model.Device{}, nil).
    AnyTimes()
```

## 5. 最佳实践

### 5.1 Mock 组织原则
1. **单一职责**：每个测试方法只测试一个功能点
2. **最小化 Mock**：只 Mock 必要的依赖
3. **明确期望**：清晰设置 Mock 的输入输出期望
4. **及时清理**：在 TearDown 中清理 Mock 资源

### 5.2 测试数据管理
```go
// 创建测试数据工厂
func CreateTestTask() *model.DetectTask {
    return &model.DetectTask{
        ID:       "test-task-" + uuid.New().String(),
        DeviceID: "test-device-001",
        Type:     model.DetectTypeQuality,
        Status:   model.TaskStatusPending,
        CreatedAt: time.Now(),
    }
}
```

### 5.3 错误场景测试
```go
func (suite *ServiceTestSuite) TestTaskService_CreateTask_DeviceNotFound() {
    // Mock 设备不存在的场景
    suite.mocks.DeviceRepo.EXPECT().
        GetByID(suite.ctx, "non-existing-device").
        Return(nil, errors.New("device not found"))
    
    req := &CreateTaskRequest{
        DeviceID: "non-existing-device",
        Type:     model.DetectTypeQuality,
    }
    
    task, err := suite.manager.Task.CreateTask(suite.ctx, req)
    suite.Error(err)
    suite.Nil(task)
    suite.Contains(err.Error(), "device not found")
}
```

## 6. 执行测试

### 6.1 运行测试命令
```bash
# 运行所有服务测试
go test ./internal/service/... -v

# 运行特定测试套件
go test ./internal/service -run TestServiceTestSuite -v

# 运行特定测试方法
go test ./internal/service -run TestServiceTestSuite/TestTaskService -v

# 生成测试覆盖率报告
go test ./internal/service/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### 6.2 CI/CD 集成
```yaml
# .github/workflows/test.yml
name: Test
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - uses: actions/setup-go@v2
      with:
        go-version: 1.21
    - name: Run tests
      run: |
        go mod download
        go generate ./...
        go test ./... -v -coverprofile=coverage.out
    - name: Upload coverage
      uses: codecov/codecov-action@v1
      with:
        file: ./coverage.out
```

## 7. 故障排除

### 7.1 常见问题
1. **Mock 未生成**：检查 mockgen 命令路径和参数
2. **接口不匹配**：确保 Mock 接口与实际接口一致
3. **期望未满足**：检查 Mock 期望设置是否正确
4. **循环依赖**：重新设计接口避免循环引用

### 7.2 调试技巧
```go
// 启用 Mock 调用日志
suite.ctrl = gomock.NewController(suite.T())
suite.ctrl.RecordCall = true // 记录所有调用

// 在测试中打印 Mock 调用
suite.mocks.TaskRepo.EXPECT().
    Create(suite.ctx, gomock.Any()).
    Do(func(ctx context.Context, task *model.DetectTask) {
        suite.T().Logf("Creating task: %+v", task)
    }).
    Return(&model.DetectTask{}, nil)
```

## 8. 总结

通过以上步骤，您可以为服务层创建完整的 Mock 测试环境：

1. ✅ 分析依赖关系，识别需要 Mock 的组件
2. ✅ 使用 GoMock 生成接口 Mock
3. ✅ 创建 Mock 辅助函数简化测试设置
4. ✅ 更新测试套件使用 Mock 依赖
5. ✅ 编写具体的测试用例
6. ✅ 遵循最佳实践确保测试质量

这样可以实现对服务层的完全隔离测试，提高测试的可靠性和执行速度。
