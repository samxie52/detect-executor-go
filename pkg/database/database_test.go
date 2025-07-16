package database

import (
	"context"
	"detect-executor-go/pkg/logger"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMySQLConnection(t *testing.T) {
	config := &MySQLConfig{
		Host:            "localhost",
		Port:            3306,
		Username:        "root",
		Password:        "password",
		Database:        "test",
		Charset:         "utf8mb4",
		MaxOpenConns:    10,
		MaxIdleConns:    10,
		ConnMaxIdleTime: 600,
		ConnMaxLifetime: 3600,
	}

	client, err := NewMySQLClient(config)
	if err != nil {
		t.Fatalf("failed to connect to MySQL: %v", err)
	}
	defer client.Close()

	err = client.Ping()
	assert.NoError(t, err)

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
		DialTimeout:  5,
		ReadTimeout:  3,
		WriteTimeout: 3,
	}

	client, err := NewRedisClient(config)
	if err != nil {
		t.Skip("Redis not available, skipping test")
		return
	}
	defer client.Close()

	ctx := context.Background()
	err = client.Ping(ctx)
	assert.NoError(t, err)

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
			Password:        "password",
			Database:        "test",
			Charset:         "utf8mb4",
			MaxOpenConns:    10,
			MaxIdleConns:    10,
			ConnMaxIdleTime: 600,
			ConnMaxLifetime: 3600,
		},
		Redis: &RedisConfig{
			Host:         "localhost",
			Port:         6379,
			Password:     "",
			Database:     0,
			PoolSize:     10,
			MinIdleConns: 5,
			DialTimeout:  5,
			ReadTimeout:  3,
			WriteTimeout: 3,
		},
	}

	manager, err := NewManager(config, log)
	if err != nil {
		t.Skip("Database not available, skipping test")
		return
	}
	defer manager.Close()

	ctx := context.Background()
	health := manager.HealthCheck(ctx)
	assert.NotEmpty(t, health)

}
