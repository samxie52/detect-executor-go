package service

import (
	"detect-executor-go/internal/repository"
	"detect-executor-go/pkg/client"
	"detect-executor-go/pkg/logger"
)

// BaseService 基础服务接口
type BaseService interface {
	Close() error
}

// ServiceContext 服务上下文
type ServiceContext struct {
	Logger     *logger.Logger
	Repository *repository.Manager
	Client     *client.Manager
}

func NewServiceContext(logger *logger.Logger, repository *repository.Manager, client *client.Manager) *ServiceContext {
	return &ServiceContext{
		Logger:     logger,
		Repository: repository,
		Client:     client,
	}
}
