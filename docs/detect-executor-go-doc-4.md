# 步骤 4: 搭建HTTP服务器

## 目标
创建HTTP服务器主程序，集成配置管理和日志系统，实现基础路由和优雅关闭机制。

## 操作步骤

### 4.1 创建服务器主程序 (`cmd/server/main.go`)
```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"detect-executor-go/pkg/config"
	"detect-executor-go/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	configLoader := config.NewLoader()
	cfg, err := configLoader.Load("")
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	loggerConfig := &logger.Config{
		Level:      cfg.Log.Level,
		Format:     cfg.Log.Format,
		Output:     cfg.Log.Output,
		Filename:   cfg.Log.Filename,
		MaxSize:    cfg.Log.MaxSize,
		MaxAge:     cfg.Log.MaxAge,
		MaxBackups: cfg.Log.MaxBackups,
		Compress:   cfg.Log.Compress,
	}

	log, err := logger.NewLogger(loggerConfig)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Close()

	log.Info("Starting detect-executor-go server")

	// 设置Gin模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 创建Gin引擎
	router := gin.New()

	// 添加中间件
	router.Use(logger.RequestID())
	router.Use(logger.GinLogger(log))
	router.Use(logger.GinRecovery(log))

	// 设置路由
	setupRoutes(router, log)

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:         cfg.Server.GetAddr(),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	// 启动服务器
	go func() {
		log.WithField("addr", srv.Addr).Info("HTTP server starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithField("error", err).Fatal("Failed to start server")
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// 优雅关闭服务器，等待现有连接完成
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.WithField("error", err).Error("Server forced to shutdown")
	} else {
		log.Info("Server exited")
	}
}

// setupRoutes 设置路由
func setupRoutes(router *gin.Engine, log *logger.Logger) {
	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
			"service":   "detect-executor-go",
		})
	})

	// API版本路由组
	v1 := router.Group("/api/v1")
	{
		// 基础信息
		v1.GET("/info", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"service":     "detect-executor-go",
				"version":     "1.0.0",
				"description": "Go-based detection execution service",
				"timestamp":   time.Now().Unix(),
			})
		})

		// 检测相关路由（后续实现）
		detect := v1.Group("/detect")
		{
			// 提交检测任务
			detect.POST("/task", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Detection task endpoint - to be implemented",
				})
			})

			// 查询任务状态
			detect.GET("/task/:id", func(c *gin.Context) {
				taskID := c.Param("id")
				c.JSON(http.StatusOK, gin.H{
					"task_id": taskID,
					"message": "Task status endpoint - to be implemented",
				})
			})

			// 获取检测结果
			detect.GET("/result/:id", func(c *gin.Context) {
				taskID := c.Param("id")
				c.JSON(http.StatusOK, gin.H{
					"task_id": taskID,
					"message": "Detection result endpoint - to be implemented",
				})
			})
		}
	}

	// 404处理
	router.NoRoute(func(c *gin.Context) {
		log.WithField("path", c.Request.URL.Path).Warn("Route not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "The requested resource was not found",
			"path":    c.Request.URL.Path,
		})
	})
}
```

### 4.2 创建服务器启动脚本 (`scripts/start.sh`)
```bash
#!/bin/bash

# 设置环境变量
export GO_ENV=${GO_ENV:-dev}
export DETECT_SERVER_PORT=${DETECT_SERVER_PORT:-8080}

# 创建必要的目录
mkdir -p logs
mkdir -p bin

# 编译并运行
echo "Building detect-executor-go..."
go build -o bin/server cmd/server/main.go

if [ $? -eq 0 ]; then
    echo "Starting server..."
    ./bin/server
else
    echo "Build failed!"
    exit 1
fi
```

### 4.3 创建开发环境启动脚本 (`scripts/dev.sh`)
```bash
#!/bin/bash

# 设置开发环境变量
export GO_ENV=dev
export DETECT_SERVER_PORT=8080
export DETECT_LOG_LEVEL=debug

# 创建必要的目录
mkdir -p logs

# 使用 go run 直接运行（开发模式）
echo "Starting development server..."
go run cmd/server/main.go
```

### 4.4 创建生产环境启动脚本 (`scripts/prod.sh`)
```bash
#!/bin/bash

# 设置生产环境变量
export GO_ENV=prod
export DETECT_SERVER_PORT=${DETECT_SERVER_PORT:-8080}
export DETECT_LOG_LEVEL=${DETECT_LOG_LEVEL:-info}

# 创建必要的目录
mkdir -p logs
mkdir -p bin

# 编译优化版本
echo "Building production server..."
go build -ldflags="-w -s" -o bin/server cmd/server/main.go

if [ $? -eq 0 ]; then
    echo "Starting production server..."
    ./bin/server
else
    echo "Build failed!"
    exit 1
fi
```

### 4.5 创建Docker支持文件 (`Dockerfile`)
```dockerfile
# 构建阶段
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的包
RUN apk add --no-cache git

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译应用
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server cmd/server/main.go

# 运行阶段
FROM alpine:latest

# 安装ca证书
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非root用户
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/server .
COPY --from=builder /app/configs ./configs

# 创建日志目录
RUN mkdir -p logs && chown -R appuser:appgroup /app

# 切换到非root用户
USER appuser

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 启动应用
CMD ["./server"]
```

### 4.6 创建Docker Compose文件 (`docker-compose.yml`)
```yaml
version: '3.8'

services:
  detect-executor:
    build: .
    ports:
      - "8080:8080"
    environment:
      - GO_ENV=prod
      - DETECT_SERVER_PORT=8080
      - DETECT_LOG_LEVEL=info
    volumes:
      - ./logs:/app/logs
      - ./configs:/app/configs
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  # Redis服务（后续步骤会用到）
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    restart: unless-stopped
    command: redis-server --appendonly yes

  # MySQL服务（后续步骤会用到）
  mysql:
    image: mysql:8.0
    ports:
      - "3306:3306"
    environment:
      - MYSQL_ROOT_PASSWORD=root123
      - MYSQL_DATABASE=detect_executor
      - MYSQL_USER=detect
      - MYSQL_PASSWORD=detect123
    volumes:
      - mysql_data:/var/lib/mysql
      - ./scripts/init.sql:/docker-entrypoint-initdb.d/init.sql
    restart: unless-stopped

volumes:
  redis_data:
  mysql_data:
```

## 验证步骤

### 4.7 创建服务器测试文件 (`cmd/server/main_test.go`)
```go
package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"detect-executor-go/pkg/config"
	"detect-executor-go/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSetupRoutes(t *testing.T) {
	// 设置测试模式
	gin.SetMode(gin.TestMode)

	// 创建测试日志
	loggerConfig := &logger.Config{
		Level:  "debug",
		Format: "json",
		Output: "stdout",
	}
	log, err := logger.NewLogger(loggerConfig)
	assert.NoError(t, err)

	// 创建路由
	router := gin.New()
	setupRoutes(router, log)

	// 测试健康检查端点
	t.Run("Health Check", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/health", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "ok")
	})

	// 测试API信息端点
	t.Run("API Info", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/info", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "detect-executor-go")
	})

	// 测试检测任务端点
	t.Run("Detection Task", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/detect/task", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "to be implemented")
	})

	// 测试404处理
	t.Run("Not Found", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/nonexistent", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Not Found")
	})
}

func TestMain(m *testing.M) {
	// 添加测试依赖
	// go get github.com/stretchr/testify/assert
	m.Run()
}
```

### 4.8 验证HTTP服务器
```bash
# 添加测试依赖
go get github.com/stretchr/testify/assert

# 运行服务器测试
go test ./cmd/server/

# 创建必要的目录
mkdir -p bin
mkdir -p scripts
mkdir -p logs

# 给脚本添加执行权限
chmod +x scripts/*.sh

# 编译服务器
go build -o bin/server cmd/server/main.go

# 运行服务器（开发模式）
GO_ENV=dev ./bin/server

# 或者使用脚本运行
./scripts/dev.sh
```

### 4.9 测试API端点
```bash
# 在另一个终端测试API
# 健康检查
curl http://localhost:8080/health

# API信息
curl http://localhost:8080/api/v1/info

# 检测任务（POST）
curl -X POST http://localhost:8080/api/v1/detect/task

# 任务状态查询
curl http://localhost:8080/api/v1/detect/task/123

# 检测结果查询
curl http://localhost:8080/api/v1/detect/result/123

# 测试404
curl http://localhost:8080/nonexistent
```

## 预期结果
- HTTP服务器成功启动，监听配置的端口
- 所有基础API端点正常响应
- 日志系统正常工作，记录请求和错误信息
- 优雅关闭机制正常工作
- Docker支持文件创建完成
- 包含完整的单元测试
- 启动脚本创建完成并可执行

## 状态
- [ ] 步骤 4 完成：HTTP服务器搭建

---

## 下一步预告
步骤 5: 数据库连接设置 (`pkg/database/`)
- MySQL连接池配置
- Redis连接配置
- 数据库健康检查
- 连接重试机制
