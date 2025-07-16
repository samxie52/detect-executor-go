package model

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

//validator 是一个验证器，用于验证结构体字段

func TestDetectTaskModel(t *testing.T) {
	task := &DetectTask{
		DeviceID: "test-device-001",
		TaskID:   "test-task-001",
		Status:   TaskStatusPending,
		Type:     DetectTypeQuality,
		VideoURL: "rtsp://test_device_id/stream",
		Priority: PriorityMedium,
		Duration: 10,
	}

	//validator.New() 创建一个验证器
	validator := validator.New()
	//validator.Struct(task) 验证 task 结构体
	err := validator.Struct(task)
	assert.NoError(t, err)

	//assert.False() 表示断言失败
	assert.False(t, task.IsCompleted())
	assert.False(t, task.CanRetry())

}

func TestDetectResultModel(t *testing.T) {
	result := &DetectResult{
		TaskID:       "test-task-001",
		DeviceID:     "test-device-001",
		ResultID:     "test-result-001",
		Type:         DetectTypeQuality,
		Score:        95.0,
		Passed:       true,
		Confidence:   0.9,
		Details:      json.RawMessage(`{"brightness": 100}`),
		Metadata:     json.RawMessage(`{"device_name": "test-device-001"}`),
		ReportURL:    "http://test.com/report.pdf",
		ThumbnailURL: "http://test.com/thumbnail.jpg",
		ProcessTime:  100,
		FileSize:     1024,
	}

	validator := validator.New()
	err := validator.Struct(result)
	assert.NoError(t, err)
}

func TestQualityResult(t *testing.T) {
	quality := QualityResult{
		Brightness: 0.70,
		Contrast:   0.82,
		Sharpness:  0.94,
		Noise:      0.12,
		Blur:       0.24,
		Saturation: 0.36,
		Issues:     []string{"low_noise", "low_saturation"},
	}

	//json.Marshal() 将结构体转换为 JSON
	data, err := json.Marshal(quality)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "brightness")
}

func TestDeviceModel(t *testing.T) {
	device := &Device{
		DeviceID: "test-device-001",
		Name:     "test-device-001",
		Type:     "camera",
		Location: "test-location",
		IP:       "127.0.0.1",
		Port:     554,
		Protocol: "rtsp",
		Username: "test-username",
		Password: "test-password",
		Status:   "active",
		LastSeen: &time.Time{},
		Config:   json.RawMessage(`{"device_name": "test-device-001"}`),
	}

	validator := validator.New()
	err := validator.Struct(device)
	assert.NoError(t, err)

	streamURL := device.GetStreamURL()
	assert.Contains(t, streamURL, "rtsp://test-username:test-password@127.0.0.1:554/stream")
}

func TestWorkerModel(t *testing.T) {
	worker := &Worker{
		WorkerID:    "test-worker-001",
		Name:        "test-worker-name",
		IP:          "127.0.0.1",
		Port:        8080,
		Status:      "idle",
		LastSeen:    &time.Time{},
		ActiveTasks: 0,
		MaxTasks:    10,
		Config:      json.RawMessage(`{"worker_name": "test-worker-name"}`),
	}

	validator := validator.New()
	err := validator.Struct(worker)
	assert.NoError(t, err)

	assert.False(t, worker.IsOnline())
	assert.False(t, worker.CanAcceptTask())
	assert.Equal(t, worker.GetAddress(), "127.0.0.1:8080")
}
