package handler

import (
	"detect-executor-go/internal/model"
	"detect-executor-go/internal/service"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	*BaseHandler
}

// NewTaskHandler 创建任务处理器
func NewTaskHandler(base *BaseHandler) *TaskHandler {
	return &TaskHandler{
		BaseHandler: base,
	}
}

// CreateTask 创建任务
// @Summary Create a new task
// @Tags Task
// @Accept json
// @Produce json
// @Param task body service.CreateTaskRequest true "Create task"
// @Success 200 {object} model.DetectTask
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/detect/task [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	if !h.CheckServiceAvailable(c) {
		return
	}

	var req service.CreateTaskRequest
	if !h.ValidateRequest(c, &req) {
		return
	}

	task, err := h.service.Task.CreateTask(c.Request.Context(), &req)
	if err != nil {
		h.InternalError(c, "create task failed: "+err.Error())
		return
	}

	h.Success(c, task)
}

// GetTask 获取任务
// @Summary Get a task by taskId
// @Tags Task
// @Accept json
// @Produce json
// @Param taskId path string true "Task ID"
// @Success 200 {object} model.DetectTask
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/detect/task/{taskId} [get]
func (h *TaskHandler) GetTask(c *gin.Context) {
	if !h.CheckServiceAvailable(c) {
		return
	}

	taskID := c.Param("taskId")
	task, err := h.service.Task.GetTask(c.Request.Context(), taskID)
	if err != nil {
		h.InternalError(c, "get task failed: "+err.Error())
		return
	}

	h.Success(c, task)
}

// ListTasks 列出任务
// @Summary List tasks
// @Tags Task
// @Accept json
// @Produce json
// @Success 200 {object} PageResponse
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/detect/task [get]
func (h *TaskHandler) ListTasks(c *gin.Context) {
	if !h.CheckServiceAvailable(c) {
		return
	}

	page, size := h.GetPageParams(c)

	req := &model.TaskListRequest{
		DeviceID: c.Query("deviceId"),
		Status:   model.TaskStatus(c.Query("status")),
		Type:     model.DetectType(c.Query("type")),
		Page:     page,
		Size:     size,
	}

	if err := h.validator.Struct(req); err != nil {
		h.BadRequest(c, "Invalid request parameters"+err.Error())
		return
	}

	tasks, total, err := h.service.Task.ListTasks(c.Request.Context(), req)
	if err != nil {
		h.InternalError(c, "list tasks failed: "+err.Error())
		return
	}

	h.PageSuccess(c, tasks, total, page, size)

}

// CancelTask 取消任务
// @Summary Cancel a task by taskId
// @Tags Task
// @Accept json
// @Produce json
// @Param taskId path string true "Task ID"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/detect/task/{taskId}/cancel [put]
func (h *TaskHandler) CancelTask(c *gin.Context) {
	if !h.CheckServiceAvailable(c) {
		return
	}

	taskID := c.Param("taskId")

	if taskID == "" {
		h.BadRequest(c, "Invalid task ID")
		return
	}

	err := h.service.Task.CancelTask(c.Request.Context(), taskID)
	if err != nil {
		h.InternalError(c, "cancel task failed: "+err.Error())
		return
	}

	// gin.H{} 是 map[string]interface{} 的简写
	h.Success(c, gin.H{"message": "Task cancelled successfully"})
}

// RetryTask 重试任务
// @Summary Retry a task by taskId
// @Tags Task
// @Accept json
// @Produce json
// @Param taskId path string true "Task ID"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/detect/task/{taskId}/retry [put]
func (h *TaskHandler) RetryTask(c *gin.Context) {
	if !h.CheckServiceAvailable(c) {
		return
	}

	taskID := c.Param("taskId")

	if taskID == "" {
		h.BadRequest(c, "Invalid task ID")
		return
	}

	err := h.service.Task.RetryTask(c.Request.Context(), taskID)
	if err != nil {
		h.InternalError(c, "retry task failed: "+err.Error())
		return
	}

	h.Success(c, gin.H{"message": "Task retried successfully"})
}
