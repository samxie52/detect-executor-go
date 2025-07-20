package service

import (
	"context"
	"detect-executor-go/internal/model"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	ctx     context.Context
	manager *Manager
}

func (suite *ServiceTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	//注意，这里需要mock repository和client
	// 在实际测试中应该使用测试数据库或mock对象
}

// TestTaskService 测试任务服务
func (suite *ServiceTestSuite) TestTaskService() {
	// 测试创建任务
	req := &CreateTaskRequest{
		DeviceID:    "device-001",
		Type:        model.DetectTypeQuality,
		VideoURL:    "rtsp://example.com/stream",
		Priority:    model.PriorityHigh,
		Description: "Test task",
	}

	// 注意：这里需要mock设备存在且在线
	// task, err := suite.manager.Task.CreateTask(suite.ctx, req)
	// suite.NoError(err)
	// suite.NotNil(task)
	// suite.Equal(req.DeviceID, task.DeviceID)
	// suite.Equal(req.Type, task.Type)

	suite.T().Log("Task service test placeholder")
}

// TestDetectService 测试检测服务
func (suite *ServiceTestSuite) TestDetectService() {
	// 测试开始检测
	taskID := "test-task-001"

	// 注意：这里需要mock任务存在
	// err := suite.manager.Detect.StartDetection(suite.ctx, taskID)
	// suite.NoError(err)

	// 测试获取检测状态
	// status, err := suite.manager.Detect.GetDetectionStatus(suite.ctx, taskID)
	// suite.NoError(err)
	// suite.NotNil(status)
	// suite.Equal(taskID, status.TaskID)

	suite.T().Log("Detect service test placeholder")
}

// TestResultService 测试结果服务
func (suite *ServiceTestSuite) TestResultService() {
	// 测试列表结果
	req := &ListResultRequest{
		Page:     1,
		PageSize: 10,
	}

	// 注意：这里需要mock数据
	// results, total, err := suite.manager.Result.ListResults(suite.ctx, req)
	// suite.NoError(err)
	// suite.NotNil(results)
	// suite.GreaterOrEqual(total, int64(0))

	suite.T().Log("Result service test placeholder")
}

// TestServiceTestSuite 运行测试套件
func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
