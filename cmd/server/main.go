package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"detect-executor-go/pkg/config"
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

	//set gin mode
	if cfg.Server.Mode == "release" {
		//release mode is production mode that gin will not print debug info
		gin.SetMode(gin.ReleaseMode)
	} else {
		//debug mode is development mode that gin will print debug info
		gin.SetMode(gin.DebugMode)
	}

	//create gin router
	router := gin.New()

	//use logger middleware
	router.Use(logger.GinLogger(log))
	//use recovery middleware
	router.Use(logger.GinRecovery(log))
	//use request id middleware
	router.Use(logger.RequestID())

	// setup routes
	setupRoutes(router, log)

	// create http server
	srv := &http.Server{
		Addr:         cfg.Server.GetAddr(),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	//run http server
	go func() {
		log.WithField("addr", cfg.Server.GetAddr()).Info("Starting http server...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithField("error", err).Error("http server error")
		}
	}()

	// wait for interrupt signal and shutdown server gracefully
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.WithField("error", err).Error("http server shutdown error")
	} else {
		log.Info("Server exited")
	}

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
