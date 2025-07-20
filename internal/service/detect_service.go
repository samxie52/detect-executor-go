package service

import (
	"context"
	"detect-executor-go/internal/model"
	"fmt"
	"sync"
	"time"
)

type DetectService interface {
	BaseService
	StartDetection(ctx context.Context, taskID string) error
	StopDetection(ctx context.Context, taskID string) error
	GetDetectionStatus(ctx context.Context, taskID string) (*DetectionStatus, error)
	MonitorTasks(ctx context.Context) error
}

// DetectionStatus 检测状态
type DetectionStatus struct {
	TaskID    string    `json:"task_id"`
	Status    string    `json:"status"`
	Progress  int       `json:"progress"`
	StartTime time.Time `json:"start_time"`
	// 任务执行时间
	ElapsedTime time.Duration `json:"elapsed_time"`
	Message     string        `json:"message,omitempty"`
}

// detectService 检测服务实现
type detectService struct {
	ctx          *ServiceContext
	runningTasks map[string]*DetectionStatus
	// 任务状态锁, 用于保护 runningTasks
	tasksMutex sync.RWMutex
	// 停止检测的通道, 用于通知检测任务停止, struct{} 表示无值
	stopChan chan struct{}
	// 监控检测任务的定时器
	monitorTicker *time.Ticker
}

// NewDetectService 创建检测服务
func NewDetectService(ctx *ServiceContext) DetectService {
	return &detectService{
		ctx:          ctx,
		runningTasks: make(map[string]*DetectionStatus),
		stopChan:     make(chan struct{}),
	}
}

// StartDetection 开始检测
func (s *detectService) StartDetection(ctx context.Context, taskID string) error {
	s.tasksMutex.Lock()
	defer s.tasksMutex.Unlock()

	if _, exists := s.runningTasks[taskID]; exists {
		return fmt.Errorf("task %s is already running", taskID)
	}

	task, err := s.ctx.Repository.Task.GetByTaskID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("get task failed: %w", err)
	}

	status := &DetectionStatus{
		TaskID:    taskID,
		Status:    string(model.TaskStatusStarting),
		Progress:  0,
		StartTime: time.Now(),
		Message:   "Detection starting",
	}

	s.runningTasks[taskID] = status

	go s.executeDetection(taskID, task)

	s.ctx.Logger.WithFields(map[string]interface{}{
		"task_id": taskID,
		"status":  status,
	}).Info("start detection")

	return nil
}

// StopDetection 停止检测
func (s *detectService) StopDetection(ctx context.Context, taskID string) error {
	s.tasksMutex.Lock()
	defer s.tasksMutex.Unlock()

	status, exists := s.runningTasks[taskID]
	if !exists {
		return fmt.Errorf("task %s is not running", taskID)
	}

	status.Status = string(model.TaskStatusStopping)
	status.Message = "Detection stopping"

	if err := s.ctx.Client.DetectEngine.CancelTask(ctx, taskID); err != nil {
		s.ctx.Logger.WithFields(map[string]interface{}{
			"task_id": taskID,
			"error":   err,
		}).Warn("cancel task failed")

	}

	delete(s.runningTasks, taskID)

	s.ctx.Logger.WithFields(map[string]interface{}{
		"task_id": taskID,
		"status":  status,
	}).Info("stop detection")

	return nil
}

// GetDetectionStatus 获取检测状态
func (s *detectService) GetDetectionStatus(ctx context.Context, taskID string) (*DetectionStatus, error) {
	s.tasksMutex.Lock()
	defer s.tasksMutex.Unlock()

	status, exists := s.runningTasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task %s is not running", taskID)
	}

	status.ElapsedTime = time.Since(status.StartTime)

	return status, nil
}

// MonitorTasks 监控检测任务
func (s *detectService) MonitorTasks(ctx context.Context) error {
	// 每5秒检查一次任务状态
	s.monitorTicker = time.NewTicker(5 * time.Second)

	go func() {
		for {
			select {
			// s.stopChan 用于通知检测任务停止
			case <-s.stopChan:
				return
			// s.monitorTicker.C 每5秒检查一次任务状态
			case <-s.monitorTicker.C:
				s.checkRunningTasks(ctx)
			}
		}
	}()

	s.ctx.Logger.Info("monitor tasks started")

	return nil
}

// executeDetection 执行检测
func (s *detectService) executeDetection(taskID string, task *model.DetectTask) {
	ctx := context.Background()

	s.updateTaskStatus(taskID, string(model.TaskStatusRunning), 10, "Detection in progress")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			engineStatus, err := s.ctx.Client.DetectEngine.GetTaskStatus(ctx, taskID)
			if err != nil {
				s.ctx.Logger.WithFields(map[string]interface{}{
					"task_id": taskID,
					"error":   err.Error(),
				}).Error("get task status from engine failed")
				s.finishDetection(taskID, "failed", "Failed to get task status from engine")
				return
			}

			// 更新任务状态
			s.updateTaskStatus(taskID, engineStatus.Status, engineStatus.Progress, engineStatus.Message)

			// 任务完成或失败
			if engineStatus.Status == "completed" || engineStatus.Status == "failed" {
				s.finishDetection(taskID, engineStatus.Status, engineStatus.Message)
				return
			}

		}
	}
}

// updateTaskStatus 更新任务状态
func (s *detectService) updateTaskStatus(taskID string, status string, progress int, message string) {
	s.tasksMutex.Lock()
	defer s.tasksMutex.Unlock()

	if taskStatus, exists := s.runningTasks[taskID]; exists {
		taskStatus.Status = status
		taskStatus.Progress = progress
		taskStatus.Message = message
	}
}

// finishDetection 完成检测
func (s *detectService) finishDetection(taskID string, status string, message string) {
	s.tasksMutex.Lock()
	defer s.tasksMutex.Unlock()

	delete(s.runningTasks, taskID)

	s.ctx.Logger.WithFields(map[string]interface{}{
		"task_id": taskID,
		"status":  status,
		"message": message,
	}).Info("Detection finished")
}

// checkRunningTasks 检查运行中的任务
func (s *detectService) checkRunningTasks(ctx context.Context) {
	s.tasksMutex.RLock()

	taskIDs := make([]string, 0, len(s.runningTasks))
	for taskID := range s.runningTasks {
		taskIDs = append(taskIDs, taskID)
	}
	s.tasksMutex.RUnlock()

	for _, taskID := range taskIDs {
		// 检查任务是否超时
		status, exists := s.runningTasks[taskID]
		if exists && time.Since(status.StartTime) > 30*time.Minute {
			s.ctx.Logger.WithFields(map[string]interface{}{
				"task_id": taskID,
				"status":  status,
			}).Warn("Detection timeout")

			s.StopDetection(ctx, taskID)
		}
	}

}

// Close 关闭检测服务
func (s *detectService) Close() error {
	if s.monitorTicker != nil {
		s.monitorTicker.Stop()
	}
	close(s.stopChan)
	return nil
}
