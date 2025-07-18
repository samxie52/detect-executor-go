package model

import (
	"encoding/json"
	"fmt"
	"time"
)

type Worker struct {
	BaseModel
	WorkerID string `gorm:"uniqueIndex;size:64;not null" json:"worker_id" validate:"required"`
	Name     string `gorm:"size:128;not null" json:"name" validate:"required"`
	IP       string `gorm:"size:64" json:"ip" validate:"required,ip"`
	Port     int    `gorm:"not null" json:"port" validate:"min=1,max=65535"`

	Status    string     `gorm:"size:32 default:'idle'" json:"status"`
	LastSeen  *time.Time `json:"last_seen"`
	StartedAt *time.Time `json:"started_at"`

	CPUUsage    float64 `gorm:"type:decimal(5,2)" json:"cpu_usage"`
	MemoryUsage float64 `gorm:"type:decimal(5,2)" json:"memory_usage"`
	LoadAvg     float64 `gorm:"type:decimal(5,2)" json:"load_avg"`

	ActiveTasks    int `gorm:"default:0" json:"active_tasks"`
	CompletedTasks int `gorm:"default:0" json:"completed_tasks"`
	FailedTasks    int `gorm:"default:0" json:"failed_tasks"`

	MaxTasks int `gorm:"default:10" json:"max_tasks"`
	//支持的功能，例如 ["quality", "integrity", "motion", "object"]
	Capabilities []string        `gorm:"type:json" json:"capabilities"`
	Config       json.RawMessage `gorm:"type:json" json:"config"`
}

// TableName 指定表名, gorm 会自动将结构体名转换为蛇形命名, 例如 Worker -> worker
func (Worker) TableName() string {
	return "workers"
}

// IsOnline 判断工作者是否在线
func (w *Worker) IsOnline() bool {
	if w.LastSeen == nil {
		return false
	}
	return time.Since(*w.LastSeen) < 30*time.Second
}

//CanAcceptTask 判断工作者是否可以接受任务

func (w *Worker) CanAcceptTask() bool {
	//idle 空闲状态
	//w.ActiveTasks < w.MaxTasks 表示当前活动任务数小于最大任务数
	return w.IsOnline() && w.Status == "idle" && w.ActiveTasks < w.MaxTasks
}

// GetAddress 获取工作者地址
func (w *Worker) GetAddress() string {
	return fmt.Sprintf("%s:%d", w.IP, w.Port)
}
