package handler

import (
	"detect-executor-go/internal/service"
	"detect-executor-go/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// BaseHandler 基础处理器
type BaseHandler struct {
	logger    *logger.Logger
	validator *validator.Validate
	service   *service.Manager
}

// NewBaseHandler 创建基础处理器
func NewBaseHandler(logger *logger.Logger, service *service.Manager) *BaseHandler {
	return &BaseHandler{
		logger:    logger,
		validator: validator.New(),
		service:   service,
	}
}

// Response 响应
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PageResponse 分页响应
type PageResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Total   int64       `json:"total"`
	Page    int         `json:"page"`
	Size    int         `json:"size"`
}

// Success 响应成功
// gin.Context: gin 框架的上下文
func (h *BaseHandler) Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// PageSuccess 分页响应成功
func (h *BaseHandler) PageSuccess(c *gin.Context, data interface{}, total int64, page int, size int) {
	c.JSON(http.StatusOK, PageResponse{
		Code:    0,
		Message: "success",
		Data:    data,
		Total:   total,
		Page:    page,
		Size:    size,
	})
}

// Error 响应错误
func (h *BaseHandler) Error(c *gin.Context, code int, message string) {
	h.logger.WithFields(map[string]interface{}{
		"path":   c.Request.URL.Path,
		"method": c.Request.Method,
		"error":  message,
	}).Error("API request error")

	c.JSON(code, Response{
		Code:    code,
		Message: message,
	})
}

// 400错误
func (h *BaseHandler) BadRequest(c *gin.Context, message string) {
	h.Error(c, http.StatusBadRequest, message)
}

// InternalError 内部错误
func (h *BaseHandler) InternalError(c *gin.Context, message string) {
	h.Error(c, http.StatusInternalServerError, message)
}

// 404错误
func (h *BaseHandler) NotFound(c *gin.Context, message string) {
	h.Error(c, http.StatusNotFound, message)
}

// ValidateRequest 验证请求
func (h *BaseHandler) ValidateRequest(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		h.BadRequest(c, "Invalid request parameters"+err.Error())
		return false
	}
	return true
}

// CheckServiceAvailable 检查服务是否可用
func (h *BaseHandler) CheckServiceAvailable(c *gin.Context) bool {
	if h.service == nil {
		h.InternalError(c, "Service unavailable - database connection required")
		return false
	}
	return true
}

// GetPageParams 获取分页参数
func (h *BaseHandler) GetPageParams(c *gin.Context) (int, int) {
	//strconv 	是Go语言中的一个包，用于字符串转换
	//strconv.Atoi()将字符串转换为整数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	return page, size
}

// GetRequestID 获取请求ID
func (h *BaseHandler) GetRequestID(c *gin.Context) string {
	return c.GetHeader("request_id")
}
