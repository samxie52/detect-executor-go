# 步骤 7: 数据访问层实现

## 目标
实现数据访问层（Repository Pattern），提供统一的数据访问接口，支持MySQL和Redis缓存集成，实现高效的数据查询和操作。

## 操作步骤

### 7.1 创建基础仓储接口 (`internal/repository/base.go`)
```go
package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"github.com/redis/go-redis/v9"
)

// BaseRepository 基础仓储接口
type BaseRepository interface {
	GetDB() *gorm.DB
	GetRedis() *redis.Client
	WithContext(ctx context.Context) BaseRepository
	WithTx(tx *gorm.DB) BaseRepository
}

// baseRepository 基础仓储实现
type baseRepository struct {
	db    *gorm.DB
	redis *redis.Client
	ctx   context.Context
}

// NewBaseRepository 创建基础仓储
func NewBaseRepository(db *gorm.DB, redis *redis.Client) BaseRepository {
	return &baseRepository{
		db:    db,
		redis: redis,
		ctx:   context.Background(),
	}
}

func (r *baseRepository) GetDB() *gorm.DB {
	return r.db.WithContext(r.ctx)
}

func (r *baseRepository) GetRedis() *redis.Client {
	return r.redis
}

func (r *baseRepository) WithContext(ctx context.Context) BaseRepository {
	return &baseRepository{
		db:    r.db,
		redis: r.redis,
		ctx:   ctx,
	}
}

func (r *baseRepository) WithTx(tx *gorm.DB) BaseRepository {
	return &baseRepository{
		db:    tx,
		redis: r.redis,
		ctx:   r.ctx,
	}
}

// CacheConfig 缓存配置
type CacheConfig struct {
	TTL    time.Duration
	Prefix string
}

// DefaultCacheConfig 默认缓存配置
var DefaultCacheConfig = CacheConfig{
	TTL:    30 * time.Minute,
	Prefix: "detect:",
}
```

### 7.2 创建任务仓储接口 (`internal/repository/task_repository.go`)
```go
package repository

import (
	"context"
	"detect-executor-go/internal/model"
)

// TaskRepository 任务仓储接口
type TaskRepository interface {
	BaseRepository
	
	// 基础CRUD操作
	Create(ctx context.Context, task *model.DetectTask) error
	GetByID(ctx context.Context, id uint) (*model.DetectTask, error)
	GetByTaskID(ctx context.Context, taskID string) (*model.DetectTask, error)
	Update(ctx context.Context, task *model.DetectTask) error
	Delete(ctx context.Context, id uint) error
	
	// 查询操作
	List(ctx context.Context, req *model.TaskListRequest) ([]model.DetectTask, int64, error)
	GetByDeviceID(ctx context.Context, deviceID string, limit int) ([]model.DetectTask, error)
	GetByStatus(ctx context.Context, status model.TaskStatus, limit int) ([]model.DetectTask, error)
	GetPendingTasks(ctx context.Context, limit int) ([]model.DetectTask, error)
	GetRetryableTasks(ctx context.Context, limit int) ([]model.DetectTask, error)
	
	// 状态更新
	UpdateStatus(ctx context.Context, taskID string, status model.TaskStatus) error
	UpdateWorker(ctx context.Context, taskID string, workerID string) error
	IncrementRetry(ctx context.Context, taskID string) error
	
	// 统计查询
	CountByStatus(ctx context.Context, status model.TaskStatus) (int64, error)
	CountByDeviceID(ctx context.Context, deviceID string) (int64, error)
	GetTaskStats(ctx context.Context, deviceID string, days int) (map[string]int64, error)
}
```

### 7.3 创建结果仓储接口 (`internal/repository/result_repository.go`)
```go
package repository

import (
	"context"
	"detect-executor-go/internal/model"
)

// ResultRepository 结果仓储接口
type ResultRepository interface {
	BaseRepository
	
	// 基础CRUD操作
	Create(ctx context.Context, result *model.DetectResult) error
	GetByID(ctx context.Context, id uint) (*model.DetectResult, error)
	GetByResultID(ctx context.Context, resultID string) (*model.DetectResult, error)
	Update(ctx context.Context, result *model.DetectResult) error
	Delete(ctx context.Context, id uint) error
	
	// 查询操作
	GetByTaskID(ctx context.Context, taskID string) (*model.DetectResult, error)
	GetByDeviceID(ctx context.Context, deviceID string, limit int) ([]model.DetectResult, error)
	List(ctx context.Context, deviceID string, detectType model.DetectType, page, size int) ([]model.DetectResult, int64, error)
	
	// 统计查询
	CountByDeviceID(ctx context.Context, deviceID string) (int64, error)
	GetResultStats(ctx context.Context, deviceID string, days int) (map[string]interface{}, error)
	GetQualityTrend(ctx context.Context, deviceID string, days int) ([]map[string]interface{}, error)
}
```

### 7.4 创建设备仓储接口 (`internal/repository/device_repository.go`)
```go
package repository

import (
	"context"
	"detect-executor-go/internal/model"
)

// DeviceRepository 设备仓储接口
type DeviceRepository interface {
	BaseRepository
	
	// 基础CRUD操作
	Create(ctx context.Context, device *model.Device) error
	GetByID(ctx context.Context, id uint) (*model.Device, error)
	GetByDeviceID(ctx context.Context, deviceID string) (*model.Device, error)
	Update(ctx context.Context, device *model.Device) error
	Delete(ctx context.Context, id uint) error
	
	// 查询操作
	List(ctx context.Context, page, size int) ([]model.Device, int64, error)
	GetOnlineDevices(ctx context.Context) ([]model.Device, error)
	GetByStatus(ctx context.Context, status string) ([]model.Device, error)
	
	// 状态更新
	UpdateLastSeen(ctx context.Context, deviceID string) error
	UpdateStatus(ctx context.Context, deviceID string, status string) error
}
```

### 7.5 创建工作节点仓储接口 (`internal/repository/worker_repository.go`)
```go
package repository

import (
	"context"
	"detect-executor-go/internal/model"
)

// WorkerRepository 工作节点仓储接口
type WorkerRepository interface {
	BaseRepository
	
	// 基础CRUD操作
	Create(ctx context.Context, worker *model.Worker) error
	GetByID(ctx context.Context, id uint) (*model.Worker, error)
	GetByWorkerID(ctx context.Context, workerID string) (*model.Worker, error)
	Update(ctx context.Context, worker *model.Worker) error
	Delete(ctx context.Context, id uint) error
	
	// 查询操作
	List(ctx context.Context, page, size int) ([]model.Worker, int64, error)
	GetOnlineWorkers(ctx context.Context) ([]model.Worker, error)
	GetAvailableWorkers(ctx context.Context) ([]model.Worker, error)
	GetByStatus(ctx context.Context, status string) ([]model.Worker, error)
	
	// 状态更新
	UpdateLastSeen(ctx context.Context, workerID string) error
	UpdateStatus(ctx context.Context, workerID string, status string) error
	UpdateStats(ctx context.Context, workerID string, stats map[string]interface{}) error
}
```

### 7.6 创建仓储管理器 (`internal/repository/manager.go`)
```go
package repository

import (
	"gorm.io/gorm"
	"github.com/redis/go-redis/v9"
)

// Manager 仓储管理器
type Manager struct {
	db    *gorm.DB
	redis *redis.Client
	
	Task   TaskRepository
	Result ResultRepository
	Device DeviceRepository
	Worker WorkerRepository
}

// NewManager 创建仓储管理器
func NewManager(db *gorm.DB, redis *redis.Client) *Manager {
	return &Manager{
		db:    db,
		redis: redis,
		
		Task:   NewTaskRepository(db, redis),
		Result: NewResultRepository(db, redis),
		Device: NewDeviceRepository(db, redis),
		Worker: NewWorkerRepository(db, redis),
	}
}

// GetDB 获取数据库连接
func (m *Manager) GetDB() *gorm.DB {
	return m.db
}

// GetRedis 获取Redis连接
func (m *Manager) GetRedis() *redis.Client {
	return m.redis
}

// Transaction 执行事务
func (m *Manager) Transaction(fn func(*Manager) error) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		txManager := &Manager{
			db:    tx,
			redis: m.redis,
			
			Task:   m.Task.WithTx(tx).(TaskRepository),
			Result: m.Result.WithTx(tx).(ResultRepository),
			Device: m.Device.WithTx(tx).(DeviceRepository),
			Worker: m.Worker.WithTx(tx).(WorkerRepository),
		}
		return fn(txManager)
	})
}
```

### 7.7 创建单元测试 (`internal/repository/repository_test.go`)
```go
package repository

import (
	"context"
	"testing"
	"time"

	"detect-executor-go/internal/model"
	"detect-executor-go/pkg/database"
	
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// RepositoryTestSuite 仓储测试套件
type RepositoryTestSuite struct {
	suite.Suite
	manager *Manager
	ctx     context.Context
}

// SetupSuite 设置测试套件
func (suite *RepositoryTestSuite) SetupSuite() {
	// 初始化测试数据库
	dbManager, err := database.NewManager(&database.Config{
		MySQL: database.MySQLConfig{
			Host:     "localhost",
			Port:     3306,
			Username: "root",
			Password: "root123",
			Database: "test",
		},
		Redis: database.RedisConfig{
			Host: "localhost",
			Port: 6379,
			DB:   1, // 使用测试数据库
		},
	})
	suite.Require().NoError(err)
	
	// 自动迁移
	err = dbManager.GetMySQL().AutoMigrate(
		&model.DetectTask{},
		&model.DetectResult{},
		&model.Device{},
		&model.Worker{},
	)
	suite.Require().NoError(err)
	
	suite.manager = NewManager(dbManager.GetMySQL(), dbManager.GetRedis())
	suite.ctx = context.Background()
}

// TearDownSuite 清理测试套件
func (suite *RepositoryTestSuite) TearDownSuite() {
	// 清理测试数据
	db := suite.manager.GetDB()
	db.Exec("DELETE FROM detect_tasks")
	db.Exec("DELETE FROM detect_results")
	db.Exec("DELETE FROM devices")
	db.Exec("DELETE FROM workers")
}

// TestTaskRepository 测试任务仓储
func (suite *RepositoryTestSuite) TestTaskRepository() {
	task := &model.DetectTask{
		DeviceID:  "device-001",
		Type:      model.DetectTypeQuality,
		VideoURL:  "rtsp://example.com/stream",
		Priority:  model.PriorityNormal,
		Duration:  300,
		Status:    model.TaskStatusPending,
	}
	
	// 测试创建
	err := suite.manager.Task.Create(suite.ctx, task)
	suite.NoError(err)
	suite.NotEmpty(task.TaskID)
	
	// 测试查询
	found, err := suite.manager.Task.GetByTaskID(suite.ctx, task.TaskID)
	suite.NoError(err)
	suite.NotNil(found)
	suite.Equal(task.DeviceID, found.DeviceID)
	
	// 测试更新状态
	err = suite.manager.Task.UpdateStatus(suite.ctx, task.TaskID, model.TaskStatusRunning)
	suite.NoError(err)
	
	// 验证状态更新
	updated, err := suite.manager.Task.GetByTaskID(suite.ctx, task.TaskID)
	suite.NoError(err)
	suite.Equal(model.TaskStatusRunning, updated.Status)
	suite.NotNil(updated.StartedAt)
	
	// 测试统计
	count, err := suite.manager.Task.CountByStatus(suite.ctx, model.TaskStatusRunning)
	suite.NoError(err)
	suite.Equal(int64(1), count)
}

// TestResultRepository 测试结果仓储
func (suite *RepositoryTestSuite) TestResultRepository() {
	result := &model.DetectResult{
		TaskID:     "task-001",
		DeviceID:   "device-001",
		Type:       model.DetectTypeQuality,
		Score:      0.85,
		Passed:     true,
		Confidence: 0.92,
	}
	
	// 测试创建
	err := suite.manager.Result.Create(suite.ctx, result)
	suite.NoError(err)
	suite.NotEmpty(result.ResultID)
	
	// 测试查询
	found, err := suite.manager.Result.GetByResultID(suite.ctx, result.ResultID)
	suite.NoError(err)
	suite.NotNil(found)
	suite.Equal(result.TaskID, found.TaskID)
}

// TestDeviceRepository 测试设备仓储
func (suite *RepositoryTestSuite) TestDeviceRepository() {
	device := &model.Device{
		DeviceID: "device-001",
		Name:     "Test Camera",
		Type:     "camera",
		IP:       "192.168.1.100",
		Port:     554,
		Status:   "active",
	}
	
	// 测试创建
	err := suite.manager.Device.Create(suite.ctx, device)
	suite.NoError(err)
	
	// 测试查询
	found, err := suite.manager.Device.GetByDeviceID(suite.ctx, device.DeviceID)
	suite.NoError(err)
	suite.NotNil(found)
	suite.Equal(device.Name, found.Name)
	
	// 测试更新最后在线时间
	err = suite.manager.Device.UpdateLastSeen(suite.ctx, device.DeviceID)
	suite.NoError(err)
}

// TestWorkerRepository 测试工作节点仓储
func (suite *RepositoryTestSuite) TestWorkerRepository() {
	worker := &model.Worker{
		WorkerID: "worker-001",
		Name:     "Test Worker",
		IP:       "192.168.1.200",
		Port:     8080,
		Status:   "idle",
		MaxTasks: 5,
	}
	
	// 测试创建
	err := suite.manager.Worker.Create(suite.ctx, worker)
	suite.NoError(err)
	
	// 测试查询
	found, err := suite.manager.Worker.GetByWorkerID(suite.ctx, worker.WorkerID)
	suite.NoError(err)
	suite.NotNil(found)
	suite.Equal(worker.Name, found.Name)
	
	// 测试获取可用工作节点
	workers, err := suite.manager.Worker.GetAvailableWorkers(suite.ctx)
	suite.NoError(err)
	suite.Len(workers, 1)
}

// TestTransaction 测试事务
func (suite *RepositoryTestSuite) TestTransaction() {
	task := &model.DetectTask{
		DeviceID: "device-002",
		Type:     model.DetectTypeQuality,
		VideoURL: "rtsp://example.com/stream2",
		Priority: model.PriorityHigh,
		Duration: 600,
		Status:   model.TaskStatusPending,
	}
	
	result := &model.DetectResult{
		DeviceID:   "device-002",
		Type:       model.DetectTypeQuality,
		Score:      0.75,
		Passed:     true,
		Confidence: 0.88,
	}
	
	// 测试事务
	err := suite.manager.Transaction(func(txManager *Manager) error {
		// 创建任务
		if err := txManager.Task.Create(suite.ctx, task); err != nil {
			return err
		}
		
		// 设置结果的任务ID
		result.TaskID = task.TaskID
		
		// 创建结果
		if err := txManager.Result.Create(suite.ctx, result); err != nil {
			return err
		}
		
		return nil
	})
	
	suite.NoError(err)
	suite.NotEmpty(task.TaskID)
	suite.NotEmpty(result.ResultID)
	suite.Equal(task.TaskID, result.TaskID)
}

// TestRepositoryTestSuite 运行测试套件
func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
```

## 验证步骤

### 7.8 运行测试验证
```bash
# 启动测试数据库
docker run -d --name mysql-test \
  -e MYSQL_ROOT_PASSWORD=root123 \
  -e MYSQL_DATABASE=test \
  -p 3306:3306 \
  mysql:8.0

docker run -d --name redis-test \
  -p 6379:6379 \
  redis:7-alpine

# 等待数据库启动
sleep 10

# 运行仓储层测试
go test ./internal/repository/ -v

# 检查依赖
go mod tidy
```

## 预期结果
- 完整的数据访问层实现，支持所有核心业务操作
- Redis缓存集成，提高查询性能
- 事务支持，保证数据一致性
- 完善的错误处理和日志记录
- 全面的单元测试覆盖
- 查询优化和索引友好的设计

## 状态
- [ ] 步骤 7 完成：数据访问层实现

---

## 下一步预告
步骤 8: 外部客户端集成 (`pkg/client/`)
- HTTP客户端封装
- 检测引擎客户端
- 文件存储客户端
- 消息队列客户端
