package handler

import "github.com/gin-gonic/gin"

// DetectHandler 检测处理器
type DetectHandler struct {
	*BaseHandler
}

// NewDetectHandler 创建检测处理器
func NewDetectHandler(base *BaseHandler) *DetectHandler {
	return &DetectHandler{
		BaseHandler: base,
	}
}

// StartDetection 开始检测
// @Summary Start detection
// @Tags Detection
// @Accept json
// @Produce json
// @Param taskId path string true "Task ID"
// @Success 200 {object} model.DetectTask
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/detect/task/{taskId}/start [post]
func (h *DetectHandler) StartDetection(c *gin.Context) {
	taskID := c.Param("taskId")

	if taskID == "" {
		h.BadRequest(c, "Invalid task ID")
		return
	}

	err := h.service.Detect.StartDetection(c.Request.Context(), taskID)
	if err != nil {
		h.InternalError(c, "start detection failed: "+err.Error())
		return
	}

	h.Success(c, gin.H{"message": "Task started successfully"})
}

// StopDetection 停止检测
// @Summary Stop detection
// @Tags Detection
// @Accept json
// @Produce json
// @Param taskId path string true "Task ID"
// @Success 200 {object} model.DetectTask
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/detect/task/{taskId}/stop [post]
func (h *DetectHandler) StopDetection(c *gin.Context) {
	taskID := c.Param("taskId")

	if taskID == "" {
		h.BadRequest(c, "Invalid task ID")
		return
	}

	err := h.service.Detect.StopDetection(c.Request.Context(), taskID)
	if err != nil {
		h.InternalError(c, "stop detection failed: "+err.Error())
		return
	}

	h.Success(c, gin.H{"message": "Task stopped successfully"})
}

// GetDetectionStatus 获取检测状态
// @Summary Get detection status
// @Tags Detection
// @Accept json
// @Produce json
// @Param taskId path string true "Task ID"
// @Success 200 {object} model.DetectTask
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/detect/task/{taskId}/status [get]
func (h *DetectHandler) GetDetectionStatus(c *gin.Context) {
	taskID := c.Param("taskId")

	if taskID == "" {
		h.BadRequest(c, "Invalid task ID")
		return
	}

	status, err := h.service.Detect.GetDetectionStatus(c.Request.Context(), taskID)
	if err != nil {
		h.InternalError(c, "get detection status failed: "+err.Error())
		return
	}

	h.Success(c, status)
}
