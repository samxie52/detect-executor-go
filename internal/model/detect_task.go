package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DetectTask struct {
	BaseModel
	//uniqueIndex 表示该字段是唯一索引，size:36 表示该字段的长度为36，not null 表示该字段不能为空
	// validate:required 表示该字段不能为空
	TaskID     string     `gorm:"uniqueIndex;size:36;not null" json:"task_id" validate:"required"`
	DeviceID   string     `gorm:"index;size:64;not null" json:"device_id" validate:"required"`
	DeviceName string     `gorm:"size:168" json:"device_name"`
	Type       DetectType `gorm:"size:32;not null" json:"type" validate:"required,oneof=quality integrity motion object"`
	Status     TaskStatus `gorm:"size:32;not null default=pending" json:"status" validate:"required,oneof=pending running success failed cancel timeout"`
	Priority   Priority   `gorm:"not null default=2" json:"priority" validate:"required,oneof=1 2 3 4"`

	//validate:required,url 表示该字段不能为空，且必须是url
	VideoURL  string     `gorm:"size:256;not null" json:"video_url" validate:"required,url"`
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
	Duration  int        `json:"duration"`
	//json.RawMessage 表示该字段是json格式的字符串
	Parameters map[string]interface{} `json:"parameters"`

	Description string `json:"description"`

	WorkerID string     `gorm:"size:64" json:"worker_id"`
	StartAt  *time.Time `json:"start_at"`
	//完成时间
	CompletedAt *time.Time `json:"completed_at"`
	RetryCount  int        `gorm:"default=0" json:"retry_count"`
	MaxRetries  int        `gorm:"default=3" json:"max_retries"`

	ResultID string `gorm:"size:36" json:"result_id"`
	ErrorMsg string `gorm:"size:1024" json:"error_msg"`

	// 临时移除外键约束以避免表创建顺序问题
	// Result []DetectResult `gorm:"foreignKey:TaskID;references:TaskID" json:"result,omitempty"`
}

// TableName 指定表名, gorm 会自动将结构体名转换为蛇形命名, 例如 DetectTask -> detect_task
// func(DetectTask) 表示该方法所属的结构体是 DetectTask
func (DetectTask) TableName() string {
	return "detect_task"
}

// BeforeCreate 在创建记录之前调用,会在 gorm 自动执行
func (dt *DetectTask) BeforeCreate(tx *gorm.DB) error {
	if dt.TaskID == "" {
		//uuid.New() 生成一个 uuid
		dt.TaskID = uuid.New().String()
	}
	return nil
}

// isCompleted 判断任务是否完成
func (dt *DetectTask) IsCompleted() bool {
	return dt.Status == TaskStatusCompleted ||
		dt.Status == TaskStatusFailed ||
		dt.Status == TaskStatusCancel ||
		dt.Status == TaskStatusTimeout
}

// CanRetry 判断任务是否可以重试
func (dt *DetectTask) CanRetry() bool {
	return dt.RetryCount < dt.MaxRetries && dt.Status == TaskStatusFailed
}

// GetDuration 获取任务持续时间
func (dt *DetectTask) GetDuration() time.Duration {
	if dt.StartAt == nil || dt.CompletedAt == nil {
		return 0
	}
	//Sub 计算两个时间的差值
	return dt.CompletedAt.Sub(*dt.StartAt)
}
