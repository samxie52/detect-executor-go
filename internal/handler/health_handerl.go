package handler

import (
	"time"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	*BaseHandler
}

func NewHealthHandler(base *BaseHandler) *HealthHandler {
	return &HealthHandler{
		BaseHandler: base,
	}
}

// Health 检查服务是否正常
// @Summary Check service health
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	h.Success(c, gin.H{"message": "Service is healthy"})
}

// Ready 检查服务是否可以接受请求
// @Summary Check service ready
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/ready [get]
func (h *HealthHandler) Ready(c *gin.Context) {
	if _, err := h.service.Task.GetTask(c.Request.Context(), "health-check"); err != nil {
		//it is normal because the task is not created
	}

	h.Success(c, gin.H{"status": "ok", "timestamp": time.Now().Unix(), "service": "detect-executor-go"})
}
