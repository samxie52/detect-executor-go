package repository

import (
	"context"
	"detect-executor-go/internal/model"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type WorkRepository interface {
	BaseRepository

	Create(ctx context.Context, worker *model.Worker) error
	GetByID(ctx context.Context, id uint) (*model.Worker, error)
	GetByWorkerID(ctx context.Context, workerID string) (*model.Worker, error)
	Update(ctx context.Context, worker *model.Worker) error
	Delete(ctx context.Context, id uint) error

	List(ctx context.Context, page int, size int) ([]model.Worker, error)
	GetOnlineWorkers(ctx context.Context) ([]model.Worker, error)
	GetAvailableWorkers(ctx context.Context) ([]model.Worker, error)
	GetByStatus(ctx context.Context, status string) ([]model.Worker, error)

	UpdateLastSeen(ctx context.Context, workerID string) error
	UpdateStatus(ctx context.Context, workerID string, status string) error
	UpdateStats(ctx context.Context, workerID string, stats map[string]interface{}) error
}

type workRepository struct {
	BaseRepository
}

func NewWorkRepository(db *gorm.DB, redis *redis.Client) WorkRepository {
	return &workRepository{
		BaseRepository: NewBaseRepository(db, redis),
	}
}

// Implement WorkRepository methods
func (r *workRepository) Create(ctx context.Context, worker *model.Worker) error {
	return r.GetDb().WithContext(ctx).Create(worker).Error
}

func (r *workRepository) GetByID(ctx context.Context, id uint) (*model.Worker, error) {
	var worker model.Worker
	err := r.GetDb().WithContext(ctx).First(&worker, id).Error
	if err != nil {
		return nil, err
	}
	return &worker, nil
}

func (r *workRepository) GetByWorkerID(ctx context.Context, workerID string) (*model.Worker, error) {
	var worker model.Worker
	err := r.GetDb().WithContext(ctx).Where("worker_id = ?", workerID).First(&worker).Error
	if err != nil {
		return nil, err
	}
	return &worker, nil
}

func (r *workRepository) Update(ctx context.Context, worker *model.Worker) error {
	return r.GetDb().WithContext(ctx).Save(worker).Error
}

func (r *workRepository) Delete(ctx context.Context, id uint) error {
	return r.GetDb().WithContext(ctx).Delete(&model.Worker{}, id).Error
}

func (r *workRepository) List(ctx context.Context, page int, size int) ([]model.Worker, error) {
	var workers []model.Worker
	offset := (page - 1) * size
	err := r.GetDb().WithContext(ctx).Offset(offset).Limit(size).Find(&workers).Error
	return workers, err
}

func (r *workRepository) GetOnlineWorkers(ctx context.Context) ([]model.Worker, error) {
	var workers []model.Worker
	err := r.GetDb().WithContext(ctx).Where("status = ?", "online").Find(&workers).Error
	return workers, err
}

func (r *workRepository) GetAvailableWorkers(ctx context.Context) ([]model.Worker, error) {
	var workers []model.Worker
	err := r.GetDb().WithContext(ctx).Where("status = ?", "available").Find(&workers).Error
	return workers, err
}

func (r *workRepository) GetByStatus(ctx context.Context, status string) ([]model.Worker, error) {
	var workers []model.Worker
	err := r.GetDb().WithContext(ctx).Where("status = ?", status).Find(&workers).Error
	return workers, err
}

func (r *workRepository) UpdateLastSeen(ctx context.Context, workerID string) error {
	return r.GetDb().WithContext(ctx).Model(&model.Worker{}).Where("worker_id = ?", workerID).Update("last_seen", "NOW()").Error
}

func (r *workRepository) UpdateStatus(ctx context.Context, workerID string, status string) error {
	return r.GetDb().WithContext(ctx).Model(&model.Worker{}).Where("worker_id = ?", workerID).Update("status", status).Error
}

func (r *workRepository) UpdateStats(ctx context.Context, workerID string, stats map[string]interface{}) error {
	// Placeholder implementation - would need actual business logic
	return r.GetDb().WithContext(ctx).Model(&model.Worker{}).Where("worker_id = ?", workerID).Update("stats", stats).Error
}
