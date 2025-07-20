# 步骤 5: 数据库连接设置

## 目标
实现MySQL和Redis数据库连接管理，包括连接池配置、健康检查和连接重试机制。

## 操作步骤

### 5.1 添加数据库相关依赖
```bash
# 添加MySQL驱动和GORM
go get gorm.io/gorm
go get gorm.io/driver/mysql

# 添加Redis客户端
go get github.com/redis/go-redis/v9

# 添加数据库连接池
go get github.com/jmoiron/sqlx
```

### 5.2 创建MySQL连接管理器 (`pkg/database/mysql.go`)
```go
package database

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MySQLConfig MySQL配置
type MySQLConfig struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	Database        string `yaml:"database"`
	Charset         string `yaml:"charset"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime int    `yaml:"conn_max_idle_time"`
}

// MySQLClient MySQL客户端
type MySQLClient struct {
	DB     *gorm.DB
	config *MySQLConfig
}

// NewMySQLClient 创建MySQL客户端
func NewMySQLClient(config *MySQLConfig) (*MySQLClient, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.Charset,
	)

	// GORM配置
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	// 获取底层sql.DB以配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// 配置连接池
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(config.ConnMaxIdleTime) * time.Second)

	return &MySQLClient{
		DB:     db,
		config: config,
	}, nil
}

// Ping 检查MySQL连接
func (c *MySQLClient) Ping() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// Close 关闭MySQL连接
func (c *MySQLClient) Close() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// GetStats 获取连接池统计信息
func (c *MySQLClient) GetStats() map[string]interface{} {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	stats := sqlDB.Stats()
	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":              stats.InUse,
		"idle":                stats.Idle,
	}
}
```

### 5.3 创建Redis连接管理器 (`pkg/database/redis.go`)
```go
package database

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisConfig Redis配置
type RedisConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Password     string `yaml:"password"`
	Database     int    `yaml:"database"`
	PoolSize     int    `yaml:"pool_size"`
	MinIdleConns int    `yaml:"min_idle_conns"`
	MaxRetries   int    `yaml:"max_retries"`
	DialTimeout  int    `yaml:"dial_timeout"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
	IdleTimeout  int    `yaml:"idle_timeout"`
}

// RedisClient Redis客户端
type RedisClient struct {
	Client *redis.Client
	config *RedisConfig
}

// NewRedisClient 创建Redis客户端
func NewRedisClient(config *RedisConfig) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password:     config.Password,
		DB:           config.Database,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		MaxRetries:   config.MaxRetries,
		DialTimeout:  time.Duration(config.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(config.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(config.IdleTimeout) * time.Second,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisClient{
		Client: rdb,
		config: config,
	}, nil
}

// Ping 检查Redis连接
func (c *RedisClient) Ping(ctx context.Context) error {
	return c.Client.Ping(ctx).Err()
}

// Close 关闭Redis连接
func (c *RedisClient) Close() error {
	return c.Client.Close()
}

// GetStats 获取连接池统计信息
func (c *RedisClient) GetStats() map[string]interface{} {
	stats := c.Client.PoolStats()
	return map[string]interface{}{
		"hits":         stats.Hits,
		"misses":       stats.Misses,
		"timeouts":     stats.Timeouts,
		"total_conns":  stats.TotalConns,
		"idle_conns":   stats.IdleConns,
		"stale_conns":  stats.StaleConns,
	}
}
```

### 5.4 创建数据库管理器 (`pkg/database/manager.go`)
```go
package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"detect-executor-go/pkg/logger"
)

// Manager 数据库管理器
type Manager struct {
	MySQL  *MySQLClient
	Redis  *RedisClient
	logger *logger.Logger
	mu     sync.RWMutex
}

// Config 数据库配置
type Config struct {
	MySQL *MySQLConfig `yaml:"mysql"`
	Redis *RedisConfig `yaml:"redis"`
}

// NewManager 创建数据库管理器
func NewManager(config *Config, log *logger.Logger) (*Manager, error) {
	manager := &Manager{
		logger: log,
	}

	// 初始化MySQL
	if config.MySQL != nil {
		mysqlClient, err := NewMySQLClient(config.MySQL)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize MySQL: %w", err)
		}
		manager.MySQL = mysqlClient
		log.Info("MySQL connection established")
	}

	// 初始化Redis
	if config.Redis != nil {
		redisClient, err := NewRedisClient(config.Redis)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Redis: %w", err)
		}
		manager.Redis = redisClient
		log.Info("Redis connection established")
	}

	return manager, nil
}

// HealthCheck 健康检查
func (m *Manager) HealthCheck(ctx context.Context) map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]interface{})

	// MySQL健康检查
	if m.MySQL != nil {
		if err := m.MySQL.Ping(); err != nil {
			result["mysql"] = map[string]interface{}{
				"status": "unhealthy",
				"error":  err.Error(),
			}
		} else {
			result["mysql"] = map[string]interface{}{
				"status": "healthy",
				"stats":  m.MySQL.GetStats(),
			}
		}
	}

	// Redis健康检查
	if m.Redis != nil {
		if err := m.Redis.Ping(ctx); err != nil {
			result["redis"] = map[string]interface{}{
				"status": "unhealthy",
				"error":  err.Error(),
			}
		} else {
			result["redis"] = map[string]interface{}{
				"status": "healthy",
				"stats":  m.Redis.GetStats(),
			}
		}
	}

	return result
}

// Close 关闭所有数据库连接
func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errors []error

	if m.MySQL != nil {
		if err := m.MySQL.Close(); err != nil {
			errors = append(errors, fmt.Errorf("MySQL close error: %w", err))
		}
	}

	if m.Redis != nil {
		if err := m.Redis.Close(); err != nil {
			errors = append(errors, fmt.Errorf("Redis close error: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("database close errors: %v", errors)
	}

	m.logger.Info("All database connections closed")
	return nil
}

// StartHealthCheckRoutine 启动健康检查例程
func (m *Manager) StartHealthCheckRoutine(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			health := m.HealthCheck(ctx)
			m.logger.WithField("health", health).Debug("Database health check")
		}
	}
}
```

### 5.5 更新配置结构体 (`pkg/config/config.go`)
在现有配置中添加数据库配置：
```go
// 在Config结构体中添加
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	Detect   DetectConfig   `yaml:"detect"`
	Log      LogConfig      `yaml:"log"`
	Metrics  MetricsConfig  `yaml:"metrics"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	Database        string `yaml:"database"`
	Charset         string `yaml:"charset"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime int    `yaml:"conn_max_idle_time"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Password     string `yaml:"password"`
	Database     int    `yaml:"database"`
	PoolSize     int    `yaml:"pool_size"`
	MinIdleConns int    `yaml:"min_idle_conns"`
	MaxRetries   int    `yaml:"max_retries"`
	DialTimeout  int    `yaml:"dial_timeout"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
	IdleTimeout  int    `yaml:"idle_timeout"`
}
```

### 5.6 更新配置文件模板
更新 `configs/config.dev.yaml`:
```yaml
# 数据库配置
database:
  host: ${DETECT_DB_HOST:localhost}
  port: ${DETECT_DB_PORT:3306}
  username: ${DETECT_DB_USER:detect}
  password: ${DETECT_DB_PASSWORD:detect123}
  database: ${DETECT_DB_NAME:detect_executor}
  charset: utf8mb4
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: 3600
  conn_max_idle_time: 600

# Redis配置
redis:
  host: ${DETECT_REDIS_HOST:localhost}
  port: ${DETECT_REDIS_PORT:6379}
  password: ${DETECT_REDIS_PASSWORD:}
  database: ${DETECT_REDIS_DB:0}
  pool_size: 100
  min_idle_conns: 10
  max_retries: 3
  dial_timeout: 5
  read_timeout: 3
  write_timeout: 3
  idle_timeout: 300
```

## 验证步骤

### 5.7 创建数据库测试文件 (`pkg/database/database_test.go`)
```go
package database

import (
	"context"
	"testing"
	"time"

	"detect-executor-go/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestMySQLConnection(t *testing.T) {
	config := &MySQLConfig{
		Host:            "localhost",
		Port:            3306,
		Username:        "root",
		Password:        "root123",
		Database:        "test",
		Charset:         "utf8mb4",
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: 3600,
		ConnMaxIdleTime: 600,
	}

	client, err := NewMySQLClient(config)
	if err != nil {
		t.Skip("MySQL not available, skipping test")
		return
	}
	defer client.Close()

	// 测试连接
	err = client.Ping()
	assert.NoError(t, err)

	// 测试统计信息
	stats := client.GetStats()
	assert.NotEmpty(t, stats)
}

func TestRedisConnection(t *testing.T) {
	config := &RedisConfig{
		Host:         "localhost",
		Port:         6379,
		Password:     "",
		Database:     0,
		PoolSize:     10,
		MinIdleConns: 5,
		MaxRetries:   3,
		DialTimeout:  5,
		ReadTimeout:  3,
		WriteTimeout: 3,
		IdleTimeout:  300,
	}

	client, err := NewRedisClient(config)
	if err != nil {
		t.Skip("Redis not available, skipping test")
		return
	}
	defer client.Close()

	// 测试连接
	ctx := context.Background()
	err = client.Ping(ctx)
	assert.NoError(t, err)

	// 测试统计信息
	stats := client.GetStats()
	assert.NotEmpty(t, stats)
}

func TestDatabaseManager(t *testing.T) {
	loggerConfig := &logger.Config{
		Level:  "debug",
		Format: "json",
		Output: "stdout",
	}
	log, err := logger.NewLogger(loggerConfig)
	assert.NoError(t, err)

	config := &Config{
		MySQL: &MySQLConfig{
			Host:            "localhost",
			Port:            3306,
			Username:        "root",
			Password:        "root123",
			Database:        "test",
			Charset:         "utf8mb4",
			MaxOpenConns:    10,
			MaxIdleConns:    5,
			ConnMaxLifetime: 3600,
			ConnMaxIdleTime: 600,
		},
		Redis: &RedisConfig{
			Host:         "localhost",
			Port:         6379,
			Password:     "",
			Database:     0,
			PoolSize:     10,
			MinIdleConns: 5,
			MaxRetries:   3,
			DialTimeout:  5,
			ReadTimeout:  3,
			WriteTimeout: 3,
			IdleTimeout:  300,
		},
	}

	manager, err := NewManager(config, log)
	if err != nil {
		t.Skip("Database not available, skipping test")
		return
	}
	defer manager.Close()

	// 测试健康检查
	ctx := context.Background()
	health := manager.HealthCheck(ctx)
	assert.NotEmpty(t, health)
}
```

### 5.8 验证数据库连接
```bash
# 运行数据库测试
go test ./pkg/database/

# 检查依赖
go mod tidy

# 启动测试数据库（使用Docker）
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

# 再次运行测试
go test ./pkg/database/
```

## 预期结果
- MySQL连接池配置完成，支持连接复用和超时管理
- Redis连接配置完成，支持连接池和重试机制
- 数据库管理器创建完成，统一管理所有数据库连接
- 健康检查机制实现，可监控数据库连接状态
- 包含完整的单元测试
- 配置文件更新，支持环境变量覆盖

## 状态
- [ ] 步骤 5 完成：数据库连接设置

---

## 下一步预告
步骤 6: 数据模型定义 (`internal/model/`)
- 检测任务模型
- 检测结果模型
- 数据验证规则
- 模型关系定义
