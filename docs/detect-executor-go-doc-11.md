# 步骤 11: 监控指标实现

## 目标
实现detect-executor-go的监控指标系统，集成Prometheus指标收集，提供健康检查、性能监控和业务指标统计，支持告警和可视化监控面板。

## 操作步骤

### 11.1 添加监控依赖
```bash
# 添加Prometheus客户端
go get github.com/prometheus/client_golang

# 添加健康检查库
go get github.com/heptiolabs/healthcheck
```

### 11.2 创建基础指标 (`pkg/metrics/base.go`)
```go
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics 指标收集器
type Metrics struct {
	// HTTP指标
	HTTPRequestsTotal     *prometheus.CounterVec
	HTTPRequestDuration   *prometheus.HistogramVec
	HTTPRequestsInFlight  prometheus.Gauge

	// 业务指标
	TasksTotal            *prometheus.CounterVec
	TasksInProgress       prometheus.Gauge
	TaskDuration          *prometheus.HistogramVec
	DetectionTotal        *prometheus.CounterVec
	DetectionDuration     *prometheus.HistogramVec
	ResultsTotal          *prometheus.CounterVec

	// 系统指标
	DatabaseConnections   prometheus.Gauge
	RedisConnections      prometheus.Gauge
	ExternalAPIRequests   *prometheus.CounterVec
	ExternalAPILatency    *prometheus.HistogramVec

	// 错误指标
	ErrorsTotal           *prometheus.CounterVec
}

// NewMetrics 创建指标收集器
func NewMetrics() *Metrics {
	return &Metrics{
		// HTTP指标
		HTTPRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		HTTPRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
		HTTPRequestsInFlight: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "http_requests_in_flight",
				Help: "Number of HTTP requests currently being processed",
			},
		),

		// 业务指标
		TasksTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "tasks_total",
				Help: "Total number of detection tasks",
			},
			[]string{"device_id", "type", "status"},
		),
		TasksInProgress: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "tasks_in_progress",
				Help: "Number of tasks currently in progress",
			},
		),
		TaskDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "task_duration_seconds",
				Help:    "Task execution duration in seconds",
				Buckets: []float64{1, 5, 10, 30, 60, 300, 600, 1800, 3600},
			},
			[]string{"device_id", "type", "status"},
		),
		DetectionTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "detections_total",
				Help: "Total number of detections performed",
			},
			[]string{"type", "result"},
		),
		DetectionDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "detection_duration_seconds",
				Help:    "Detection execution duration in seconds",
				Buckets: []float64{1, 5, 10, 30, 60, 300, 600, 1800, 3600},
			},
			[]string{"type"},
		),
		ResultsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "results_total",
				Help: "Total number of detection results",
			},
			[]string{"device_id", "type", "passed"},
		),

		// 系统指标
		DatabaseConnections: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "database_connections",
				Help: "Number of active database connections",
			},
		),
		RedisConnections: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "redis_connections",
				Help: "Number of active Redis connections",
			},
		),
		ExternalAPIRequests: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "external_api_requests_total",
				Help: "Total number of external API requests",
			},
			[]string{"service", "method", "status"},
		),
		ExternalAPILatency: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "external_api_latency_seconds",
				Help:    "External API request latency in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"service", "method"},
		),

		// 错误指标
		ErrorsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "errors_total",
				Help: "Total number of errors",
			},
			[]string{"component", "type"},
		),
	}
}

// RecordHTTPRequest 记录HTTP请求指标
func (m *Metrics) RecordHTTPRequest(method, path, status string, duration float64) {
	m.HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
	m.HTTPRequestDuration.WithLabelValues(method, path).Observe(duration)
}

// IncHTTPRequestsInFlight 增加进行中的HTTP请求数
func (m *Metrics) IncHTTPRequestsInFlight() {
	m.HTTPRequestsInFlight.Inc()
}

// DecHTTPRequestsInFlight 减少进行中的HTTP请求数
func (m *Metrics) DecHTTPRequestsInFlight() {
	m.HTTPRequestsInFlight.Dec()
}

// RecordTask 记录任务指标
func (m *Metrics) RecordTask(deviceID, taskType, status string, duration float64) {
	m.TasksTotal.WithLabelValues(deviceID, taskType, status).Inc()
	m.TaskDuration.WithLabelValues(deviceID, taskType, status).Observe(duration)
}

// IncTasksInProgress 增加进行中的任务数
func (m *Metrics) IncTasksInProgress() {
	m.TasksInProgress.Inc()
}

// DecTasksInProgress 减少进行中的任务数
func (m *Metrics) DecTasksInProgress() {
	m.TasksInProgress.Dec()
}

// RecordDetection 记录检测指标
func (m *Metrics) RecordDetection(detectType, result string, duration float64) {
	m.DetectionTotal.WithLabelValues(detectType, result).Inc()
	m.DetectionDuration.WithLabelValues(detectType).Observe(duration)
}

// RecordResult 记录结果指标
func (m *Metrics) RecordResult(deviceID, resultType string, passed bool) {
	passedStr := "false"
	if passed {
		passedStr = "true"
	}
	m.ResultsTotal.WithLabelValues(deviceID, resultType, passedStr).Inc()
}

// SetDatabaseConnections 设置数据库连接数
func (m *Metrics) SetDatabaseConnections(count float64) {
	m.DatabaseConnections.Set(count)
}

// SetRedisConnections 设置Redis连接数
func (m *Metrics) SetRedisConnections(count float64) {
	m.RedisConnections.Set(count)
}

// RecordExternalAPIRequest 记录外部API请求指标
func (m *Metrics) RecordExternalAPIRequest(service, method, status string, latency float64) {
	m.ExternalAPIRequests.WithLabelValues(service, method, status).Inc()
	m.ExternalAPILatency.WithLabelValues(service, method).Observe(latency)
}

// RecordError 记录错误指标
func (m *Metrics) RecordError(component, errorType string) {
	m.ErrorsTotal.WithLabelValues(component, errorType).Inc()
}
```

### 11.3 创建HTTP指标中间件 (`pkg/metrics/middleware.go`)
```go
package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// HTTPMetricsMiddleware HTTP指标中间件
func HTTPMetricsMiddleware(metrics *Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// 增加进行中的请求数
		metrics.IncHTTPRequestsInFlight()
		defer metrics.DecHTTPRequestsInFlight()

		// 处理请求
		c.Next()

		// 记录指标
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		
		metrics.RecordHTTPRequest(
			c.Request.Method,
			c.FullPath(),
			status,
			duration,
		)
	}
}
```

### 11.4 创建健康检查 (`pkg/metrics/health.go`)
```go
package metrics

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/heptiolabs/healthcheck"
	"github.com/redis/go-redis/v9"

	"detect-executor-go/pkg/database"
	"detect-executor-go/pkg/logger"
)

// HealthChecker 健康检查器
type HealthChecker struct {
	handler healthcheck.Handler
	logger  *logger.Logger
}

// NewHealthChecker 创建健康检查器
func NewHealthChecker(logger *logger.Logger, dbManager *database.Manager) *HealthChecker {
	health := healthcheck.NewHandler()

	checker := &HealthChecker{
		handler: health,
		logger:  logger,
	}

	// 添加数据库健康检查
	if dbManager.MySQL != nil {
		health.AddReadinessCheck("mysql", checker.mysqlCheck(dbManager.MySQL))
	}

	if dbManager.Redis != nil {
		health.AddReadinessCheck("redis", checker.redisCheck(dbManager.Redis))
	}

	// 添加存活检查
	health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(100))

	return checker
}

// mysqlCheck MySQL健康检查
func (h *HealthChecker) mysqlCheck(db *sql.DB) healthcheck.Check {
	return func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			h.logger.Error("MySQL health check failed", map[string]interface{}{
				"error": err.Error(),
			})
			return fmt.Errorf("mysql ping failed: %w", err)
		}

		return nil
	}
}

// redisCheck Redis健康检查
func (h *HealthChecker) redisCheck(client *redis.Client) healthcheck.Check {
	return func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.Ping(ctx).Err(); err != nil {
			h.logger.Error("Redis health check failed", map[string]interface{}{
				"error": err.Error(),
			})
			return fmt.Errorf("redis ping failed: %w", err)
		}

		return nil
	}
}

// LivenessHandler 存活检查处理器
func (h *HealthChecker) LivenessHandler() gin.HandlerFunc {
	return gin.WrapF(h.handler.LiveEndpoint)
}

// ReadinessHandler 就绪检查处理器
func (h *HealthChecker) ReadinessHandler() gin.HandlerFunc {
	return gin.WrapF(h.handler.ReadyEndpoint)
}

// HealthHandler 健康检查处理器
func (h *HealthChecker) HealthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "detect-executor-go",
			"version": "1.0.0",
			"time":    time.Now().UTC(),
		})
	}
}
```

### 11.5 创建监控管理器 (`pkg/metrics/manager.go`)
```go
package metrics

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"detect-executor-go/pkg/database"
	"detect-executor-go/pkg/logger"
)

// Manager 监控管理器
type Manager struct {
	metrics      *Metrics
	healthCheck  *HealthChecker
	logger       *logger.Logger
	dbManager    *database.Manager
	updateTicker *time.Ticker
	stopChan     chan struct{}
}

// NewManager 创建监控管理器
func NewManager(logger *logger.Logger, dbManager *database.Manager) *Manager {
	metrics := NewMetrics()
	healthCheck := NewHealthChecker(logger, dbManager)

	manager := &Manager{
		metrics:     metrics,
		healthCheck: healthCheck,
		logger:      logger,
		dbManager:   dbManager,
		stopChan:    make(chan struct{}),
	}

	// 启动系统指标更新
	manager.startSystemMetricsUpdater()

	return manager
}

// GetMetrics 获取指标收集器
func (m *Manager) GetMetrics() *Metrics {
	return m.metrics
}

// GetHealthChecker 获取健康检查器
func (m *Manager) GetHealthChecker() *HealthChecker {
	return m.healthCheck
}

// PrometheusHandler Prometheus指标处理器
func (m *Manager) PrometheusHandler() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}

// startSystemMetricsUpdater 启动系统指标更新器
func (m *Manager) startSystemMetricsUpdater() {
	m.updateTicker = time.NewTicker(30 * time.Second)

	go func() {
		for {
			select {
			case <-m.stopChan:
				return
			case <-m.updateTicker.C:
				m.updateSystemMetrics()
			}
		}
	}()
}

// updateSystemMetrics 更新系统指标
func (m *Manager) updateSystemMetrics() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 更新数据库连接数
	if m.dbManager.MySQL != nil {
		if db, err := m.dbManager.MySQL.DB(); err == nil {
			stats := db.Stats()
			m.metrics.SetDatabaseConnections(float64(stats.OpenConnections))
		}
	}

	// 更新Redis连接数
	if m.dbManager.Redis != nil {
		if stats := m.dbManager.Redis.PoolStats(); stats != nil {
			m.metrics.SetRedisConnections(float64(stats.TotalConns))
		}
	}

	m.logger.Debug("System metrics updated", nil)
}

// Close 关闭监控管理器
func (m *Manager) Close() error {
	if m.updateTicker != nil {
		m.updateTicker.Stop()
	}
	close(m.stopChan)
	return nil
}
```

### 11.6 创建业务指标装饰器 (`pkg/metrics/decorators.go`)
```go
package metrics

import (
	"context"
	"time"

	"detect-executor-go/internal/model"
	"detect-executor-go/internal/service"
)

// TaskServiceDecorator 任务服务指标装饰器
type TaskServiceDecorator struct {
	service.TaskService
	metrics *Metrics
}

// NewTaskServiceDecorator 创建任务服务指标装饰器
func NewTaskServiceDecorator(taskService service.TaskService, metrics *Metrics) service.TaskService {
	return &TaskServiceDecorator{
		TaskService: taskService,
		metrics:     metrics,
	}
}

// CreateTask 创建任务（带指标记录）
func (d *TaskServiceDecorator) CreateTask(ctx context.Context, req *service.CreateTaskRequest) (*model.DetectTask, error) {
	start := time.Now()
	
	task, err := d.TaskService.CreateTask(ctx, req)
	
	duration := time.Since(start).Seconds()
	status := "success"
	if err != nil {
		status = "error"
		d.metrics.RecordError("task_service", "create_task")
	}
	
	if task != nil {
		d.metrics.RecordTask(task.DeviceID, string(task.Type), status, duration)
	}
	
	return task, err
}

// UpdateTaskStatus 更新任务状态（带指标记录）
func (d *TaskServiceDecorator) UpdateTaskStatus(ctx context.Context, taskID string, status model.TaskStatus, message string) error {
	err := d.TaskService.UpdateTaskStatus(ctx, taskID, status, message)
	
	if err != nil {
		d.metrics.RecordError("task_service", "update_status")
	} else {
		// 根据状态更新进行中任务数
		switch status {
		case model.TaskStatusRunning:
			d.metrics.IncTasksInProgress()
		case model.TaskStatusCompleted, model.TaskStatusFailed, model.TaskStatusCancelled:
			d.metrics.DecTasksInProgress()
		}
	}
	
	return err
}

// DetectServiceDecorator 检测服务指标装饰器
type DetectServiceDecorator struct {
	service.DetectService
	metrics *Metrics
}

// NewDetectServiceDecorator 创建检测服务指标装饰器
func NewDetectServiceDecorator(detectService service.DetectService, metrics *Metrics) service.DetectService {
	return &DetectServiceDecorator{
		DetectService: detectService,
		metrics:       metrics,
	}
}

// StartDetection 开始检测（带指标记录）
func (d *DetectServiceDecorator) StartDetection(ctx context.Context, taskID string) error {
	start := time.Now()
	
	err := d.DetectService.StartDetection(ctx, taskID)
	
	duration := time.Since(start).Seconds()
	result := "success"
	if err != nil {
		result = "error"
		d.metrics.RecordError("detect_service", "start_detection")
	}
	
	d.metrics.RecordDetection("start", result, duration)
	
	return err
}

// ResultServiceDecorator 结果服务指标装饰器
type ResultServiceDecorator struct {
	service.ResultService
	metrics *Metrics
}

// NewResultServiceDecorator 创建结果服务指标装饰器
func NewResultServiceDecorator(resultService service.ResultService, metrics *Metrics) service.ResultService {
	return &ResultServiceDecorator{
		ResultService: resultService,
		metrics:       metrics,
	}
}

// ListResults 列表结果（带指标记录）
func (d *ResultServiceDecorator) ListResults(ctx context.Context, req *service.ListResultsRequest) ([]*model.DetectResult, int64, error) {
	start := time.Now()
	
	results, total, err := d.ResultService.ListResults(ctx, req)
	
	duration := time.Since(start).Seconds()
	if err != nil {
		d.metrics.RecordError("result_service", "list_results")
	}
	
	// 记录查询性能指标
	d.metrics.RecordExternalAPIRequest("result_service", "list", "200", duration)
	
	// 记录结果统计
	for _, result := range results {
		d.metrics.RecordResult(result.DeviceID, string(result.Type), result.Passed)
	}
	
	return results, total, err
}
```

### 11.7 更新配置文件 (`configs/config.yaml`)
```yaml
# 监控配置
metrics:
  enabled: true
  port: 9090
  path: "/metrics"
  
health:
  enabled: true
  liveness_path: "/live"
  readiness_path: "/ready"
  
# 告警配置
alerting:
  enabled: true
  webhook_url: "http://localhost:9093/api/v1/alerts"
  rules:
    - name: "high_error_rate"
      condition: "rate(errors_total[5m]) > 0.1"
      severity: "warning"
    - name: "task_queue_full"
      condition: "tasks_in_progress > 100"
      severity: "critical"
```

### 11.8 更新配置结构 (`pkg/config/config.go`)
```go
// 在原有Config结构体中添加监控配置
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	Detect   DetectConfig   `yaml:"detect"`
	Log      LogConfig      `yaml:"log"`
	Metrics  MetricsConfig  `yaml:"metrics"`
	Client   ClientConfig   `yaml:"client"`
	Health   HealthConfig   `yaml:"health"`
	Alerting AlertingConfig `yaml:"alerting"`
}

// MetricsConfig 监控配置
type MetricsConfig struct {
	Enabled bool   `yaml:"enabled"`
	Port    int    `yaml:"port"`
	Path    string `yaml:"path"`
}

// HealthConfig 健康检查配置
type HealthConfig struct {
	Enabled       bool   `yaml:"enabled"`
	LivenessPath  string `yaml:"liveness_path"`
	ReadinessPath string `yaml:"readiness_path"`
}

// AlertingConfig 告警配置
type AlertingConfig struct {
	Enabled    bool        `yaml:"enabled"`
	WebhookURL string      `yaml:"webhook_url"`
	Rules      []AlertRule `yaml:"rules"`
}

// AlertRule 告警规则
type AlertRule struct {
	Name      string `yaml:"name"`
	Condition string `yaml:"condition"`
	Severity  string `yaml:"severity"`
}
```

### 11.9 创建告警系统 (`pkg/metrics/alerting.go`)
```go
package metrics

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"detect-executor-go/pkg/config"
	"detect-executor-go/pkg/logger"
)

// AlertManager 告警管理器
type AlertManager struct {
	config     *config.AlertingConfig
	logger     *logger.Logger
	httpClient *http.Client
	alerts     chan Alert
	stopChan   chan struct{}
}

// Alert 告警信息
type Alert struct {
	Name        string            `json:"name"`
	Severity    string            `json:"severity"`
	Message     string            `json:"message"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	Timestamp   time.Time         `json:"timestamp"`
}

// NewAlertManager 创建告警管理器
func NewAlertManager(config *config.AlertingConfig, logger *logger.Logger) *AlertManager {
	if !config.Enabled {
		return nil
	}

	am := &AlertManager{
		config: config,
		logger: logger,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		alerts:   make(chan Alert, 100),
		stopChan: make(chan struct{}),
	}

	// 启动告警处理器
	go am.processAlerts()

	return am
}

// SendAlert 发送告警
func (am *AlertManager) SendAlert(alert Alert) {
	if am == nil {
		return
	}

	alert.Timestamp = time.Now()

	select {
	case am.alerts <- alert:
	default:
		am.logger.Warn("Alert queue full, dropping alert", map[string]interface{}{
			"alert_name": alert.Name,
		})
	}
}

// processAlerts 处理告警
func (am *AlertManager) processAlerts() {
	for {
		select {
		case <-am.stopChan:
			return
		case alert := <-am.alerts:
			am.sendToWebhook(alert)
		}
	}
}

// sendToWebhook 发送告警到Webhook
func (am *AlertManager) sendToWebhook(alert Alert) {
	if am.config.WebhookURL == "" {
		return
	}

	payload, err := json.Marshal(alert)
	if err != nil {
		am.logger.Error("Failed to marshal alert", map[string]interface{}{
			"error": err.Error(),
			"alert": alert.Name,
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", am.config.WebhookURL, bytes.NewBuffer(payload))
	if err != nil {
		am.logger.Error("Failed to create webhook request", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := am.httpClient.Do(req)
	if err != nil {
		am.logger.Error("Failed to send alert webhook", map[string]interface{}{
			"error": err.Error(),
			"alert": alert.Name,
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		am.logger.Error("Webhook returned error status", map[string]interface{}{
			"status_code": resp.StatusCode,
			"alert":       alert.Name,
		})
		return
	}

	am.logger.Info("Alert sent successfully", map[string]interface{}{
		"alert":    alert.Name,
		"severity": alert.Severity,
	})
}

// Close 关闭告警管理器
func (am *AlertManager) Close() error {
	if am == nil {
		return nil
	}
	close(am.stopChan)
	return nil
}
```

### 11.10 更新主程序 (`cmd/server/main.go`)
```go
// 在原有main函数中添加监控初始化
func main() {
	// ... 原有初始化代码 ...

	// 初始化监控管理器
	metricsManager := metrics.NewManager(log, dbManager)
	defer metricsManager.Close()

	// 初始化告警管理器
	alertManager := metrics.NewAlertManager(&cfg.Alerting, log)
	defer alertManager.Close()

	// 初始化服务层（使用指标装饰器）
	serviceManager, err := service.NewManager(log, repoManager, clientManager)
	if err != nil {
		log.Fatal("Failed to create service manager", map[string]interface{}{"error": err})
	}

	// 装饰服务以添加指标收集
	serviceManager.Task = metrics.NewTaskServiceDecorator(serviceManager.Task, metricsManager.GetMetrics())
	serviceManager.Detect = metrics.NewDetectServiceDecorator(serviceManager.Detect, metricsManager.GetMetrics())
	serviceManager.Result = metrics.NewResultServiceDecorator(serviceManager.Result, metricsManager.GetMetrics())

	defer serviceManager.Close()

	// 初始化路由（添加指标中间件）
	router := handler.NewRouter(log, serviceManager)
	
	// 添加监控中间件
	router.GetEngine().Use(metrics.HTTPMetricsMiddleware(metricsManager.GetMetrics()))
	
	// 添加监控路由
	router.GetEngine().GET("/metrics", metricsManager.PrometheusHandler())
	router.GetEngine().GET("/live", metricsManager.GetHealthChecker().LivenessHandler())
	router.GetEngine().GET("/ready", metricsManager.GetHealthChecker().ReadinessHandler())

	// ... 剩余代码保持不变 ...
}
```

## 验证步骤

### 11.11 运行测试验证
```bash
# 添加Prometheus依赖
go get github.com/prometheus/client_golang
go get github.com/heptiolabs/healthcheck

# 检查依赖
go mod tidy

# 运行监控测试
go test ./pkg/metrics/ -v

# 启动服务器
go run cmd/server/main.go

# 测试监控接口
curl http://localhost:8080/metrics
curl http://localhost:8080/live
curl http://localhost:8080/ready
```

## 预期结果
- 完整的Prometheus指标集成
- HTTP请求和业务指标自动收集
- 完善的健康检查和就绪检查
- 系统资源监控和告警
- 错误追踪和性能分析
- 可视化监控面板支持
- 自动化告警和通知机制
- 完善的单元测试覆盖

## 状态
- [ ] 步骤 11 完成：监控指标实现

---

## 下一步预告
步骤 12: 测试用例实现 (`tests/`)
- 单元测试和集成测试
- 性能测试和压力测试
- 测试数据和模拟环境
- CI/CD集成和自动化测试
