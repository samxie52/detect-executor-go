package model

import (
	"fmt"
	"time"

	"encoding/json"
)

type Device struct {
	BaseModel

	DeviceID    string `gorm:"uniqueIndex;size:64;not null" json:"device_id" validate:"required"`
	Name        string `gorm:"size:64;not null" json:"name" validate:"required"`
	Type        string `gorm:"size:64;not null" json:"type" validate:"required"`
	Location    string `gorm:"size:64" json:"location"`
	Description string `gorm:"size:256" json:"description"`

	//IP 地址 validate:"omitempty,ip" 表示 IP 地址是可选的，且必须是 IP 地址
	IP       string `gorm:"size:64" json:"ip" validate:"omitempty,ip"`
	Port     int    `gorm:"default:554" json:"port" validate:"min=1,max=65535"`
	Protocol string `gorm:"size:16 default:'rtsp'" json:"protocol"`

	Username string `gorm:"size:64" json:"username"`
	Password string `gorm:"size:128" json:"-"`

	Status string `gorm:"size:32 default:'active'" json:"status"`
	//最后在线时间
	LastSeen *time.Time `json:"last_seen"`

	Config json.RawMessage `gorm:"type:json" json:"config"`

	//外键关联 foreignKey:DeviceID 表示外键是 DeviceID references:DeviceID 表示引用的表是 DeviceID
	// 临时移除外键约束以避免表创建顺序问题
	// Tasks []DetectTask `gorm:"foreignKey:DeviceID;references:DeviceID" json:"tasks,omitempty"`
}

// TableName 指定表名, gorm 会自动将结构体名转换为蛇形命名, 例如 Device -> device
func (Device) TableName() string {
	return "devices"
}

// IsOnline 判断设备是否在线
func (d *Device) IsOnline() bool {
	if d.LastSeen == nil {
		return false
	}
	// time.Since() 计算当前时间与 LastSeen 的时间差
	return time.Since(*d.LastSeen) < 5*time.Minute
}

// GetStreamURL 获取设备流地址
func (d *Device) GetStreamURL() string {
	if d.Username != "" && d.Password != "" {
		// example: rtsp://username:password@ip:port/stream
		return fmt.Sprintf("%s://%s:%s@%s:%d/stream", d.Protocol, d.Username, d.Password, d.IP, d.Port)
	}
	// example: rtsp://ip:port/stream
	return fmt.Sprintf("%s://%s:%d/stream", d.Protocol, d.IP, d.Port)
}
