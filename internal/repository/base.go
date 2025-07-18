package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// BaseRepository 是一个接口，用于定义数据库和 Redis 的操作
type BaseRepository interface {
	GetDb() *gorm.DB
	GetRedis() *redis.Client
	// WithContext 是一个链式调用，用于设置上下文
	// WithContext is a chainable call, used to set the context
	//context 是 gorm 的上下文，用于在 gorm 操作中传递上下文信息
	//context is gorm's context, used to pass context information in gorm operations
	//context.Context 是 golang 的上下文，用于在 golang 操作中传递上下文信息
	//context.Context is golang's context, used to pass context information in golang operations
	WithContext(ctx context.Context) BaseRepository
	// WithTx 是一个链式调用，用于设置事务
	// WithTx is a chainable call, used to set the transaction
	//tx 是 gorm 的事务，用于在 gorm 操作中传递事务信息
	//tx is gorm's transaction, used to pass transaction information in gorm operations
	WithTx(tx *gorm.DB) BaseRepository
}

// baseRepository 是 BaseRepository 的实现
type baseRepository struct {
	db    *gorm.DB
	redis *redis.Client
	ctx   context.Context
}

// NewBaseRepository 创建一个 BaseRepository
func NewBaseRepository(db *gorm.DB, redis *redis.Client) BaseRepository {
	return &baseRepository{
		db:    db,
		redis: redis,
		ctx:   context.Background(),
	}
}

func (b *baseRepository) GetDb() *gorm.DB {
	return b.db.WithContext(b.ctx)
}

func (b *baseRepository) GetRedis() *redis.Client {
	return b.redis
}

func (b *baseRepository) WithContext(ctx context.Context) BaseRepository {
	return &baseRepository{
		db:    b.db.WithContext(ctx),
		redis: b.redis,
		ctx:   ctx,
	}
}

func (b *baseRepository) WithTx(tx *gorm.DB) BaseRepository {
	return &baseRepository{
		db:    tx,
		redis: b.redis,
		ctx:   b.ctx,
	}
}

// CacheConfig 是一个结构体，用于定义缓存的配置
// CacheConfig is a struct, used to define the cache configuration
type CacheConfig struct {
	TTL    time.Duration
	Prefix string
}

// DefaultCacheConfig 是默认的缓存配置, 30分钟过期, 前缀为 detect:
// DefaultCacheConfig is the default cache configuration, 30 minutes expiration, prefix is detect:
var DefaultCacheConfig = CacheConfig{
	TTL:    time.Minute * 30,
	Prefix: "detect:",
}
