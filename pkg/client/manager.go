package client

import (
	"context"
	"detect-executor-go/pkg/config"
	"fmt"
	"sync"
)

// Manager 客户端管理器
type Manager struct {
	//mu 用于保护 HTTP、DetectEngine、Storage 和 MessageQueue 的访问
	mu sync.RWMutex

	HTTP         HttpClient
	DetectEngine DetectEngineClient
	Storage      StorageClient
	MessageQueue MessageQueueClient

	config *Config
}

// Config 客户端配置
type Config struct {
	HTTP         config.HTTPClientConfig
	DetectEngine config.DetectEngineConfig
	Storage      config.StorageConfig
	MessageQueue config.MessageQueueConfig
}

func NewManager(config *Config) (*Manager, error) {
	manager := &Manager{
		config: config,
	}
	// init http
	manager.HTTP = NewHTTPClient(config.HTTP)

	// init detect engine
	manager.DetectEngine = NewDetectEngineClient(
		config.DetectEngine.BaseURL,
		config.DetectEngine.APIKey,
		config.DetectEngine.Client,
	)

	// init storage
	storageClient, err := NewStorageClient(config.Storage)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage client: %v", err)
	}
	manager.Storage = storageClient

	// init message queue
	mqClient, err := NewMessageQueueClient(config.MessageQueue)
	if err != nil {
		return nil, fmt.Errorf("failed to create message queue client: %v", err)
	}
	manager.MessageQueue = mqClient

	return manager, nil

}

// Close 关闭客户端
func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	var errors []error

	if m.HTTP != nil {
		if err := m.HTTP.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close HTTP client: %v", err))
		}
	}

	if m.DetectEngine != nil {
		if err := m.DetectEngine.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close detect engine client: %v", err))
		}
	}

	if m.Storage != nil {
		if err := m.Storage.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close storage client: %v", err))
		}
	}

	if m.MessageQueue != nil {
		if err := m.MessageQueue.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close message queue client: %v", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to close clients: %v", errors)
	}

	return nil
}

// HealthCheck 健康检查
func (m *Manager) HealthCheck(ctx context.Context) map[string]error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make(map[string]error)

	if m.HTTP != nil {
		if err := m.HTTP.Ping(ctx); err != nil {
			result["http"] = err
		}
	}

	if m.DetectEngine != nil {
		if err := m.DetectEngine.Ping(ctx); err != nil {
			result["detect_engine"] = err
		}
	}

	if m.Storage != nil {
		if err := m.Storage.Ping(ctx); err != nil {
			result["storage"] = err
		}
	}

	if m.MessageQueue != nil {
		if err := m.MessageQueue.Ping(ctx); err != nil {
			result["message_queue"] = err
		}
	}

	return result
}
