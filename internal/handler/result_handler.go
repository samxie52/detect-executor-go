package handler

import (
	"detect-executor-go/internal/model"
	"detect-executor-go/internal/service"
	"time"

	"github.com/gin-gonic/gin"
)

// ResultHandler 结果处理器
type ResultHandler struct {
	*BaseHandler
}

// NewResultHandler 创建结果处理器
func NewResultHandler(base *BaseHandler) *ResultHandler {
	return &ResultHandler{
		BaseHandler: base,
	}
}

// GetResultByTaskId 获取结果
// @Summary Get result by task id
// @Tags Result
// @Accept json
// @Produce json
// @Param taskId path string true "Task ID"
// @Success 200 {object} model.DetectResult
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/detect/result/{taskId} [get]
func (h *ResultHandler) GetResultByTaskId(c *gin.Context) {
	taskID := c.Param("taskId")

	if taskID == "" {
		h.BadRequest(c, "Invalid task ID")
		return
	}

	result, err := h.service.Result.GetResultByTaskId(c.Request.Context(), taskID)
	if err != nil {
		h.InternalError(c, "get result failed: "+err.Error())
		return
	}

	h.Success(c, result)
}

// ListResults 列出结果
// @Summary List results
// @Tags Result
// @Accept json
// @Produce json
// @Param taskId path string true "Task ID"
// @Success 200 {object} model.DetectResult
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/detect/task/{taskId}/results [get]
func (h *ResultHandler) ListResults(c *gin.Context) {
	page, size := h.GetPageParams(c)

	req := &service.ListResultRequest{
		DeviceID: c.Query("deviceId"),
		Type:     model.DetectType(c.Query("type")),
		Page:     page,
		PageSize: size,
	}

	if passedStr := c.Query("passed"); passedStr != "" {
		passed := passedStr == "true"
		req.Passed = &passed
	}

	if startTimeStr := c.Query("startTime"); startTimeStr != "" {
		startTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			h.BadRequest(c, "Invalid start time format")
			return
		}
		req.StartTime = &startTime
	}

	if endTimeStr := c.Query("endTime"); endTimeStr != "" {
		endTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			h.BadRequest(c, "Invalid end time format")
			return
		}
		req.EndTime = &endTime
	}

	results, total, err := h.service.Result.ListResults(c.Request.Context(), req)
	if err != nil {
		h.InternalError(c, "list results failed: "+err.Error())
		return
	}

	h.PageSuccess(c, results, total, page, size)
}

// DeleteResult 删除结果
// @Summary Delete result
// @Tags Result
// @Accept json
// @Produce json
// @Param resultId path string true "Result ID"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/detect/result/{resultId} [delete]
func (h *ResultHandler) DeleteResult(c *gin.Context) {
	resultID := c.Param("resultId")

	if resultID == "" {
		h.BadRequest(c, "Invalid result ID")
		return
	}

	err := h.service.Result.DeleteResult(c.Request.Context(), resultID)
	if err != nil {
		h.InternalError(c, "delete result failed: "+err.Error())
		return
	}

	h.Success(c, gin.H{"message": "Result deleted successfully"})
}
