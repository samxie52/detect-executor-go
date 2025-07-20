package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"detect-executor-go/internal/handler"
	"detect-executor-go/internal/repository"
	"detect-executor-go/internal/service"
	"detect-executor-go/pkg/client"
	"detect-executor-go/pkg/config"
	"detect-executor-go/pkg/database"
	"detect-executor-go/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	//load config
	configLoader := config.NewLoader()
	cfg, err := configLoader.Load("")
	if err != nil {
		fmt.Println("load config failed: ", err)
		os.Exit(1)
	}

	//init logger
	loggerConfig := &logger.Config{
		Level:      cfg.Log.Level,
		Format:     cfg.Log.Format,
		Output:     cfg.Log.Output,
		Filename:   cfg.Log.Filename,
		MaxSize:    cfg.Log.MaxSize,
		MaxAge:     cfg.Log.MaxAge,
		MaxBackups: cfg.Log.MaxBackups,
		Compress:   cfg.Log.Compress,
	}

	log, err := logger.NewLogger(loggerConfig)
	if err != nil {
		fmt.Println("failed to init logger: ", err)
		os.Exit(1)
	}
	defer log.Close()

	log.Info("Starting detect executor server...")

	//将cfg转换为json格式
	jsonCfg, err := json.Marshal(cfg)
	if err != nil {
		log.Error("failed to marshal config: ", err)
		os.Exit(1)
	}
	log.Info("config>>>>>>>>>>>> ", string(jsonCfg))

	dbManager, err := database.NewManager(&cfg.Database, &cfg.Redis, log)
	if err != nil {
		log.Error("failed to init database: ", err)
		log.Warn("continuing without database connection...")
	}
	if dbManager != nil {
		defer dbManager.Close()
	}

	var repoManage *repository.Manager
	if dbManager != nil {
		repoManage = repository.NewManager(dbManager.MySQL.DB, dbManager.Redis.Client)
		log.Info("repoManage init success")
	} else {
		log.Warn("running without database - some features may not work")
	}

	clientManager, err := client.NewManager(&client.Config{
		HTTP:         cfg.Client.HTTP,
		DetectEngine: cfg.Client.DetectEngine,
		Storage:      cfg.Client.Storage,
		MessageQueue: cfg.Client.MessageQueue,
	})
	if err != nil {
		log.Error("failed to init client manager: ", err)
	}

	log.Info("clientManager init success")
	defer clientManager.Close()

	var serviceManager *service.Manager
	if repoManage != nil {
		serviceManager, err = service.NewManager(log, repoManage, clientManager)
		if err != nil {
			log.Error("failed to init service manager: ", err)
			os.Exit(1)
		}
		log.Info("serviceManager init success")
		defer serviceManager.Close()
	} else {
		log.Warn("running without service manager - API endpoints will return errors")
	}

	router := handler.NewRouter(log, serviceManager)

	server := &http.Server{
		Addr:    cfg.Server.GetAddr(),
		Handler: router.GetEngine(),
	}

	go func() {

		log.Info("Starting http server...", map[string]interface{}{
			"addr": cfg.Server.GetAddr(),
		})
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("http server error: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("http server shutdown error: ", err)
	}
	log.Info("Server exited")

}

func setupRoutes(router *gin.Engine, log *logger.Logger) {
	//setup health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
			"service":   "detect-executor-go",
		})
	})

	//setup api v1
	v1 := router.Group("/api/v1")
	{
		// base msg
		v1.GET("/info", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":      "ok",
				"version":     "1.0.0",
				"timestamp":   time.Now().Unix(),
				"description": "Go-based detect executor",
			})
		})

		// detect
		detect := v1.Group("/detect")
		{
			// common detect task
			detect.POST("/task", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"description": "Detect task",
				})
			})

			// query the task status
			detect.GET("/task/:id", func(c *gin.Context) {
				taskID := c.Param("id")
				c.JSON(http.StatusOK, gin.H{
					"description": "Query task status",
					"task_id":     taskID,
				})
			})

			// get detect result
			detect.GET("/result/:id", func(c *gin.Context) {
				taskID := c.Param("id")
				c.JSON(http.StatusOK, gin.H{
					"description": "Get detect result",
					"task_id":     taskID,
				})
			})
		}
	}

	// handle 404
	router.NoRoute(func(c *gin.Context) {
		log.WithField("path", c.Request.URL.Path).Error("Route Not Found")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Route Not Found",
			"message": "The requested route does not exist",
			"path":    c.Request.URL.Path,
		})
	})

}
