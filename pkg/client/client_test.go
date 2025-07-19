package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// ClientTestSuite is a test suite for client
type ClientTestSuite struct {
	suite.Suite
	ctx context.Context
}

func (suite *ClientTestSuite) SetupSuite() {
	//context.Background() 返回一个空的 Context，通常用于 main 函数、初始化和测试
	suite.ctx = context.Background()
}

func (suite *ClientTestSuite) TestHTTPClient() {
	//创建一个测试 HTTP 服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	}))
	defer server.Close()

	client := NewHTTPClient(DefaultClientConfig)

	resp, err := client.Get(suite.ctx, server.URL, nil)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.Contains(string(resp.Body), "success")

	resp, err = client.Post(suite.ctx, server.URL, nil, nil)
	suite.NoError(err)
	suite.Equal(http.StatusOK, resp.StatusCode)

}

func (suite *ClientTestSuite) TestDetectEngineClient() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/detect/task":
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"task_id": "123","status": "pending"}`))
		case "/api/v1/health":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "ok"}`))
		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"status": "not found"}`))
		}
	}))
	defer server.Close()

	client := NewDetectEngineClient(server.URL, "test-key", DefaultClientConfig)

	err := client.Ping(suite.ctx)
	suite.NoError(err)

	task := &DetectTaskRequest{
		TaskID:   "123",
		Type:     "quality",
		VideoURL: "rtsp://test_device_id/stream",
	}

	resp, err := client.SubmitTask(suite.ctx, task)
	suite.NoError(err)
	suite.Equal("123", resp.TaskID)
	suite.Equal("pending", resp.Status)

}

func (suite *ClientTestSuite) TestClientManager() {
	config := &Config{
		HTTP: ClientConfig{
			Timeout:       10 * time.Second,
			RetryCount:    3,
			RetryInterval: 1 * time.Second,
		},
		DetectEngine: DetectEngineConfig{
			BaseURL: "http://localhost:8080",
			APIKey:  "test-key",
			Client: ClientConfig{
				Timeout: 10 * time.Second,
			},
		},
		Storage: StorageConfig{
			Endpoint:  "localhost:9000",
			AccessKey: "test-key",
			SecretKey: "test-key",
			UseSSL:    false,
		},
		MessageQueue: MessageQueueConfig{
			URL: "amqp://guest:guest@localhost:5672/",
		},
	}

	suite.NotNil(config.HTTP)
	suite.NotNil(config.DetectEngine)
	suite.NotNil(config.Storage)
	suite.NotNil(config.MessageQueue)
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}
