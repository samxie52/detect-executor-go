package service

import (
	"context"
	"detect-executor-go/internal/repository"
	"detect-executor-go/pkg/client"
	"detect-executor-go/pkg/logger"
	"fmt"
	"sync"
)

// Manager 服务管理器
type Manager struct {
	mu sync.RWMutex

	Task   TaskService
	Detect DetectService
	Result ResultService

	ctx *ServiceContext
}

func NewManager(logger *logger.Logger, repository *repository.Manager, client *client.Manager) (*Manager, error) {
	ctx := NewServiceContext(logger, repository, client)
	manager := &Manager{
		ctx: ctx,
	}
	manager.Task = NewTaskService(ctx)
	manager.Detect = NewDetectService(ctx)
	manager.Result = NewResultService(ctx)

	if err := manager.Detect.MonitorTasks(context.Background()); err != nil {
		return nil, fmt.Errorf("manager monitor tasks failed: %w", err)
	}

	return manager, nil
}

// Close 关闭服务管理器
func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errors []error

	if m.Task != nil {
		if err := m.Task.Close(); err != nil {
			errors = append(errors, fmt.Errorf("manager close task failed: %w", err))
		}
	}
	if m.Detect != nil {
		if err := m.Detect.Close(); err != nil {
			errors = append(errors, fmt.Errorf("manager close detect failed: %w", err))
		}
	}
	if m.Result != nil {
		if err := m.Result.Close(); err != nil {
			errors = append(errors, fmt.Errorf("manager close result failed: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("manager close failed: %v", errors)
	}
	return nil
}
