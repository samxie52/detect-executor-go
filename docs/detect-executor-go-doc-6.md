# 步骤 6: 数据模型定义

## 目标
定义检测相关的数据模型，包括检测任务、检测结果、设备信息等核心数据结构，并实现数据验证规则。

## 操作步骤

### 6.1 添加数据验证依赖
```bash
# 添加数据验证库
go get github.com/go-playground/validator/v10

# 添加UUID生成库
go get github.com/google/uuid

# 添加时间处理库
go get github.com/jinzhu/now
```

### 6.2 创建基础模型 (`internal/model/base.go`)
```go
package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TaskStatus 任务状态枚举
type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"    // 待处理
	TaskStatusRunning    TaskStatus = "running"    // 执行中
	TaskStatusCompleted  TaskStatus = "completed"  // 已完成
	TaskStatusFailed     TaskStatus = "failed"     // 失败
	TaskStatusCancelled  TaskStatus = "cancelled"  // 已取消
	TaskStatusTimeout    TaskStatus = "timeout"    // 超时
)

// DetectType 检测类型枚举
type DetectType string

const (
	DetectTypeQuality    DetectType = "quality"    // 视频质量检测
	DetectTypeIntegrity  DetectType = "integrity"  // 录像完整性检测
	DetectTypeMotion     DetectType = "motion"     // 运动检测
	DetectTypeObject     DetectType = "object"     // 目标检测
)

// Priority 优先级枚举
type Priority int

const (
	PriorityLow    Priority = 1 // 低优先级
	PriorityNormal Priority = 2 // 普通优先级
	PriorityHigh   Priority = 3 // 高优先级
	PriorityUrgent Priority = 4 // 紧急优先级
)
```

### 6.3 创建检测任务模型 (`internal/model/detect_task.go`)
```go
package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// DetectTask 检测任务模型
type DetectTask struct {
	BaseModel
	TaskID      string     `gorm:"uniqueIndex;size:36;not null" json:"task_id" validate:"required"`
	DeviceID    string     `gorm:"index;size:64;not null" json:"device_id" validate:"required"`
	DeviceName  string     `gorm:"size:128" json:"device_name"`
	Type        DetectType `gorm:"size:32;not null" json:"type" validate:"required,oneof=quality integrity motion object"`
	Status      TaskStatus `gorm:"size:32;not null;default:pending" json:"status"`
	Priority    Priority   `gorm:"not null;default:2" json:"priority" validate:"min=1,max=4"`
	
	// 检测参数
	VideoURL    string          `gorm:"size:512;not null" json:"video_url" validate:"required,url"`
	StartTime   *time.Time      `json:"start_time"`
	EndTime     *time.Time      `json:"end_time"`
	Duration    int             `json:"duration"` // 检测时长（秒）
	Parameters  json.RawMessage `gorm:"type:json" json:"parameters"`
	
	// 执行信息
	WorkerID    string     `gorm:"size:64" json:"worker_id"`
	StartedAt   *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
	RetryCount  int        `gorm:"default:0" json:"retry_count"`
	MaxRetries  int        `gorm:"default:3" json:"max_retries"`
	
	// 结果信息
	ResultID    string `gorm:"size:36" json:"result_id"`
	ErrorMsg    string `gorm:"size:1024" json:"error_msg"`
	
	// 关联关系
	Results []DetectResult `gorm:"foreignKey:TaskID;references:TaskID" json:"results,omitempty"`
}

// TableName 指定表名
func (DetectTask) TableName() string {
	return "detect_tasks"
}

// BeforeCreate 创建前钩子
func (dt *DetectTask) BeforeCreate(tx *gorm.DB) error {
	if dt.TaskID == "" {
		dt.TaskID = uuid.New().String()
	}
	return nil
}

// IsCompleted 检查任务是否完成
func (dt *DetectTask) IsCompleted() bool {
	return dt.Status == TaskStatusCompleted || 
		   dt.Status == TaskStatusFailed || 
		   dt.Status == TaskStatusCancelled ||
		   dt.Status == TaskStatusTimeout
}

// CanRetry 检查是否可以重试
func (dt *DetectTask) CanRetry() bool {
	return dt.Status == TaskStatusFailed && dt.RetryCount < dt.MaxRetries
}

// GetDuration 获取执行时长
func (dt *DetectTask) GetDuration() time.Duration {
	if dt.StartedAt == nil || dt.CompletedAt == nil {
		return 0
	}
	return dt.CompletedAt.Sub(*dt.StartedAt)
}
```

### 6.4 创建检测结果模型 (`internal/model/detect_result.go`)
```go
package model

import (
	"encoding/json"
)

// DetectResult 检测结果模型
type DetectResult struct {
	BaseModel
	ResultID    string          `gorm:"uniqueIndex;size:36;not null" json:"result_id" validate:"required"`
	TaskID      string          `gorm:"index;size:36;not null" json:"task_id" validate:"required"`
	DeviceID    string          `gorm:"index;size:64;not null" json:"device_id" validate:"required"`
	Type        DetectType      `gorm:"size:32;not null" json:"type"`
	
	// 检测结果
	Score       float64         `gorm:"type:decimal(5,4)" json:"score"` // 检测得分 0-1
	Passed      bool            `json:"passed"`                         // 是否通过检测
	Confidence  float64         `gorm:"type:decimal(5,4)" json:"confidence"` // 置信度
	
	// 详细结果数据
	Details     json.RawMessage `gorm:"type:json" json:"details"`
	Metadata    json.RawMessage `gorm:"type:json" json:"metadata"`
	
	// 文件信息
	ReportURL   string          `gorm:"size:512" json:"report_url"`
	ThumbnailURL string         `gorm:"size:512" json:"thumbnail_url"`
	
	// 统计信息
	ProcessTime int             `json:"process_time"` // 处理时间（毫秒）
	FileSize    int64           `json:"file_size"`    // 文件大小（字节）
	
	// 关联关系
	Task        *DetectTask     `gorm:"foreignKey:TaskID;references:TaskID" json:"task,omitempty"`
}

// TableName 指定表名
func (DetectResult) TableName() string {
	return "detect_results"
}

// BeforeCreate 创建前钩子
func (dr *DetectResult) BeforeCreate(tx *gorm.DB) error {
	if dr.ResultID == "" {
		dr.ResultID = uuid.New().String()
	}
	return nil
}

// QualityResult 视频质量检测结果
type QualityResult struct {
	Brightness  float64 `json:"brightness"`  // 亮度
	Contrast    float64 `json:"contrast"`    // 对比度
	Sharpness   float64 `json:"sharpness"`   // 清晰度
	Noise       float64 `json:"noise"`       // 噪声
	Blur        float64 `json:"blur"`        // 模糊度
	Saturation  float64 `json:"saturation"`  // 饱和度
	Issues      []string `json:"issues"`     // 检测到的问题
}

// IntegrityResult 录像完整性检测结果
type IntegrityResult struct {
	TotalFrames   int64    `json:"total_frames"`   // 总帧数
	LostFrames    int64    `json:"lost_frames"`    // 丢失帧数
	CorruptFrames int64    `json:"corrupt_frames"` // 损坏帧数
	Gaps          []Gap    `json:"gaps"`           // 时间间隙
	Integrity     float64  `json:"integrity"`      // 完整性百分比
}

// Gap 时间间隙
type Gap struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Duration  int       `json:"duration"` // 秒
}
```

### 6.5 创建设备模型 (`internal/model/device.go`)
```go
package model

// Device 设备模型
type Device struct {
	BaseModel
	DeviceID    string `gorm:"uniqueIndex;size:64;not null" json:"device_id" validate:"required"`
	Name        string `gorm:"size:128;not null" json:"name" validate:"required"`
	Type        string `gorm:"size:32;not null" json:"type" validate:"required"`
	Location    string `gorm:"size:256" json:"location"`
	Description string `gorm:"size:512" json:"description"`
	
	// 网络信息
	IP          string `gorm:"size:45" json:"ip" validate:"omitempty,ip"`
	Port        int    `gorm:"default:554" json:"port" validate:"min=1,max=65535"`
	Protocol    string `gorm:"size:16;default:rtsp" json:"protocol"`
	
	// 认证信息
	Username    string `gorm:"size:64" json:"username"`
	Password    string `gorm:"size:128" json:"-"` // 不返回密码
	
	// 状态信息
	Status      string `gorm:"size:32;default:active" json:"status"`
	LastSeen    *time.Time `json:"last_seen"`
	
	// 配置信息
	Config      json.RawMessage `gorm:"type:json" json:"config"`
	
	// 关联关系
	Tasks       []DetectTask `gorm:"foreignKey:DeviceID;references:DeviceID" json:"tasks,omitempty"`
}

// TableName 指定表名
func (Device) TableName() string {
	return "devices"
}

// IsOnline 检查设备是否在线
func (d *Device) IsOnline() bool {
	if d.LastSeen == nil {
		return false
	}
	return time.Since(*d.LastSeen) < 5*time.Minute
}

// GetStreamURL 获取流媒体URL
func (d *Device) GetStreamURL() string {
	if d.Username != "" && d.Password != "" {
		return fmt.Sprintf("%s://%s:%s@%s:%d/stream", 
			d.Protocol, d.Username, d.Password, d.IP, d.Port)
	}
	return fmt.Sprintf("%s://%s:%d/stream", d.Protocol, d.IP, d.Port)
}
```

### 6.6 创建工作节点模型 (`internal/model/worker.go`)
```go
package model

// Worker 工作节点模型
type Worker struct {
	BaseModel
	WorkerID    string `gorm:"uniqueIndex;size:64;not null" json:"worker_id" validate:"required"`
	Name        string `gorm:"size:128;not null" json:"name" validate:"required"`
	IP          string `gorm:"size:45;not null" json:"ip" validate:"required,ip"`
	Port        int    `gorm:"not null" json:"port" validate:"min=1,max=65535"`
	
	// 状态信息
	Status      string     `gorm:"size:32;default:idle" json:"status"`
	LastSeen    *time.Time `json:"last_seen"`
	StartedAt   time.Time  `json:"started_at"`
	
	// 性能信息
	CPUUsage    float64 `gorm:"type:decimal(5,2)" json:"cpu_usage"`
	MemoryUsage float64 `gorm:"type:decimal(5,2)" json:"memory_usage"`
	LoadAvg     float64 `gorm:"type:decimal(5,2)" json:"load_avg"`
	
	// 任务统计
	ActiveTasks    int `gorm:"default:0" json:"active_tasks"`
	CompletedTasks int `gorm:"default:0" json:"completed_tasks"`
	FailedTasks    int `gorm:"default:0" json:"failed_tasks"`
	
	// 配置信息
	MaxTasks    int             `gorm:"default:10" json:"max_tasks"`
	Capabilities []string       `gorm:"type:json" json:"capabilities"`
	Config      json.RawMessage `gorm:"type:json" json:"config"`
}

// TableName 指定表名
func (Worker) TableName() string {
	return "workers"
}

// IsOnline 检查工作节点是否在线
func (w *Worker) IsOnline() bool {
	if w.LastSeen == nil {
		return false
	}
	return time.Since(*w.LastSeen) < 30*time.Second
}

// CanAcceptTask 检查是否可以接受新任务
func (w *Worker) CanAcceptTask() bool {
	return w.IsOnline() && w.Status == "idle" && w.ActiveTasks < w.MaxTasks
}

// GetAddress 获取工作节点地址
func (w *Worker) GetAddress() string {
	return fmt.Sprintf("%s:%d", w.IP, w.Port)
}
```

### 6.7 创建数据传输对象 (`internal/model/dto.go`)
```go
package model

import (
	"time"
)

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	DeviceID   string                 `json:"device_id" validate:"required"`
	Type       DetectType             `json:"type" validate:"required,oneof=quality integrity motion object"`
	VideoURL   string                 `json:"video_url" validate:"required,url"`
	StartTime  *time.Time             `json:"start_time"`
	EndTime    *time.Time             `json:"end_time"`
	Duration   int                    `json:"duration" validate:"min=1"`
	Priority   Priority               `json:"priority" validate:"min=1,max=4"`
	Parameters map[string]interface{} `json:"parameters"`
}

// TaskResponse 任务响应
type TaskResponse struct {
	TaskID     string     `json:"task_id"`
	Status     TaskStatus `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	StartedAt  *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Progress   int        `json:"progress"` // 0-100
	Message    string     `json:"message,omitempty"`
}

// TaskListRequest 任务列表请求
type TaskListRequest struct {
	DeviceID string     `form:"device_id"`
	Type     DetectType `form:"type"`
	Status   TaskStatus `form:"status"`
	Page     int        `form:"page" validate:"min=1"`
	Size     int        `form:"size" validate:"min=1,max=100"`
	StartDate *time.Time `form:"start_date"`
	EndDate   *time.Time `form:"end_date"`
}

// TaskListResponse 任务列表响应
type TaskListResponse struct {
	Tasks      []DetectTask `json:"tasks"`
	Total      int64        `json:"total"`
	Page       int          `json:"page"`
	Size       int          `json:"size"`
	TotalPages int          `json:"total_pages"`
}

// ResultResponse 结果响应
type ResultResponse struct {
	ResultID     string          `json:"result_id"`
	TaskID       string          `json:"task_id"`
	Type         DetectType      `json:"type"`
	Score        float64         `json:"score"`
	Passed       bool            `json:"passed"`
	Confidence   float64         `json:"confidence"`
	Details      json.RawMessage `json:"details"`
	ReportURL    string          `json:"report_url,omitempty"`
	ThumbnailURL string          `json:"thumbnail_url,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
}

// HealthCheckResponse 健康检查响应
type HealthCheckResponse struct {
	Status    string                 `json:"status"`
	Timestamp int64                  `json:"timestamp"`
	Services  map[string]interface{} `json:"services"`
	Version   string                 `json:"version"`
}
```

## 验证步骤

### 6.8 创建模型测试文件 (`internal/model/model_test.go`)
```go
package model

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestDetectTaskModel(t *testing.T) {
	task := &DetectTask{
		DeviceID:  "device-001",
		Type:      DetectTypeQuality,
		VideoURL:  "rtsp://example.com/stream",
		Priority:  PriorityNormal,
		Duration:  300,
	}

	// 测试验证
	validate := validator.New()
	err := validate.Struct(task)
	assert.NoError(t, err)

	// 测试方法
	assert.False(t, task.IsCompleted())
	assert.False(t, task.CanRetry())
}

func TestDetectResultModel(t *testing.T) {
	result := &DetectResult{
		TaskID:     "task-001",
		DeviceID:   "device-001",
		Type:       DetectTypeQuality,
		Score:      0.85,
		Passed:     true,
		Confidence: 0.92,
	}

	// 测试验证
	validate := validator.New()
	err := validate.Struct(result)
	assert.NoError(t, err)
}

func TestQualityResult(t *testing.T) {
	quality := QualityResult{
		Brightness: 0.7,
		Contrast:   0.8,
		Sharpness:  0.9,
		Noise:      0.1,
		Blur:       0.05,
		Saturation: 0.75,
		Issues:     []string{"low_brightness"},
	}

	data, err := json.Marshal(quality)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "brightness")
}

func TestDeviceModel(t *testing.T) {
	device := &Device{
		DeviceID: "device-001",
		Name:     "Test Camera",
		Type:     "camera",
		IP:       "192.168.1.100",
		Port:     554,
		Username: "admin",
		Password: "password",
	}

	// 测试验证
	validate := validator.New()
	err := validate.Struct(device)
	assert.NoError(t, err)

	// 测试方法
	streamURL := device.GetStreamURL()
	assert.Contains(t, streamURL, "admin:password")
	assert.Contains(t, streamURL, "192.168.1.100:554")
}

func TestWorkerModel(t *testing.T) {
	worker := &Worker{
		WorkerID: "worker-001",
		Name:     "Test Worker",
		IP:       "192.168.1.200",
		Port:     8080,
		MaxTasks: 5,
	}

	// 测试验证
	validate := validator.New()
	err := validate.Struct(worker)
	assert.NoError(t, err)

	// 测试方法
	assert.False(t, worker.IsOnline()) // LastSeen为nil
	assert.False(t, worker.CanAcceptTask())
	
	address := worker.GetAddress()
	assert.Equal(t, "192.168.1.200:8080", address)
}
```

### 6.9 验证数据模型
```bash
# 运行模型测试
go test ./internal/model/

# 检查依赖
go mod tidy

# 验证数据结构
go run -c "
package main
import (
    \"fmt\"
    \"detect-executor-go/internal/model\"
)
func main() {
    task := &model.DetectTask{}
    fmt.Printf(\"DetectTask fields: %+v\", task)
}
"
```

## 预期结果
- 完整的数据模型定义，包含所有核心业务实体
- 数据验证规则实现，确保数据完整性
- 模型关系定义，支持关联查询
- 枚举类型定义，提高代码可读性
- 数据传输对象，规范API接口
- 完整的单元测试，验证模型功能

## 状态
- [ ] 步骤 6 完成：数据模型定义

---

## 下一步预告
步骤 7: 数据访问层实现 (`internal/repository/`)
- 任务数据访问接口
- 结果数据访问接口
- 设备数据访问接口
- 缓存集成和查询优化
