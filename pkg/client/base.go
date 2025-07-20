package client

import (
	"context"
)

// BaseClient 基础客户端接口
type BaseClient interface {
	Close() error
	Ping(ctx context.Context) error
}
