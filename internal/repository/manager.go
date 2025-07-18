package repository

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Manager struct {
	db    *gorm.DB
	redis *redis.Client

	Task   TaskRepository
	Result ResultRepository
	Device DeviceRepository
	Worker WorkRepository
}

func NewManager(db *gorm.DB, redis *redis.Client) *Manager {
	return &Manager{
		db:     db,
		redis:  redis,
		Task:   NewTaskRepository(db, redis),
		Result: NewResultRepository(db, redis),
		Device: NewDeviceRepository(db, redis),
		Worker: NewWorkRepository(db, redis),
	}
}

func (m *Manager) GetDb() *gorm.DB {
	return m.db
}

func (m *Manager) GetRedis() *redis.Client {
	return m.redis
}

// Transaction
func (m *Manager) Transaction(fn func(*Manager) error) error {
	return m.db.Transaction(func(tx *gorm.DB) error {

		txManager := &Manager{
			db:    tx,
			redis: m.redis,

			Task:   NewTaskRepository(tx, m.redis),
			Result: NewResultRepository(tx, m.redis),
			Device: NewDeviceRepository(tx, m.redis),
			Worker: NewWorkRepository(tx, m.redis),
		}

		return fn(txManager)
	})
}
