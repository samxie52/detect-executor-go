package service

import (
	"context"
	"detect-executor-go/internal/model"
	"fmt"
	"time"
)

type ResultService interface {
	BaseService
	GetResult(ctx context.Context, resultID string) (*model.DetectResult, error)
	GetResultByTaskId(ctx context.Context, taskID string) (*model.DetectResult, error)
	ListResults(ctx context.Context, req *ListResultRequest) ([]model.DetectResult, int64, error)
	DeleteResult(ctx context.Context, resultID string) error
}

// ListResultRequest 列出检测结果的请求
type ListResultRequest struct {
	DeviceID  string           `json:"device_id,omitempty"`
	Type      model.DetectType `json:"type,omitempty"`
	Passed    *bool            `json:"passed,omitempty"`
	StartTime *time.Time       `json:"start_time,omitempty"`
	EndTime   *time.Time       `json:"end_time,omitempty"`
	Page      int              `json:"page,omitempty"`
	PageSize  int              `json:"page_size,omitempty"`
}

// // StatisticsRequest 统计检测结果的请求
// type StatisticsRequest struct {
// 	DeviceID  string           `json:"device_id,omitempty"`
// 	Type      model.DetectType `json:"type,omitempty"`
// 	StartTime *time.Time       `json:"start_time,omitempty"`
// 	EndTime   *time.Time       `json:"end_time,omitempty"`
// 	GroupBy   string           `json:"group_by"` // day, week, month
// }

// // ResultStatistics 统计检测结果
// type ResultStatistics struct {
// 	TotalCount    int64                 `json:"total_count"`
// 	PassedCount   int64                 `json:"passed_count"`
// 	FailedCount   int64                 `json:"failed_count"`
// 	PassRate      float64               `json:"pass_rate"`
// 	AvgScore      float64               `json:"avg_score"`
// 	AvgConfidence float64               `json:"avg_confidence"`
// 	GroupData     []StatisticsGroupData `json:"group_data"`
// }

// // StatisticsGroupData 统计检测结果的分组数据
// type StatisticsGroupData struct {
// 	Date        string  `json:"date"`
// 	Count       int64   `json:"count"`
// 	PassedCount int64   `json:"passed_count"`
// 	PassRate    float64 `json:"pass_rate"`
// 	AvgScore    float64 `json:"avg_score"`
// }

// resultService 结果服务
type resultService struct {
	ctx *ServiceContext
}

func NewResultService(ctx *ServiceContext) ResultService {
	return &resultService{
		ctx: ctx,
	}
}

// GetResult
func (s *resultService) GetResult(ctx context.Context, resultID string) (*model.DetectResult, error) {
	result, err := s.ctx.Repository.Result.GetByResultID(ctx, resultID)
	if err != nil {
		return nil, fmt.Errorf("get result failed: %w", err)
	}
	return result, nil
}

// GetResultByTaskId
func (s *resultService) GetResultByTaskId(ctx context.Context, taskID string) (*model.DetectResult, error) {
	result, err := s.ctx.Repository.Result.GetByTaskID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("get result failed: %w", err)
	}
	return result, nil
}

// ListResults
func (s *resultService) ListResults(ctx context.Context, req *ListResultRequest) ([]model.DetectResult, int64, error) {

	results, err := s.ctx.Repository.Result.List(ctx, req.DeviceID, req.Type, req.Page, req.PageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list results failed: %w", err)
	}
	return results, 0, nil
}

// DeleteResult
func (s *resultService) DeleteResult(ctx context.Context, resultID string) error {
	if err := s.ctx.Repository.Result.DeleteByResultID(ctx, resultID); err != nil {
		return fmt.Errorf("delete result failed: %w", err)
	}
	return nil
}

// Close
func (s *resultService) Close() error {
	return nil
}
