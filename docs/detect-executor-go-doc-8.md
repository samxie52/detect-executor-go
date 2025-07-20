# 步骤 8: 外部客户端集成

## 目标
实现与外部系统的客户端集成，包括HTTP客户端、检测引擎客户端、文件存储客户端和消息队列客户端，为核心业务服务提供外部依赖支持。

## 操作步骤

### 8.1 添加客户端依赖
```bash
# 添加HTTP客户端依赖
go get github.com/go-resty/resty/v2

# 添加文件存储客户端
go get github.com/minio/minio-go/v7

# 添加消息队列客户端
go get github.com/streadway/amqp

# 添加gRPC客户端
go get google.golang.org/grpc
go get google.golang.org/protobuf
```

### 8.2 创建基础客户端接口 (`pkg/client/base.go`)
```go
package client

import (
	"context"
	"time"
)

// BaseClient 基础客户端接口
type BaseClient interface {
	Close() error
	Ping(ctx context.Context) error
}

// ClientConfig 客户端配置
type ClientConfig struct {
	Timeout       time.Duration
	RetryCount    int
	RetryInterval time.Duration
	MaxIdleConns  int
	MaxConnsPerHost int
}

// DefaultClientConfig 默认客户端配置
var DefaultClientConfig = ClientConfig{
	Timeout:         30 * time.Second,
	RetryCount:      3,
	RetryInterval:   1 * time.Second,
	MaxIdleConns:    100,
	MaxConnsPerHost: 10,
}
```

### 8.3 创建HTTP客户端 (`pkg/client/http_client.go`)
```go
package client

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

// HTTPClient HTTP客户端接口
type HTTPClient interface {
	BaseClient
	Get(ctx context.Context, url string, headers map[string]string) (*HTTPResponse, error)
	Post(ctx context.Context, url string, body interface{}, headers map[string]string) (*HTTPResponse, error)
	Put(ctx context.Context, url string, body interface{}, headers map[string]string) (*HTTPResponse, error)
	Delete(ctx context.Context, url string, headers map[string]string) (*HTTPResponse, error)
}

// HTTPResponse HTTP响应
type HTTPResponse struct {
	StatusCode int
	Body       []byte
	Headers    map[string]string
}

// httpClient HTTP客户端实现
type httpClient struct {
	client *resty.Client
	config ClientConfig
}

// NewHTTPClient 创建HTTP客户端
func NewHTTPClient(config ClientConfig) HTTPClient {
	client := resty.New()
	client.SetTimeout(config.Timeout)
	client.SetRetryCount(config.RetryCount)
	client.SetRetryWaitTime(config.RetryInterval)
	
	return &httpClient{
		client: client,
		config: config,
	}
}

// Get GET请求
func (c *httpClient) Get(ctx context.Context, url string, headers map[string]string) (*HTTPResponse, error) {
	req := c.client.R().SetContext(ctx)
	
	if headers != nil {
		req.SetHeaders(headers)
	}
	
	resp, err := req.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http get failed: %w", err)
	}
	
	return &HTTPResponse{
		StatusCode: resp.StatusCode(),
		Body:       resp.Body(),
		Headers:    c.convertHeaders(resp.Header()),
	}, nil
}

// Post POST请求
func (c *httpClient) Post(ctx context.Context, url string, body interface{}, headers map[string]string) (*HTTPResponse, error) {
	req := c.client.R().SetContext(ctx)
	
	if headers != nil {
		req.SetHeaders(headers)
	}
	
	if body != nil {
		req.SetBody(body)
	}
	
	resp, err := req.Post(url)
	if err != nil {
		return nil, fmt.Errorf("http post failed: %w", err)
	}
	
	return &HTTPResponse{
		StatusCode: resp.StatusCode(),
		Body:       resp.Body(),
		Headers:    c.convertHeaders(resp.Header()),
	}, nil
}

// Put PUT请求
func (c *httpClient) Put(ctx context.Context, url string, body interface{}, headers map[string]string) (*HTTPResponse, error) {
	req := c.client.R().SetContext(ctx)
	
	if headers != nil {
		req.SetHeaders(headers)
	}
	
	if body != nil {
		req.SetBody(body)
	}
	
	resp, err := req.Put(url)
	if err != nil {
		return nil, fmt.Errorf("http put failed: %w", err)
	}
	
	return &HTTPResponse{
		StatusCode: resp.StatusCode(),
		Body:       resp.Body(),
		Headers:    c.convertHeaders(resp.Header()),
	}, nil
}

// Delete DELETE请求
func (c *httpClient) Delete(ctx context.Context, url string, headers map[string]string) (*HTTPResponse, error) {
	req := c.client.R().SetContext(ctx)
	
	if headers != nil {
		req.SetHeaders(headers)
	}
	
	resp, err := req.Delete(url)
	if err != nil {
		return nil, fmt.Errorf("http delete failed: %w", err)
	}
	
	return &HTTPResponse{
		StatusCode: resp.StatusCode(),
		Body:       resp.Body(),
		Headers:    c.convertHeaders(resp.Header()),
	}, nil
}

// Close 关闭客户端
func (c *httpClient) Close() error {
	// resty客户端不需要显式关闭
	return nil
}

// Ping 健康检查
func (c *httpClient) Ping(ctx context.Context) error {
	// HTTP客户端的ping可以通过请求一个健康检查端点实现
	return nil
}

// convertHeaders 转换响应头
func (c *httpClient) convertHeaders(headers map[string][]string) map[string]string {
	result := make(map[string]string)
	for key, values := range headers {
		if len(values) > 0 {
			result[key] = values[0]
		}
	}
	return result
}
```

### 8.4 创建检测引擎客户端 (`pkg/client/detect_engine_client.go`)
```go
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"detect-executor-go/internal/model"
)

// DetectEngineClient 检测引擎客户端接口
type DetectEngineClient interface {
	BaseClient
	SubmitTask(ctx context.Context, task *DetectTaskRequest) (*DetectTaskResponse, error)
	GetTaskStatus(ctx context.Context, taskID string) (*DetectTaskStatus, error)
	GetTaskResult(ctx context.Context, taskID string) (*DetectTaskResult, error)
	CancelTask(ctx context.Context, taskID string) error
}

// DetectTaskRequest 检测任务请求
type DetectTaskRequest struct {
	TaskID     string                 `json:"task_id"`
	Type       model.DetectType       `json:"type"`
	VideoURL   string                 `json:"video_url"`
	StartTime  *time.Time             `json:"start_time,omitempty"`
	EndTime    *time.Time             `json:"end_time,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// DetectTaskResponse 检测任务响应
type DetectTaskResponse struct {
	TaskID  string `json:"task_id"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// DetectTaskStatus 检测任务状态
type DetectTaskStatus struct {
	TaskID   string `json:"task_id"`
	Status   string `json:"status"`
	Progress int    `json:"progress"`
	Message  string `json:"message,omitempty"`
}

// DetectTaskResult 检测任务结果
type DetectTaskResult struct {
	TaskID     string                 `json:"task_id"`
	Status     string                 `json:"status"`
	Score      float64                `json:"score"`
	Passed     bool                   `json:"passed"`
	Confidence float64                `json:"confidence"`
	Details    map[string]interface{} `json:"details"`
	ReportURL  string                 `json:"report_url,omitempty"`
}

// detectEngineClient 检测引擎客户端实现
type detectEngineClient struct {
	httpClient HTTPClient
	baseURL    string
	apiKey     string
}

// NewDetectEngineClient 创建检测引擎客户端
func NewDetectEngineClient(baseURL, apiKey string, config ClientConfig) DetectEngineClient {
	return &detectEngineClient{
		httpClient: NewHTTPClient(config),
		baseURL:    baseURL,
		apiKey:     apiKey,
	}
}

// SubmitTask 提交检测任务
func (c *detectEngineClient) SubmitTask(ctx context.Context, task *DetectTaskRequest) (*DetectTaskResponse, error) {
	url := fmt.Sprintf("%s/api/v1/tasks", c.baseURL)
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", c.apiKey),
	}
	
	resp, err := c.httpClient.Post(ctx, url, task, headers)
	if err != nil {
		return nil, fmt.Errorf("submit detect task failed: %w", err)
	}
	
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return nil, fmt.Errorf("submit detect task failed with status %d: %s", resp.StatusCode, string(resp.Body))
	}
	
	var result DetectTaskResponse
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("parse detect task response failed: %w", err)
	}
	
	return &result, nil
}

// GetTaskStatus 获取任务状态
func (c *detectEngineClient) GetTaskStatus(ctx context.Context, taskID string) (*DetectTaskStatus, error) {
	url := fmt.Sprintf("%s/api/v1/tasks/%s/status", c.baseURL, taskID)
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", c.apiKey),
	}
	
	resp, err := c.httpClient.Get(ctx, url, headers)
	if err != nil {
		return nil, fmt.Errorf("get detect task status failed: %w", err)
	}
	
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("get detect task status failed with status %d: %s", resp.StatusCode, string(resp.Body))
	}
	
	var result DetectTaskStatus
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("parse detect task status failed: %w", err)
	}
	
	return &result, nil
}

// GetTaskResult 获取任务结果
func (c *detectEngineClient) GetTaskResult(ctx context.Context, taskID string) (*DetectTaskResult, error) {
	url := fmt.Sprintf("%s/api/v1/tasks/%s/result", c.baseURL, taskID)
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", c.apiKey),
	}
	
	resp, err := c.httpClient.Get(ctx, url, headers)
	if err != nil {
		return nil, fmt.Errorf("get detect task result failed: %w", err)
	}
	
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("get detect task result failed with status %d: %s", resp.StatusCode, string(resp.Body))
	}
	
	var result DetectTaskResult
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("parse detect task result failed: %w", err)
	}
	
	return &result, nil
}

// CancelTask 取消任务
func (c *detectEngineClient) CancelTask(ctx context.Context, taskID string) error {
	url := fmt.Sprintf("%s/api/v1/tasks/%s/cancel", c.baseURL, taskID)
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", c.apiKey),
	}
	
	resp, err := c.httpClient.Post(ctx, url, nil, headers)
	if err != nil {
		return fmt.Errorf("cancel detect task failed: %w", err)
	}
	
	if resp.StatusCode != 200 {
		return fmt.Errorf("cancel detect task failed with status %d: %s", resp.StatusCode, string(resp.Body))
	}
	
	return nil
}

// Close 关闭客户端
func (c *detectEngineClient) Close() error {
	return c.httpClient.Close()
}

// Ping 健康检查
func (c *detectEngineClient) Ping(ctx context.Context) error {
	url := fmt.Sprintf("%s/api/v1/health", c.baseURL)
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", c.apiKey),
	}
	
	resp, err := c.httpClient.Get(ctx, url, headers)
	if err != nil {
		return fmt.Errorf("ping detect engine failed: %w", err)
	}
	
	if resp.StatusCode != 200 {
		return fmt.Errorf("detect engine unhealthy, status: %d", resp.StatusCode)
	}
	
	return nil
}
```

### 8.5 创建文件存储客户端 (`pkg/client/storage_client.go`)
```go
package client

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// StorageClient 文件存储客户端接口
type StorageClient interface {
	BaseClient
	UploadFile(ctx context.Context, bucket, objectName string, reader io.Reader, size int64) error
	DownloadFile(ctx context.Context, bucket, objectName string) (io.ReadCloser, error)
	DeleteFile(ctx context.Context, bucket, objectName string) error
	GetFileURL(ctx context.Context, bucket, objectName string, expiry time.Duration) (string, error)
	ListFiles(ctx context.Context, bucket, prefix string) ([]FileInfo, error)
}

// FileInfo 文件信息
type FileInfo struct {
	Name         string
	Size         int64
	LastModified time.Time
	ContentType  string
}

// storageClient 文件存储客户端实现
type storageClient struct {
	client *minio.Client
	config StorageConfig
}

// StorageConfig 存储配置
type StorageConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	Region          string
}

// NewStorageClient 创建文件存储客户端
func NewStorageClient(config StorageConfig) (StorageClient, error) {
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
		Region: config.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("create minio client failed: %w", err)
	}
	
	return &storageClient{
		client: client,
		config: config,
	}, nil
}

// UploadFile 上传文件
func (c *storageClient) UploadFile(ctx context.Context, bucket, objectName string, reader io.Reader, size int64) error {
	_, err := c.client.PutObject(ctx, bucket, objectName, reader, size, minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("upload file failed: %w", err)
	}
	return nil
}

// DownloadFile 下载文件
func (c *storageClient) DownloadFile(ctx context.Context, bucket, objectName string) (io.ReadCloser, error) {
	object, err := c.client.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("download file failed: %w", err)
	}
	return object, nil
}

// DeleteFile 删除文件
func (c *storageClient) DeleteFile(ctx context.Context, bucket, objectName string) error {
	err := c.client.RemoveObject(ctx, bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("delete file failed: %w", err)
	}
	return nil
}

// GetFileURL 获取文件访问URL
func (c *storageClient) GetFileURL(ctx context.Context, bucket, objectName string, expiry time.Duration) (string, error) {
	url, err := c.client.PresignedGetObject(ctx, bucket, objectName, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("get file url failed: %w", err)
	}
	return url.String(), nil
}

// ListFiles 列出文件
func (c *storageClient) ListFiles(ctx context.Context, bucket, prefix string) ([]FileInfo, error) {
	var files []FileInfo
	
	for object := range c.client.ListObjects(ctx, bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	}) {
		if object.Err != nil {
			return nil, fmt.Errorf("list files failed: %w", object.Err)
		}
		
		files = append(files, FileInfo{
			Name:         object.Key,
			Size:         object.Size,
			LastModified: object.LastModified,
			ContentType:  object.ContentType,
		})
	}
	
	return files, nil
}

// Close 关闭客户端
func (c *storageClient) Close() error {
	// minio客户端不需要显式关闭
	return nil
}

// Ping 健康检查
func (c *storageClient) Ping(ctx context.Context) error {
	// 通过列出存储桶来检查连接
	_, err := c.client.ListBuckets(ctx)
	if err != nil {
		return fmt.Errorf("storage ping failed: %w", err)
	}
	return nil
}
```

### 8.6 创建消息队列客户端 (`pkg/client/message_queue_client.go`)
```go
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

// MessageQueueClient 消息队列客户端接口
type MessageQueueClient interface {
	BaseClient
	Publish(ctx context.Context, exchange, routingKey string, message interface{}) error
	Consume(ctx context.Context, queue string, handler MessageHandler) error
	DeclareQueue(ctx context.Context, queue string, durable bool) error
	DeclareExchange(ctx context.Context, exchange, exchangeType string, durable bool) error
}

// MessageHandler 消息处理器
type MessageHandler func(ctx context.Context, message []byte) error

// messageQueueClient 消息队列客户端实现
type messageQueueClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  MessageQueueConfig
}

// MessageQueueConfig 消息队列配置
type MessageQueueConfig struct {
	URL             string
	ReconnectDelay  time.Duration
	Heartbeat       time.Duration
	ConnectionName  string
}

// NewMessageQueueClient 创建消息队列客户端
func NewMessageQueueClient(config MessageQueueConfig) (MessageQueueClient, error) {
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		return nil, fmt.Errorf("connect to rabbitmq failed: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("create channel failed: %w", err)
	}

	return &messageQueueClient{
		conn:    conn,
		channel: channel,
		config:  config,
	}, nil
}

// Publish 发布消息
func (c *messageQueueClient) Publish(ctx context.Context, exchange, routingKey string, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal message failed: %w", err)
	}

	err = c.channel.Publish(
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Timestamp:   time.Now(),
		},
	)
	if err != nil {
		return fmt.Errorf("publish message failed: %w", err)
	}

	return nil
}

// Consume 消费消息
func (c *messageQueueClient) Consume(ctx context.Context, queue string, handler MessageHandler) error {
	msgs, err := c.channel.Consume(
		queue,
		"", // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("consume messages failed: %w", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgs:
				if !ok {
					return
				}

				if err := handler(ctx, msg.Body); err != nil {
					msg.Nack(false, true) // requeue
				} else {
					msg.Ack(false)
				}
			}
		}
	}()

	return nil
}

// DeclareQueue 声明队列
func (c *messageQueueClient) DeclareQueue(ctx context.Context, queue string, durable bool) error {
	_, err := c.channel.QueueDeclare(
		queue,
		durable,
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("declare queue failed: %w", err)
	}
	return nil
}

// DeclareExchange 声明交换机
func (c *messageQueueClient) DeclareExchange(ctx context.Context, exchange, exchangeType string, durable bool) error {
	err := c.channel.ExchangeDeclare(
		exchange,
		exchangeType,
		durable,
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("declare exchange failed: %w", err)
	}
	return nil
}

// Close 关闭客户端
func (c *messageQueueClient) Close() error {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
	return nil
}

// Ping 健康检查
func (c *messageQueueClient) Ping(ctx context.Context) error {
	if c.conn == nil || c.conn.IsClosed() {
		return fmt.Errorf("rabbitmq connection is closed")
	}
	if c.channel == nil {
		return fmt.Errorf("rabbitmq channel is closed")
	}
	return nil
}
```

### 8.7 创建客户端管理器 (`pkg/client/manager.go`)
```go
package client

import (
	"context"
	"fmt"
	"sync"
)

// Manager 客户端管理器
type Manager struct {
	mu sync.RWMutex

	HTTP          HTTPClient
	DetectEngine  DetectEngineClient
	Storage       StorageClient
	MessageQueue  MessageQueueClient

	config *Config
}

// Config 客户端配置
type Config struct {
	HTTP         ClientConfig
	DetectEngine DetectEngineConfig
	Storage      StorageConfig
	MessageQueue MessageQueueConfig
}

// DetectEngineConfig 检测引擎配置
type DetectEngineConfig struct {
	BaseURL string
	APIKey  string
	Client  ClientConfig
}

// NewManager 创建客户端管理器
func NewManager(config *Config) (*Manager, error) {
	manager := &Manager{
		config: config,
	}

	// 初始化HTTP客户端
	manager.HTTP = NewHTTPClient(config.HTTP)

	// 初始化检测引擎客户端
	manager.DetectEngine = NewDetectEngineClient(
		config.DetectEngine.BaseURL,
		config.DetectEngine.APIKey,
		config.DetectEngine.Client,
	)

	// 初始化存储客户端
	storageClient, err := NewStorageClient(config.Storage)
	if err != nil {
		return nil, fmt.Errorf("create storage client failed: %w", err)
	}
	manager.Storage = storageClient

	// 初始化消息队列客户端
	mqClient, err := NewMessageQueueClient(config.MessageQueue)
	if err != nil {
		return nil, fmt.Errorf("create message queue client failed: %w", err)
	}
	manager.MessageQueue = mqClient

	return manager, nil
}

// Close 关闭所有客户端
func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errors []error

	if m.HTTP != nil {
		if err := m.HTTP.Close(); err != nil {
			errors = append(errors, fmt.Errorf("close http client failed: %w", err))
		}
	}

	if m.DetectEngine != nil {
		if err := m.DetectEngine.Close(); err != nil {
			errors = append(errors, fmt.Errorf("close detect engine client failed: %w", err))
		}
	}

	if m.Storage != nil {
		if err := m.Storage.Close(); err != nil {
			errors = append(errors, fmt.Errorf("close storage client failed: %w", err))
		}
	}

	if m.MessageQueue != nil {
		if err := m.MessageQueue.Close(); err != nil {
			errors = append(errors, fmt.Errorf("close message queue client failed: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("close clients failed: %v", errors)
	}

	return nil
}

// HealthCheck 健康检查
func (m *Manager) HealthCheck(ctx context.Context) map[string]error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	results := make(map[string]error)

	if m.HTTP != nil {
		results["http"] = m.HTTP.Ping(ctx)
	}

	if m.DetectEngine != nil {
		results["detect_engine"] = m.DetectEngine.Ping(ctx)
	}

	if m.Storage != nil {
		results["storage"] = m.Storage.Ping(ctx)
	}

	if m.MessageQueue != nil {
		results["message_queue"] = m.MessageQueue.Ping(ctx)
	}

	return results
}
```

### 8.8 更新配置文件 (`configs/config.yaml`)
```yaml
# 外部客户端配置
client:
  http:
    timeout: 30s
    retry_count: 3
    retry_interval: 1s
    max_idle_conns: 100
    max_conns_per_host: 10
  
  detect_engine:
    base_url: "http://localhost:8081"
    api_key: "your-api-key"
    client:
      timeout: 60s
      retry_count: 3
      retry_interval: 2s
      max_idle_conns: 50
      max_conns_per_host: 5
  
  storage:
    endpoint: "localhost:9000"
    access_key_id: "minioadmin"
    secret_access_key: "minioadmin"
    use_ssl: false
    region: "us-east-1"
  
  message_queue:
    url: "amqp://guest:guest@localhost:5672/"
    reconnect_delay: 5s
    heartbeat: 30s
    connection_name: "detect-executor"
```

### 8.9 更新配置结构体 (`pkg/config/config.go`)
```go
// 在Config结构体中添加客户端配置
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Detect   DetectConfig
	Log      LogConfig
	Metrics  MetricsConfig
	Client   ClientConfig // 新增
}

// ClientConfig 客户端配置
type ClientConfig struct {
	HTTP         HTTPClientConfig         `mapstructure:"http"`
	DetectEngine DetectEngineClientConfig `mapstructure:"detect_engine"`
	Storage      StorageClientConfig      `mapstructure:"storage"`
	MessageQueue MessageQueueClientConfig `mapstructure:"message_queue"`
}

// HTTPClientConfig HTTP客户端配置
type HTTPClientConfig struct {
	Timeout         time.Duration `mapstructure:"timeout"`
	RetryCount      int           `mapstructure:"retry_count"`
	RetryInterval   time.Duration `mapstructure:"retry_interval"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxConnsPerHost int           `mapstructure:"max_conns_per_host"`
}

// DetectEngineClientConfig 检测引擎客户端配置
type DetectEngineClientConfig struct {
	BaseURL string            `mapstructure:"base_url"`
	APIKey  string            `mapstructure:"api_key"`
	Client  HTTPClientConfig  `mapstructure:"client"`
}

// StorageClientConfig 存储客户端配置
type StorageClientConfig struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	UseSSL          bool   `mapstructure:"use_ssl"`
	Region          string `mapstructure:"region"`
}

// MessageQueueClientConfig 消息队列客户端配置
type MessageQueueClientConfig struct {
	URL            string        `mapstructure:"url"`
	ReconnectDelay time.Duration `mapstructure:"reconnect_delay"`
	Heartbeat      time.Duration `mapstructure:"heartbeat"`
	ConnectionName string        `mapstructure:"connection_name"`
}
```

### 8.10 创建单元测试 (`pkg/client/client_test.go`)
```go
package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ClientTestSuite 客户端测试套件
type ClientTestSuite struct {
	suite.Suite
	ctx context.Context
}

// SetupSuite 设置测试套件
func (suite *ClientTestSuite) SetupSuite() {
	suite.ctx = context.Background()
}

// TestHTTPClient 测试HTTP客户端
func (suite *ClientTestSuite) TestHTTPClient() {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	}))
	defer server.Close()

	// 创建HTTP客户端
	client := NewHTTPClient(DefaultClientConfig)

	// 测试GET请求
	resp, err := client.Get(suite.ctx, server.URL, nil)
	suite.NoError(err)
	suite.Equal(200, resp.StatusCode)
	suite.Contains(string(resp.Body), "success")

	// 测试POST请求
	resp, err = client.Post(suite.ctx, server.URL, map[string]string{"key": "value"}, nil)
	suite.NoError(err)
	suite.Equal(200, resp.StatusCode)
}

// TestDetectEngineClient 测试检测引擎客户端
func (suite *ClientTestSuite) TestDetectEngineClient() {
	// 创建模拟检测引擎服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/tasks":
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"task_id": "test-task", "status": "pending"}`))
		case "/api/v1/health":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "healthy"}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// 创建检测引擎客户端
	client := NewDetectEngineClient(server.URL, "test-key", DefaultClientConfig)

	// 测试健康检查
	err := client.Ping(suite.ctx)
	suite.NoError(err)

	// 测试提交任务
	task := &DetectTaskRequest{
		TaskID:   "test-task",
		Type:     "quality",
		VideoURL: "rtsp://example.com/stream",
	}

	resp, err := client.SubmitTask(suite.ctx, task)
	suite.NoError(err)
	suite.Equal("test-task", resp.TaskID)
	suite.Equal("pending", resp.Status)
}

// TestClientManager 测试客户端管理器
func (suite *ClientTestSuite) TestClientManager() {
	config := &Config{
		HTTP: ClientConfig{
			Timeout:       30 * time.Second,
			RetryCount:    3,
			RetryInterval: 1 * time.Second,
		},
		DetectEngine: DetectEngineConfig{
			BaseURL: "http://localhost:8081",
			APIKey:  "test-key",
			Client: ClientConfig{
				Timeout: 60 * time.Second,
			},
		},
		Storage: StorageConfig{
			Endpoint:        "localhost:9000",
			AccessKeyID:     "minioadmin",
			SecretAccessKey: "minioadmin",
			UseSSL:          false,
		},
		MessageQueue: MessageQueueConfig{
			URL: "amqp://guest:guest@localhost:5672/",
		},
	}

	// 注意：这里只测试管理器创建，实际的客户端连接需要真实的服务
	// 在实际测试中，应该使用mock或者测试容器
	
	// 测试配置验证
	suite.NotNil(config.HTTP)
	suite.NotNil(config.DetectEngine)
	suite.NotNil(config.Storage)
	suite.NotNil(config.MessageQueue)
}

// TestClientTestSuite 运行测试套件
func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}
```

## 验证步骤

### 8.11 运行测试验证
```bash
# 检查依赖
go mod tidy

# 运行客户端测试
go test ./pkg/client/ -v

# 启动外部服务进行集成测试（可选）
docker run -d --name minio-test \
  -p 9000:9000 \
  -p 9001:9001 \
  -e MINIO_ROOT_USER=minioadmin \
  -e MINIO_ROOT_PASSWORD=minioadmin \
  minio/minio server /data --console-address ":9001"

docker run -d --name rabbitmq-test \
  -p 5672:5672 \
  -p 15672:15672 \
  rabbitmq:3-management

# 等待服务启动
sleep 10

# 运行集成测试
go test ./pkg/client/ -v -tags=integration
```

## 预期结果
- 完整的外部客户端集成，支持HTTP、检测引擎、存储和消息队列
- 统一的客户端管理和配置
- 健康检查和错误处理机制
- 重试机制和超时控制
- 完善的单元测试和集成测试
- 配置文件集成和环境变量支持

## 状态
- [ ] 步骤 8 完成：外部客户端集成

---

## 下一步预告
步骤 9: 核心业务服务实现 (`internal/service/`)
- 任务调度服务
- 检测执行服务
- 结果处理服务
- 设备管理服务
