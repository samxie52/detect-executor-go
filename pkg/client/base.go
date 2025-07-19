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
	Timeout         time.Duration
	RetryCount      int
	RetryInterval   time.Duration
	MaxIdleConns    int
	MaxConnsPerHost int
}

// DefaultClientConfig 默认客户端配置
var DefaultClientConfig = ClientConfig{
	Timeout:         time.Second * 30,
	RetryCount:      3,
	RetryInterval:   time.Second * 1,
	MaxIdleConns:    100,
	MaxConnsPerHost: 10,
}
