package service

import (
	"context"
	"detect-executor-go/internal/model"
	"detect-executor-go/pkg/client"
	"encoding/json"
	"fmt"
	"time"
)

type TaskService interface {
	BaseService

	CreateTask(ctx context.Context, req *CreateTaskRequest) (*model.DetectTask, error)
	GetTask(ctx context.Context, taskID string) (*model.DetectTask, error)
	ListTasks(ctx context.Context, req *model.TaskListRequest) ([]model.DetectTask, int64, error)
	UpdateTaskStatus(ctx context.Context, taskID string, status model.TaskStatus, message string) error
	CancelTask(ctx context.Context, taskID string) error
	RetryTask(ctx context.Context, taskID string) error
	SubmitTaskToEngine(ctx context.Context, task *model.DetectTask) error
	ProcessTaskResult(ctx context.Context, taskID string) error
}

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	DeviceID    string                 `json:"device_id" validate:"required"`
	Type        model.DetectType       `json:"type" validate:"required"`
	VideoURL    string                 `json:"video_url" validate:"required "`
	StartTime   *time.Time             `json:"start_time,omitempty"`
	EndTime     *time.Time             `json:"end_time,omitempty"`
	Priority    model.Priority         `json:"priority"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	Description string                 `json:"description,omitempty"`
}

// taskService implements TaskService
type taskService struct {
	ctx *ServiceContext
}

// NewTaskService 创建任务服务
func NewTaskService(ctx *ServiceContext) TaskService {
	return &taskService{ctx: ctx}
}

// CreateTask 创建任务
func (s *taskService) CreateTask(ctx context.Context, req *CreateTaskRequest) (*model.DetectTask, error) {
	device, err := s.ctx.Repository.Device.GetByDeviceID(ctx, req.DeviceID)
	if err != nil {
		return nil, fmt.Errorf("device not found: %w", err)
	}

	if device.Status != "active" {
		return nil, fmt.Errorf("device %s is not active", device.DeviceID)
	}
	task := &model.DetectTask{
		DeviceID:    req.DeviceID,
		Type:        req.Type,
		VideoURL:    req.VideoURL,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Priority:    req.Priority,
		Description: req.Description,
	}

	if err := s.ctx.Repository.Task.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("create task failed: %w", err)
	}

	s.ctx.Logger.WithFields(map[string]interface{}{
		"task_id":     task.TaskID,
		"device_id":   task.DeviceID,
		"type":        task.Type,
		"video_url":   task.VideoURL,
		"start_time":  task.StartTime,
		"end_time":    task.EndTime,
		"priority":    task.Priority,
		"description": task.Description,
	}).Info("create task success")

	return task, nil
}

func (s *taskService) GetTask(ctx context.Context, taskID string) (*model.DetectTask, error) {
	return s.ctx.Repository.Task.GetByTaskID(ctx, taskID)
}

func (s *taskService) ListTasks(ctx context.Context, req *model.TaskListRequest) ([]model.DetectTask, int64, error) {

	tasks, err := s.ctx.Repository.Task.List(ctx, req)
	if err != nil {
		return nil, 0, fmt.Errorf("list tasks failed: %w", err)
	}
	return tasks, 0, nil
}

func (s *taskService) UpdateTaskStatus(ctx context.Context, taskID string, status model.TaskStatus, message string) error {
	task, err := s.ctx.Repository.Task.GetByTaskID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("get task failed: %w", err)
	}
	task.Status = status
	task.UpdatedAt = time.Now()
	if err := s.ctx.Repository.Task.Update(ctx, task); err != nil {
		return fmt.Errorf("update task status failed: %w", err)
	}

	statusMsg := map[string]interface{}{
		"task_id":   taskID,
		"status":    status,
		"message":   message,
		"timestamp": time.Now(),
	}

	if err := s.ctx.Client.MessageQueue.Publish(ctx, "task.status", "task.status.changed", statusMsg); err != nil {
		s.ctx.Logger.WithFields(map[string]interface{}{
			"task_id": taskID,
			"status":  status,
			"error":   err,
		}).Warn("publish task status failed")
	}

	s.ctx.Logger.WithFields(map[string]interface{}{
		"task_id": taskID,
		"status":  status,
		"message": message,
	}).Info("update task status success")

	return nil
}

func (s *taskService) CancelTask(ctx context.Context, taskID string) error {
	task, err := s.ctx.Repository.Task.GetByTaskID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("get task failed: %w", err)
	}
	if task.Status != model.TaskStatusPending {
		return fmt.Errorf("task %s is not pending", taskID)
	}
	if err := s.ctx.Client.DetectEngine.CancelTask(ctx, taskID); err != nil {
		return fmt.Errorf("cancel task failed: %w", err)
	}
	return s.UpdateTaskStatus(ctx, taskID, model.TaskStatusCancel, "task cancelled")
}

func (s *taskService) RetryTask(ctx context.Context, taskID string) error {
	task, err := s.ctx.Repository.Task.GetByTaskID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("get task failed: %w", err)
	}
	if task.Status != model.TaskStatusFailed {
		return fmt.Errorf("task %s is not failed", taskID)
	}
	if err := s.UpdateTaskStatus(ctx, taskID, model.TaskStatusPending, "task retry"); err != nil {
		return fmt.Errorf("retry task failed: %w", err)
	}
	return s.SubmitTaskToEngine(ctx, task)
}

func (s *taskService) SubmitTaskToEngine(ctx context.Context, task *model.DetectTask) error {
	engineTask := &client.DetectTaskRequest{
		TaskID:    task.TaskID,
		Type:      task.Type,
		VideoURL:  task.VideoURL,
		StartTime: task.StartTime,
		EndTime:   task.EndTime,
	}

	_, err := s.ctx.Client.DetectEngine.SubmitTask(ctx, engineTask)
	if err != nil {
		if updateErr := s.UpdateTaskStatus(ctx, task.TaskID, model.TaskStatusFailed, "submit task failed"); updateErr != nil {
			return fmt.Errorf("submit task and update task status failed: %w", updateErr)
		}
		return fmt.Errorf("submit task failed: %w", err)
	}
	if err := s.UpdateTaskStatus(ctx, task.TaskID, model.TaskStatusRunning, "task submitted to engine"); err != nil {
		return fmt.Errorf("submit task and update task status failed: %w", err)
	}
	return nil
}

func (s *taskService) ProcessTaskResult(ctx context.Context, taskID string) error {
	result, err := s.ctx.Client.DetectEngine.GetTaskResult(ctx, taskID)
	if err != nil {
		return fmt.Errorf("get task result failed: %w", err)
	}

	detectResult := &model.DetectResult{
		TaskID:     taskID,
		Type:       model.DetectType(result.Status),
		Score:      result.Score,
		Passed:     result.Passed,
		Confidence: result.Confidence,
		Details:    json.RawMessage(result.Details),
		ReportURL:  result.ReportURL,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.ctx.Repository.Result.Create(ctx, detectResult); err != nil {
		return fmt.Errorf("create detect result failed: %w", err)
	}

	status := model.TaskStatusCompleted
	message := "task completed"
	if !result.Passed {
		message = "process task failed"
	}
	if err := s.UpdateTaskStatus(ctx, taskID, status, message); err != nil {
		return fmt.Errorf("update task status failed: %w", err)
	}
	s.ctx.Logger.WithFields(map[string]interface{}{
		"task_id": taskID,
		"status":  status,
		"message": message,
	}).Info("update task status success")
	return nil
}

// Close
func (s *taskService) Close() error {
	return nil
}
