package database

import (
	"context"
	"detect-executor-go/pkg/logger"
	"fmt"
	"sync"
	"time"
)

// sync 是 Go 语言中的一个包，提供了同步原语，包括互斥锁（Mutex）和读写锁（RWMutex）
// RWMutex 用于保护 MySQL 和 Redis 的访问
// 在多线程环境下，读写锁可以提高并发性能
// RLock() 用于读取操作
// RUnlock() 用于释放读取操作
// Lock() 用于写入操作
// Unlock() 用于释放写入操作
type Manager struct {
	MySQL  *MySQLClient
	Redis  *RedisClient
	logger *logger.Logger
	//mu 是读写锁，用于保护 MySQL 和 Redis 的访问
	//在多线程环境下，读写锁可以提高并发性能
	//mu.RLock() 用于读取操作
	//mu.RUnlock() 用于释放读取操作
	//mu.Lock() 用于写入操作
	//mu.Unlock() 用于释放写入操作
	mu *sync.RWMutex
}

type Config struct {
	MySQL *MySQLConfig `yaml:"mysql"`
	Redis *RedisConfig `yaml:"redis"`
}

func NewManager(config *Config, log *logger.Logger) (*Manager, error) {
	manager := &Manager{
		logger: log,
	}

	// init mysql
	if config.MySQL != nil {
		mysqlClient, err := NewMySQLClient(config.MySQL)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize MySQL: %w ", err)
		}
		manager.MySQL = mysqlClient
		log.Info("MySQL connection established")
	}

	// initialize Redis
	if config.Redis != nil {
		redisClient, err := NewRedisClient(config.Redis)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Redis: %w ", err)
		}
		manager.Redis = redisClient
		log.Info("Redis connection established")
	}

	return manager, nil

}

func (m *Manager) HealthCheck(ctx context.Context) map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]interface{})

	// mysql health check
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

	// redis health check
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

func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errors []error

	if m.MySQL != nil {
		if err := m.MySQL.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close MySQL: %w", err))
		}
	}

	if m.Redis != nil {
		if err := m.Redis.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close Redis: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to close database: %v", errors)
	}

	m.logger.Info("All database connections closed")

	return nil
}

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
