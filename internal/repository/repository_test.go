package repository

import (
	"context"
	"detect-executor-go/internal/model"
	"detect-executor-go/pkg/database"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// RepositoryTestSuite is a test suite for repository
// suite 是 testify 的测试套件，用于组织和运行测试
type RepositoryTestSuite struct {
	suite.Suite
	manager *Manager
	ctx     context.Context
}

func (suite *RepositoryTestSuite) SetupSuite() {
	dbManager, err := database.NewManager(&database.Config{
		MySQL: &database.MySQLConfig{
			Host:     "192.168.5.16",
			Port:     3306,
			Username: "root",
			Password: "root123",
			Database: "test",
			Charset:  "utf8mb4",
		},
		Redis: &database.RedisConfig{
			Host:     "192.168.5.16",
			Port:     6379,
			Password: "",
			Database: 0,
		},
	}, nil)

	suite.Require().NoError(err)

	// auto migrate
	err = dbManager.GetMySQL().DB.AutoMigrate(
		&model.DetectTask{},
		&model.DetectResult{},
		&model.Device{},
		&model.Worker{},
	)

	suite.Require().NoError(err)

	suite.manager = NewManager(dbManager.GetMySQL().DB, dbManager.GetRedis().Client)

	suite.ctx = context.Background()
}

func (suite *RepositoryTestSuite) TeardownSuite() {
	// clear test data

	db := suite.manager.GetDb()
	db.Exec("DELETE FROM detect_task")
	db.Exec("DELETE FROM detect_result")
	db.Exec("DELETE FROM devices")
	db.Exec("DELETE FROM workers")
}

func (suite *RepositoryTestSuite) TestTaskRepository() {
	// 清理可能存在的测试数据
	suite.manager.GetDb().Exec("DELETE FROM detect_task WHERE task_id = ?", "test_task_id")

	task := &model.DetectTask{
		TaskID:   "test_task_id",
		Type:     model.DetectTypeQuality,
		VideoURL: "rtsp://test_device_id/stream",
		Priority: model.PriorityMedium,
		Duration: 100,
		Status:   model.TaskStatusPending,
	}

	err := suite.manager.Task.Create(suite.ctx, task)
	suite.NoError(err)
	suite.NotEmpty(task.TaskID)

	found, err := suite.manager.Task.GetByTaskID(suite.ctx, task.TaskID)
	suite.NoError(err)
	suite.NotNil(found)
	suite.Equal(task.DeviceID, found.DeviceID)

	err = suite.manager.Task.UpdateStatus(suite.ctx, task.TaskID, model.TaskStatusCompleted)
	suite.NoError(err)

	updated, err := suite.manager.Task.GetByTaskID(suite.ctx, task.TaskID)
	suite.NoError(err)
	suite.NotNil(updated)
	suite.Equal(model.TaskStatusCompleted, updated.Status)

	// 使用重试机制处理可能的数据库连接问题
	var count int64
	var countErr error
	for i := 0; i < 3; i++ {
		count, countErr = suite.manager.Task.CountByStatus(suite.ctx, model.TaskStatusCompleted)
		if countErr == nil {
			break
		}
		// 如果是连接错误，等待一下再重试
		if strings.Contains(countErr.Error(), "connection") {
			time.Sleep(time.Second)
			continue
		}
		break
	}
	suite.NoError(countErr)
	suite.Equal(int64(1), count)
}

func (suite *RepositoryTestSuite) TestResultRepositoryList() {
	// 清理可能存在的测试数据
	suite.manager.GetDb().Exec("DELETE FROM detect_result WHERE result_id = ?", "test_result_id")

	result := &model.DetectResult{
		ResultID: "test_result_id",
		TaskID:   "test_task_id",
		DeviceID: "test_device_id",
		Type:     model.DetectTypeQuality,
		Score:    0.9, // 修复：使用0.9而不是90，符合decimal(5,4)范围
		Passed:   true,
	}

	err := suite.manager.Result.Create(suite.ctx, result)
	suite.NoError(err)
	suite.NotEmpty(result.ResultID)

	found, err := suite.manager.Result.GetByResultID(suite.ctx, result.ResultID)
	suite.NoError(err)
	suite.NotNil(found)
	suite.Equal(result.DeviceID, found.DeviceID)

	result.Passed = false
	err = suite.manager.Result.Update(suite.ctx, result)
	suite.NoError(err)

	updated, err := suite.manager.Result.GetByResultID(suite.ctx, result.ResultID)
	suite.NoError(err)
	suite.NotNil(updated)
	suite.Equal(result.Passed, updated.Passed)

	count, err := suite.manager.Result.CountByDeviceID(suite.ctx, result.DeviceID)
	suite.NoError(err)
	suite.Equal(int64(1), count)
}

func (suite *RepositoryTestSuite) TestDeviceRepositoryList() {
	// 清理可能存在的测试数据
	suite.manager.GetDb().Exec("DELETE FROM devices WHERE device_id = ?", "test_device_id")

	// 使用当前时间而不是零值时间
	now := time.Now()
	device := &model.Device{
		DeviceID: "test_device_id",
		Name:     "test_device_name",
		Type:     "camera",
		Status:   "online",
		LastSeen: &now,
	}

	err := suite.manager.Device.Create(suite.ctx, device)
	suite.NoError(err)
	suite.NotEmpty(device.DeviceID)

	found, err := suite.manager.Device.GetByDeviceID(suite.ctx, device.DeviceID)
	suite.NoError(err)
	suite.NotNil(found)
	suite.Equal(device.DeviceID, found.DeviceID)

	device.Status = "offline"
	err = suite.manager.Device.Update(suite.ctx, device)
	suite.NoError(err)

	updated, err := suite.manager.Device.GetByDeviceID(suite.ctx, device.DeviceID)
	suite.NoError(err)
	suite.NotNil(updated)
	suite.Equal(device.Status, updated.Status)

}

func (suite *RepositoryTestSuite) TestWorkerRepositoryList() {
	// 清理可能存在的测试数据
	suite.manager.GetDb().Exec("DELETE FROM workers WHERE worker_id = ?", "test_worker_id")

	// 使用当前时间而不是零值时间
	now := time.Now()
	worker := &model.Worker{
		WorkerID: "test_worker_id",
		Name:     "test_worker_name",
		IP:       "127.0.0.1",
		Port:     8080,
		Status:   "online",
		LastSeen: &now,
	}

	err := suite.manager.Worker.Create(suite.ctx, worker)
	suite.NoError(err)
	suite.NotEmpty(worker.WorkerID)

	found, err := suite.manager.Worker.GetByWorkerID(suite.ctx, worker.WorkerID)
	suite.NoError(err)
	suite.NotNil(found)
	suite.Equal(worker.WorkerID, found.WorkerID)

	worker.Status = "offline"
	err = suite.manager.Worker.Update(suite.ctx, worker)
	suite.NoError(err)

	updated, err := suite.manager.Worker.GetByWorkerID(suite.ctx, worker.WorkerID)
	suite.NoError(err)
	suite.NotNil(updated)
	suite.Equal(worker.Status, updated.Status)
}

func (suite *RepositoryTestSuite) TestTransaction() {
	// 清理可能存在的测试数据
	suite.manager.GetDb().Exec("DELETE FROM detect_task WHERE task_id = ?", "test_task_id_transaction")
	suite.manager.GetDb().Exec("DELETE FROM detect_result WHERE result_id = ?", "test_result_id_transaction")

	task := &model.DetectTask{
		TaskID:   "test_task_id_transaction",
		Type:     model.DetectTypeQuality,
		VideoURL: "rtsp://test_device_id/stream",
		Priority: model.PriorityMedium,
		Duration: 100,
		Status:   model.TaskStatusPending,
	}

	result := &model.DetectResult{
		ResultID: "test_result_id_transaction",
		DeviceID: "test_device_id_transaction",
		Type:     model.DetectTypeQuality,
		Score:    0.9, // 修复：使用0.9而不是90，符合decimal(5,4)范围
		Passed:   true,
	}

	err := suite.manager.Transaction(func(tx *Manager) error {
		if err := tx.Task.Create(suite.ctx, task); err != nil {
			return err
		}
		result.TaskID = task.TaskID
		if err := tx.Result.Create(suite.ctx, result); err != nil {
			return err
		}
		return nil
	})
	suite.NoError(err)
	suite.NotEmpty(task.TaskID)
	suite.NotEmpty(result.ResultID)
	suite.Equal(task.TaskID, result.TaskID)
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
