package config

import (
	"os"
	"testing"
	"time"
)

func TestConfigLoad(t *testing.T) {
	configContent := `
	server:
	  host: "0.0.0.0"
	  port: 8080
	  read_timeout: "30s"
	  
	database:
	  host: "localhost"
	  port: 3306
	  username: "root"
	  password: "password"
	  database: "detect_db_test"
	
	detect:
	  max_workers: 10
	  batch_size: 10
`

	tmpFile, err := os.CreateTemp("", "config_test_*.yaml")
	if err != nil {
		t.Fatalf("create temp file failed: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatalf("write config content failed: %v", err)
	}
	tmpFile.Close()

	loader := NewLoader()
	config, err := loader.Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("load config failed: %v", err)
	}

	if config.Server.Host != "0.0.0.0" {
		t.Errorf("expected server host to be 0.0.0.0, got %s", config.Server.Host)
	}
	if config.Server.Port != 8080 {
		t.Errorf("expected server port to be 8080, got %d", config.Server.Port)
	}
	if config.Server.ReadTimeout != 30*time.Second {
		t.Errorf("expected server read timeout to be 30s, got %v", config.Server.ReadTimeout)
	}
	if config.Server.WriteTimeout != 30*time.Second {
		t.Errorf("expected server write timeout to be 30s, got %v", config.Server.WriteTimeout)
	}
	if config.Server.IdleTimeout != 5*time.Minute {
		t.Errorf("expected server idle timeout to be 5m, got %v", config.Server.IdleTimeout)
	}

}

func TestConfigDefault(t *testing.T) {
	loader := NewLoader()
	config, err := loader.Load("nonexist.yaml")
	if err != nil {
		t.Fatalf("load config failed: %v", err)
	}

	if config.Server.Host != "0.0.0.0" {
		t.Errorf("expected default server host to be 0.0.0.0, got %s", config.Server.Host)
	}
	if config.Server.Port != 8080 {
		t.Errorf("expected default server port to be 8080, got %d", config.Server.Port)
	}
	if config.Server.ReadTimeout != 30*time.Second {
		t.Errorf("expected default server read timeout to be 30s, got %v", config.Server.ReadTimeout)
	}

}

func TestConfigMethods(t *testing.T) {
	config := &Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
		},
		Database: DatabaseConfig{
			Host:     "dbhost",
			Port:     3306,
			Username: "root",
			Password: "password",
			Database: "detect_db_test",
		},
		Detect: DetectConfig{
			MaxWorkers:     10,
			BatchSize:      10,
			QueueSize:      100,
			Timeout:        10 * time.Second,
			RetryCount:     3,
			RetryDelay:     1 * time.Second,
			CppEngineURL:   "http://localhost:8081",
			VideoStreamURL: "http://localhost:8082",
		},
		Log: LogConfig{
			Level:      "info",
			Format:     "json",
			Output:     "stdout",
			Filename:   "logs/detect-executor-go.log",
			MaxSize:    100,
			MaxAge:     28,
			MaxBackups: 3,
			Compress:   true,
		},
		Metrics: MetricsConfig{
			Enabled: false,
			Path:    "/metrics",
			Port:    9090,
		},
		Redis: RedisConfig{
			Host:         "redishost",
			Port:         6379,
			Password:     "password",
			DB:           0,
			PoolSize:     10,
			MinIdleConns: 10,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		},
	}

	if config.Server.GetAddr() != "0.0.0.0:8080" {
		t.Errorf("expected server addr to be 0.0.0.0:8080, got %s", config.Server.GetAddr())
	}
	if config.Database.GetDSN() != "root:password@tcp(dbhost:3306)/detect_db_test?charset=utf8mb4&parseTime=True&loc=Local" {
		t.Errorf("expected database dsn to be root:password@tcp(dbhost:3306)/detect_db_test?charset=utf8mb4&parseTime=True&loc=Local, got %s", config.Database.GetDSN())
	}

	if config.Metrics.GetAddr() != ":9090" {
		t.Errorf("expected metrics addr to be :9090, got %s", config.Metrics.GetAddr())
	}
	if config.Redis.GetAddr() != "redishost:6379" {
		t.Errorf("expected redis addr to be redishost:6379, got %s", config.Redis.GetAddr())
	}
}
