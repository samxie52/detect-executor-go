package repository

import (
	"context"
	"detect-executor-go/internal/model"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type TaskRepository interface {
	BaseRepository

	// why need ctx?
	// ctx is used to pass context information in gorm operations
	// base Create, Update, Delete methods
	Create(ctx context.Context, task *model.DetectTask) error
	GetByID(ctx context.Context, id uint) (*model.DetectTask, error)
	GetByTaskID(ctx context.Context, taskID string) (*model.DetectTask, error)
	Update(ctx context.Context, task *model.DetectTask) error
	Delete(ctx context.Context, id uint) error

	//query
	List(ctx context.Context, req *model.TaskListRequest) ([]model.DetectTask, error)
	GetByDeviceId(ctx context.Context, deviceID string, limit int) ([]model.DetectTask, error)
	GetByStatus(ctx context.Context, status model.TaskStatus, limit int) ([]model.DetectTask, error)
	GetPendingTasks(ctx context.Context, limit int) ([]model.DetectTask, error)
	GetRetryTasks(ctx context.Context, limit int) ([]model.DetectTask, error)

	//update status
	UpdateStatus(ctx context.Context, taskID string, status model.TaskStatus) error
	UpdateWorker(ctx context.Context, taskID string, workerID string) error
	IncrementRetry(ctx context.Context, taskID string) error

	//query count
	CountByStatus(ctx context.Context, status model.TaskStatus) (int64, error)
	CountByDeviceID(ctx context.Context, deviceID string) (int64, error)
	GetTaskStats(ctx context.Context, deviceID string, days int) (map[string]int64, error)
}

type taskRepository struct {
	BaseRepository
}

func NewTaskRepository(db *gorm.DB, redis *redis.Client) TaskRepository {
	return &taskRepository{
		BaseRepository: NewBaseRepository(db, redis),
	}
}

// Implement TaskRepository methods
func (r *taskRepository) Create(ctx context.Context, task *model.DetectTask) error {
	return r.GetDb().WithContext(ctx).Create(task).Error
}

func (r *taskRepository) GetByID(ctx context.Context, id uint) (*model.DetectTask, error) {
	var task model.DetectTask
	err := r.GetDb().WithContext(ctx).First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (t *taskRepository) GetByTaskID(ctx context.Context, taskID string) (*model.DetectTask, error) {
	var task model.DetectTask
	err := t.GetDb().WithContext(ctx).Where("task_id = ?", taskID).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) Update(ctx context.Context, task *model.DetectTask) error {
	return r.GetDb().WithContext(ctx).Save(task).Error
}

func (r *taskRepository) Delete(ctx context.Context, id uint) error {
	return r.GetDb().WithContext(ctx).Delete(&model.DetectTask{}, id).Error
}

func (r *taskRepository) List(ctx context.Context, req *model.TaskListRequest) ([]model.DetectTask, error) {
	var tasks []model.DetectTask
	// Placeholder implementation - would need actual business logic based on req
	err := r.GetDb().WithContext(ctx).Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) GetByDeviceId(ctx context.Context, deviceID string, limit int) ([]model.DetectTask, error) {
	var tasks []model.DetectTask
	err := r.GetDb().WithContext(ctx).Where("device_id = ?", deviceID).Limit(limit).Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) GetByStatus(ctx context.Context, status model.TaskStatus, limit int) ([]model.DetectTask, error) {
	var tasks []model.DetectTask
	err := r.GetDb().WithContext(ctx).Where("status = ?", status).Limit(limit).Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) GetPendingTasks(ctx context.Context, limit int) ([]model.DetectTask, error) {
	var tasks []model.DetectTask
	err := r.GetDb().WithContext(ctx).Where("status = ?", "pending").Limit(limit).Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) GetRetryTasks(ctx context.Context, limit int) ([]model.DetectTask, error) {
	var tasks []model.DetectTask
	err := r.GetDb().WithContext(ctx).Where("status = ?", "retry").Limit(limit).Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) UpdateStatus(ctx context.Context, taskID string, status model.TaskStatus) error {
	return r.GetDb().WithContext(ctx).Model(&model.DetectTask{}).Where("task_id = ?", taskID).Update("status", status).Error
}

func (r *taskRepository) UpdateWorker(ctx context.Context, taskID string, workerID string) error {
	return r.GetDb().WithContext(ctx).Model(&model.DetectTask{}).Where("task_id = ?", taskID).Update("worker_id", workerID).Error
}

func (r *taskRepository) IncrementRetry(ctx context.Context, taskID string) error {
	return r.GetDb().WithContext(ctx).Model(&model.DetectTask{}).Where("task_id = ?", taskID).Update("retry_count", "retry_count + 1").Error
}

func (r *taskRepository) CountByStatus(ctx context.Context, status model.TaskStatus) (int64, error) {
	var count int64
	err := r.GetDb().WithContext(ctx).Model(&model.DetectTask{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

func (r *taskRepository) CountByDeviceID(ctx context.Context, deviceID string) (int64, error) {
	var count int64
	err := r.GetDb().WithContext(ctx).Model(&model.DetectTask{}).Where("device_id = ?", deviceID).Count(&count).Error
	return count, err
}

func (r *taskRepository) GetTaskStats(ctx context.Context, deviceID string, days int) (map[string]int64, error) {
	// Placeholder implementation - would need actual business logic
	stats := make(map[string]int64)
	return stats, nil
}
