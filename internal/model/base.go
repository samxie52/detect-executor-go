package model

import (
	"time"

	"gorm.io/gorm"
)

//grom 是一个ORM框架，用来操作数据库
// primarykey 是主键
//json 用来序列化和反序列化

// BaseModel 基础模型
type BaseModel struct {
	//uint 无符号整数
	ID uint `gorm:"primarykey" json:"id"`
	//time.Time 表示时间类型，自动填充 例如 2025-07-16 15:09:01
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	//gorm:"index" 表示该字段是索引
	//json:"-" 表示该字段不参与json序列化
	//gorm.DeletedAt 表示软删除
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TaskStatus 任务状态枚举
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"   //待处理
	TaskStatusRunning   TaskStatus = "running"   //处理中
	TaskStatusCompleted TaskStatus = "completed" //完成
	TaskStatusFailed    TaskStatus = "failed"    //失败
	TaskStatusCancel    TaskStatus = "cancel"    //取消
	TaskStatusTimeout   TaskStatus = "timeout"   //超时
)

// DetectType 检测类型枚举
type DetectType string

const (
	DetectTypeQuality   DetectType = "quality"   //质量检测
	DetectTypeIntegrity DetectType = "integrity" //完整性检测
	DetectTypeMotion    DetectType = "motion"    //运动检测
	DetectTypeObject    DetectType = "object"    //目标检测
)

type Priority int

const (
	PriorityLow    Priority = 1 //低优先级
	PriorityMedium Priority = 2 //中优先级
	PriorityHigh   Priority = 3 //高优先级
	PriorityUrgent Priority = 4 //紧急优先级
)
