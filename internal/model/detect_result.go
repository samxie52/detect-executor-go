package model

import (
	"encoding/json"
	"time"
)

type DetectResult struct {
	BaseModel

	ResultID string     `gorm:"uniqueIndex;size:36;not null" json:"result_id" validate:"required"`
	TaskID   string     `gorm:"index;size:36;not null" json:"task_id" validate:"required"`
	DeviceID string     `gorm:"index;size:64;not null" json:"device_id" validate:"required"`
	Type     DetectType `gorm:"size:32;not null" json:"type"`

	//分数 5位小数 type:decimal(5,4) 中 5 表示总位数，4 表示小数位数
	Score      float64 `gorm:"type:decimal(5,4)" json:"score"`
	Passed     bool    `json:"passed"`
	Confidence float64 `gorm:"type:decimal(5,4)" json:"confidence"`

	Details  json.RawMessage `gorm:"type:json" json:"details"`
	Metadata json.RawMessage `gorm:"type:json" json:"metadata"`

	ReportURL    string `gorm:"size:256" json:"report_url"`
	ThumbnailURL string `gorm:"size:256" json:"thumbnail_url"`

	ProcessTime int   `json:"process_time"` //处理时间(毫秒)
	FileSize    int64 `json:"file_size"`    //文件大小(字节)

	//外键关联 foreignKey:TaskID 表示外键是 TaskID references:TaskID 表示引用的表是 TaskID
	//json:"task,omitempty" 表示在序列化时，如果 Task 为空，不序列化
	// 临时移除外键约束以避免表创建顺序问题
	// Task *DetectTask `gorm:"foreignKey:TaskID;references:TaskID" json:"task,omitempty"`
}

// TableName 指定表名, gorm 会自动将结构体名转换为蛇形命名, 例如 DetectResult -> detect_result
func (DetectResult) TableName() string {
	return "detect_result"
}

// QualityResult 质量检测结果
type QualityResult struct {
	Brightness float64  `json:"brightness"` //亮度
	Contrast   float64  `json:"contrast"`   //对比度
	Sharpness  float64  `json:"sharpness"`  //清晰度
	Noise      float64  `json:"noise"`      //噪生
	Blur       float64  `json:"blur"`       //模糊
	Saturation float64  `json:"saturation"` //饱和度
	Issues     []string `json:"issues"`     //问题

}

// IntegrityResult 完整性检测结果
type IntegrityResult struct {
	TotalFrames     int64   `json:"total_frames"`     //总帧数
	LostFrames      int64   `json:"lost_frames"`      //丢帧数
	CorruptedFrames int64   `json:"corrupted_frames"` //破损帧数
	Gaps            []Gaps  `json:"gaps"`             //时间间隙
	Integrity       float64 `json:"integrity"`        //完整性百分比
}

type Gaps struct {
	StartTime time.Time `json:"start_time"` //开始时间
	EndTime   time.Time `json:"end_time"`   //结束时间
	Duration  int       `json:"duration"`   //持续时间(秒)
}
