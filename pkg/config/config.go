package config

import (
	"fmt"
	"time"
)

// Config 应用程序配置结构
type Config struct {
	Server   ServerConfig   `yaml:"server" mapstructure:"server"`
	Database DatabaseConfig `yaml:"database" mapstructure:"database"`
	Redis    RedisConfig    `yaml:"redis" mapstructure:"redis"`
	Detect   DetectConfig   `yaml:"detect" mapstructure:"detect"`
	Log      LogConfig      `yaml:"log" mapstructure:"log"`
	Metrics  MetricsConfig  `yaml:"metrics" mapstructure:"metrics"`
	Client   ClientConfig   `yaml:"client" mapstructure:"client"`
}

// ClientConfig 客户端配置
type ClientConfig struct {
	HTTP         HTTPClientConfig   `yaml:"http" mapstructure:"http"`
	DetectEngine DetectEngineConfig `yaml:"detect_engine" mapstructure:"detect_engine"`
	Storage      StorageConfig      `yaml:"storage" mapstructure:"storage"`
	MessageQueue MessageQueueConfig `yaml:"message_queue" mapstructure:"message_queue"`
}

// HTTPClientConfig HTTP 客户端配置
type HTTPClientConfig struct {
	Timeout         time.Duration `mapstructure:"timeout"`
	RetryCount      int           `mapstructure:"retry_count"`
	RetryInterval   time.Duration `mapstructure:"retry_interval"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxConnsPerHost int           `mapstructure:"max_conns_per_host"`
}

// DefaultHTTPClientConfig 默认 HTTP 客户端配置
var DefaultHTTPClientConfig = HTTPClientConfig{
	Timeout:         time.Second * 30,
	RetryCount:      3,
	RetryInterval:   time.Second * 1,
	MaxIdleConns:    100,
	MaxConnsPerHost: 10,
}

// DetectEngineConfig 检测引擎配置
type DetectEngineConfig struct {
	BaseURL string           `mapstructure:"base_url"`
	APIKey  string           `mapstructure:"api_key"`
	Client  HTTPClientConfig `mapstructure:"client"`
}

// StorageConfig 存储配置
type StorageConfig struct {
	// MinIO 服务器地址
	Endpoint string `mapstructure:"endpoint"`
	// MinIO 访问密钥
	AccessKey string `mapstructure:"access_key"`
	// MinIO 访问密钥
	SecretKey string `mapstructure:"secret_key"`
	// MinIO 是否使用 SSL
	UseSSL bool `mapstructure:"use_ssl"`
	// MinIO 区域,表示存储桶的地理位置
	Region string `mapstructure:"region"`
}

// MessageQueueConfig 消息队列配置
type MessageQueueConfig struct {
	// RabbitMQ 服务器地址
	URL string `mapstructure:"url"`
	// 重连延迟时间
	ReconnectDelay time.Duration `mapstructure:"reconnect_delay"`
	// 心跳时间
	Heartbeat time.Duration `mapstructure:"heartbeat"`
	// 连接名称
	ConnectionName string `mapstructure:"connection_name"`
	// 连接数量
	ConnectionNumber int `mapstructure:"connection_number"`
}

// ServerConfig HTTP Server configuration
type ServerConfig struct {
	Mode         string        `yaml:"mode" mapstructure:"mode"`
	Host         string        `yaml:"host" mapstructure:"host"`
	Port         int           `yaml:"port" mapstructure:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout" mapstructure:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout" mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" mapstructure:"idle_timeout"`
}

// DatabaseConfig Database configuration
type DatabaseConfig struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Port     int    `yaml:"port" mapstructure:"port"`
	Username string `yaml:"username" mapstructure:"username"`
	Password string `yaml:"password" mapstructure:"password"`
	Database string `yaml:"database" mapstructure:"database"`
	//数据库字符集 character set utf8mb4
	Charset string `yaml:"charset" mapstructure:"charset"`
	//连接池配置 connection pool configuration
	MaxIdleConns int `yaml:"max_idle_conns" mapstructure:"max_idle_conns"`
	//连接池最大打开连接数 maximum number of open connections
	MaxOpenConns int `yaml:"max_open_conns" mapstructure:"max_open_conns"`
	//连接池连接最大生命周期 maximum lifetime of a connection
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" mapstructure:"conn_max_lifetime"`
}

// RedisConfig Redis配置 redis configuration
type RedisConfig struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Port     int    `yaml:"port" mapstructure:"port"`
	Password string `yaml:"password" mapstructure:"password"`
	DB       int    `yaml:"db" mapstructure:"db"`
	//连接池配置 connection pool configuration
	PoolSize int `yaml:"pool_size" mapstructure:"pool_size"`
	//连接池最小空闲连接数 minimum number of idle connections
	MinIdleConns int           `yaml:"min_idle_conns" mapstructure:"min_idle_conns"`
	DialTimeout  time.Duration `yaml:"dial_timeout" mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout" mapstructure:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout" mapstructure:"write_timeout"`
}

// DetectConfig 检测相关配置
type DetectConfig struct {
	MaxWorkers     int           `yaml:"max_workers" mapstructure:"max_workers"`
	BatchSize      int           `yaml:"batch_size" mapstructure:"batch_size"`
	QueueSize      int           `yaml:"queue_size" mapstructure:"queue_size"`
	Timeout        time.Duration `yaml:"timeout" mapstructure:"timeout"`
	RetryCount     int           `yaml:"retry_count" mapstructure:"retry_count"`
	RetryDelay     time.Duration `yaml:"retry_delay" mapstructure:"retry_delay"`
	CppEngineURL   string        `yaml:"cpp_engine_url" mapstructure:"cpp_engine_url"`
	VideoStreamURL string        `yaml:"video_stream_url" mapstructure:"video_stream_url"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `yaml:"level" mapstructure:"level"`
	Format     string `yaml:"format" mapstructure:"format"`
	Output     string `yaml:"output" mapstructure:"output"`
	Filename   string `yaml:"filename" mapstructure:"filename"`
	MaxSize    int    `yaml:"max_size" mapstructure:"max_size"`
	MaxAge     int    `yaml:"max_age" mapstructure:"max_age"`
	MaxBackups int    `yaml:"max_backups" mapstructure:"max_backups"`
	Compress   bool   `yaml:"compress" mapstructure:"compress"`
}

// MetricsConfig 监控配置
type MetricsConfig struct {
	Enabled bool   `yaml:"enabled" mapstructure:"enabled"`
	Path    string `yaml:"path" mapstructure:"path"`
	Port    int    `yaml:"port" mapstructure:"port"`
}

// GetAddr 获取服务器地址
func (s *ServerConfig) GetAddr() string {
	// fmt.Sprintf() method is used to format a string
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// GetDSN 获取数据库DSN
func (d *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", d.Username, d.Password, d.Host, d.Port, d.Database, d.Charset)
}

// GetAddr 获取Redis地址
func (r *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

// GetAddr 获取监控服务地址
func (m *MetricsConfig) GetAddr() string {
	return fmt.Sprintf(":%d", m.Port)
}
