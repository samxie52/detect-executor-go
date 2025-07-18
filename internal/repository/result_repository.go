package repository

import (
	"context"
	"detect-executor-go/internal/model"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ResultRepository interface {
	BaseRepository

	Create(ctx context.Context, result *model.DetectResult) error
	GetByID(ctx context.Context, id uint) (*model.DetectResult, error)
	GetByResultID(ctx context.Context, resultID string) (*model.DetectResult, error)
	Update(ctx context.Context, result *model.DetectResult) error
	Delete(ctx context.Context, id uint) error

	GetByTaskID(ctx context.Context, taskID string) (*model.DetectResult, error)
	GetByDeviceID(ctx context.Context, deviceID string) (*model.DetectResult, error)
	List(ctx context.Context, deviceID string, detectType model.DetectType, page int, size int) ([]model.DetectResult, error)

	CountByDeviceID(ctx context.Context, deviceID string) (int64, error)
	GetResultStats(ctx context.Context, deviceID string, days int) (map[string]int64, error)
	GetQualityTrend(ctx context.Context, deviceID string, days int) ([]map[string]interface{}, error)
}

type resultRepository struct {
	BaseRepository
}

func NewResultRepository(db *gorm.DB, redis *redis.Client) ResultRepository {
	return &resultRepository{
		BaseRepository: NewBaseRepository(db, redis),
	}
}

// Implement ResultRepository methods
func (r *resultRepository) Create(ctx context.Context, result *model.DetectResult) error {
	return r.GetDb().WithContext(ctx).Create(result).Error
}

func (r *resultRepository) GetByID(ctx context.Context, id uint) (*model.DetectResult, error) {
	var result model.DetectResult
	err := r.GetDb().WithContext(ctx).First(&result, id).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *resultRepository) GetByResultID(ctx context.Context, resultID string) (*model.DetectResult, error) {
	var result model.DetectResult
	err := r.GetDb().WithContext(ctx).Where("result_id = ?", resultID).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *resultRepository) Update(ctx context.Context, result *model.DetectResult) error {
	return r.GetDb().WithContext(ctx).Save(result).Error
}

func (r *resultRepository) Delete(ctx context.Context, id uint) error {
	return r.GetDb().WithContext(ctx).Delete(&model.DetectResult{}, id).Error
}

func (r *resultRepository) GetByTaskID(ctx context.Context, taskID string) (*model.DetectResult, error) {
	var result model.DetectResult
	err := r.GetDb().WithContext(ctx).Where("task_id = ?", taskID).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *resultRepository) GetByDeviceID(ctx context.Context, deviceID string) (*model.DetectResult, error) {
	var result model.DetectResult
	err := r.GetDb().WithContext(ctx).Where("device_id = ?", deviceID).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *resultRepository) List(ctx context.Context, deviceID string, detectType model.DetectType, page int, size int) ([]model.DetectResult, error) {
	var results []model.DetectResult
	offset := (page - 1) * size
	query := r.GetDb().WithContext(ctx).Offset(offset).Limit(size)
	if deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}
	if detectType != "" {
		query = query.Where("detect_type = ?", detectType)
	}
	err := query.Find(&results).Error
	return results, err
}

func (r *resultRepository) CountByDeviceID(ctx context.Context, deviceID string) (int64, error) {
	var count int64
	err := r.GetDb().WithContext(ctx).Model(&model.DetectResult{}).Where("device_id = ?", deviceID).Count(&count).Error
	return count, err
}

func (r *resultRepository) GetResultStats(ctx context.Context, deviceID string, days int) (map[string]int64, error) {
	// Placeholder implementation - would need actual business logic
	stats := make(map[string]int64)
	return stats, nil
}

func (r *resultRepository) GetQualityTrend(ctx context.Context, deviceID string, days int) ([]map[string]interface{}, error) {
	// Placeholder implementation - would need actual business logic
	trend := make([]map[string]interface{}, 0)
	return trend, nil
}
