package handler

import (
	"detect-executor-go/internal/service"
	"detect-executor-go/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Router 路由器
type Router struct {
	engine *gin.Engine
	logger *logger.Logger
}

func NewRouter(log *logger.Logger, service *service.Manager) *Router {
	base := NewBaseHandler(log, service)

	taskHandler := NewTaskHandler(base)
	detectHandler := NewDetectHandler(base)
	resultHandler := NewResultHandler(base)
	healthHandler := NewHealthHandler(base)

	engine := gin.New()

	engine.Use(gin.Recovery())
	engine.Use(CORSMiddleware())
	engine.Use(logger.GinLogger(log))
	engine.Use(logger.RequestID())

	//health check
	engine.GET("/api/v1/health", healthHandler.Health)
	engine.GET("/api/v1/ready", healthHandler.Ready)

	v1 := engine.Group("/api/v1")
	{
		task := v1.Group("/task")
		{
			task.GET("", taskHandler.ListTasks)
			task.POST("", taskHandler.CreateTask)
			task.GET("/:taskId", taskHandler.GetTask)
			task.PUT("/:taskId/cancel", taskHandler.CancelTask)
			task.PUT("/:taskId/retry", taskHandler.RetryTask)
		}

		detection := v1.Group("/detection")
		{
			detection.POST("/:taskId/start", detectHandler.StartDetection)
			detection.POST("/:taskId/stop", detectHandler.StopDetection)
			detection.GET("/:taskId/status", detectHandler.GetDetectionStatus)
		}

		result := v1.Group("/result")
		{
			result.GET("/:taskId", resultHandler.GetResultByTaskId)
			result.GET("", resultHandler.ListResults)
			result.DELETE("/:resultId", resultHandler.DeleteResult)
		}

		// Swagger 文档路由（需要时可以添加）
		// engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return &Router{
		engine: engine,
		logger: log,
	}
}

// CORS 中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

// GetEngine 获取引擎
func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}
