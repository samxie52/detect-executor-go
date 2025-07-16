package model

import (
	"encoding/json"
	"time"
)

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	DeviceID   string                 `json:"device_id" validate:"required"`
	Type       DetectType             `json:"type" validate:"required,oneof=quality integrity motion object"`
	VideoURL   string                 `json:"video_url" validate:"required,url"`
	StartTime  time.Time              `json:"start_time"`
	EndTime    time.Time              `json:"end_time"`
	Duration   int                    `json:"duration" validate:"min=1"`
	Priority   Priority               `json:"priority" validate:"required,min=1,max=4"`
	Parameters map[string]interface{} `json:"parameters"`
}

// TaskResponse 任务响应
type TaskResponse struct {
	TaskID      string    `json:"task_id"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	StartedAt   time.Time `json:"started_at"`
	CompletedAt time.Time `json:"completed_at"`

	Progress int    `json:"progress"` //1-100
	Message  string `json:"message,omitempty"`
}

// TaskListRequest 任务列表请求
type TaskListRequest struct {
	DeviceID  string     `json:"device_id"`
	Type      DetectType `json:"type"`
	Status    TaskStatus `json:"status"`
	Page      int        `json:"page"`
	Size      int        `json:"size"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
}

// TaskListResponse 任务列表响应
type TaskListResponse struct {
	Tasks      []DetectTask `json:"tasks"`
	Total      int          `json:"total"`
	Page       int          `json:"page"`
	Size       int          `json:"size"`
	TotalPages int          `json:"total_pages"`
}

// ResultResponse 检测结果响应
type ResultResponse struct {
	ResultID     string          `json:"result_id"`
	TaskID       string          `json:"task_id"`
	Type         DetectType      `json:"type"`
	Score        float64         `json:"score"`
	Passed       bool            `json:"passed"`
	Confidence   float64         `json:"confidence"`
	Details      json.RawMessage `json:"details"`
	ReportURL    string          `json:"report_url"`
	ThumbnailURL string          `json:"thumbnail_url"`
	CreatedAt    time.Time       `json:"created_at"`
}

// HealthCheckResponse 健康检查响应
type HealthCheckResponse struct {
	Status    string                 `json:"status"`
	Timestamp int64                  `json:"timestamp"`
	Service   map[string]interface{} `json:"service"`
	Version   string                 `json:"version"`
}
