package client

import (
	"context"
	"detect-executor-go/internal/model"
	"encoding/json"
	"fmt"
)

// DetectEngineClient 检测引擎客户端接口
type DetectEngineClient interface {
	BaseClient
	SubmitTask(ctx context.Context, task *DetectTaskRequest) (*DetectTaskResponse, error)
	GetTaskStatus(ctx context.Context, taskID string) (*DetectTaskStatus, error)
	GetTaskResult(ctx context.Context, taskID string) (*DetectTaskResult, error)
	CancelTask(ctx context.Context, taskID string) error
}

// DetectTaskRequest 检测任务请求
type DetectTaskRequest struct {
	TaskID     string                 `json:"task_id"`
	Type       model.DetectType       `json:"type"`
	VideoURL   string                 `json:"video_url"`
	StartTime  string                 `json:"start_time"`
	EndTime    string                 `json:"end_time"`
	Parameters map[string]interface{} `json:"parameters"`
}

// DetectTaskResponse 检测任务响应
type DetectTaskResponse struct {
	TaskID  string `json:"task_id"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// DetectTaskStatus 检测任务状态
type DetectTaskStatus struct {
	TaskID   string `json:"task_id"`
	Status   string `json:"status"`
	Progress int    `json:"progress"`
	Message  string `json:"message,omitempty"`
}

// DetectTaskResult 检测任务结果
type DetectTaskResult struct {
	TaskID     string                 `json:"task_id"`
	Status     string                 `json:"status"`
	Score      float64                `json:"score"`
	Passed     bool                   `json:"passed"`
	Confidence float64                `json:"confidence"`
	Details    map[string]interface{} `json:"details"`
	ReportURL  string                 `json:"report_url"`
}

// detectEngineClient 检测引擎客户端实现
type detectEngineClient struct {
	httpClient HttpClient
	baseURL    string
	apiKey     string
}

// NewDetectEngineClient 创建检测引擎客户端
func NewDetectEngineClient(baseURL, apiKey string, config ClientConfig) DetectEngineClient {
	return &detectEngineClient{
		httpClient: NewHTTPClient(config),
		baseURL:    baseURL,
		apiKey:     apiKey,
	}

}

// SubmitTask 提交检测任务
func (c *detectEngineClient) SubmitTask(ctx context.Context, task *DetectTaskRequest) (*DetectTaskResponse, error) {
	url := fmt.Sprintf("%s/api/v1/detect/task", c.baseURL)
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", c.apiKey),
	}

	resp, err := c.httpClient.Post(ctx, url, task, headers)
	if err != nil {
		return nil, fmt.Errorf("submit task failed: %v", err)
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return nil, fmt.Errorf("submit task failed with status code %v and message %s", resp.StatusCode, string(resp.Body))
	}

	var result DetectTaskResponse
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}

	return &result, nil
}

// GetTaskStatus 获取检测任务状态
func (c *detectEngineClient) GetTaskStatus(ctx context.Context, taskID string) (*DetectTaskStatus, error) {
	url := fmt.Sprintf("%s/api/v1/detect/task/%s", c.baseURL, taskID)
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", c.apiKey),
	}

	resp, err := c.httpClient.Get(ctx, url, headers)
	if err != nil {
		return nil, fmt.Errorf("get task status failed: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("get task status failed with status code %v and message %s", resp.StatusCode, string(resp.Body))
	}

	var result DetectTaskStatus
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}

	return &result, nil
}

// GetTaskResult 获取检测任务结果
func (c *detectEngineClient) GetTaskResult(ctx context.Context, taskID string) (*DetectTaskResult, error) {
	url := fmt.Sprintf("%s/api/v1/detect/result/%s", c.baseURL, taskID)
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", c.apiKey),
	}

	resp, err := c.httpClient.Get(ctx, url, headers)
	if err != nil {
		return nil, fmt.Errorf("get task result failed: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("get task result failed with status code %v and message %s", resp.StatusCode, string(resp.Body))
	}

	var result DetectTaskResult
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}

	return &result, nil
}

// CancelTask 取消检测任务
func (c *detectEngineClient) CancelTask(ctx context.Context, taskID string) error {
	url := fmt.Sprintf("%s/api/v1/detect/task/cancel/%s", c.baseURL, taskID)
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", c.apiKey),
	}

	resp, err := c.httpClient.Delete(ctx, url, headers)
	if err != nil {
		return fmt.Errorf("cancel task failed: %v", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("cancel task failed with status code %v and message %s", resp.StatusCode, string(resp.Body))
	}

	return nil
}

// Close 关闭客户端
func (c *detectEngineClient) Close() error {
	return c.httpClient.Close()
}

// Ping 测试连接
func (c *detectEngineClient) Ping(ctx context.Context) error {
	url := fmt.Sprintf("%s/api/v1/health", c.baseURL)
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", c.apiKey),
	}

	resp, err := c.httpClient.Get(ctx, url, headers)
	if err != nil {
		return fmt.Errorf("ping failed: %v", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("ping failed with status code %v and message %s", resp.StatusCode, string(resp.Body))
	}

	return nil
}
