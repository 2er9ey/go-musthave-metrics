package main

import (
	"io"

	"github.com/2er9ey/go-musthave-metrics/internal/handler"
	"github.com/2er9ey/go-musthave-metrics/internal/logger"
	"github.com/gin-gonic/gin"
)

func SetupRouter(metricsHandler handler.MetricHandler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	router := gin.New()
	//router := gin.Default()

	router.Use(gin.Recovery())
	router.Use(logger.LoggerMiddleware())
	router.Use(handler.GzipMiddleware())

	router.GET("/", func(c *gin.Context) {
		metricsHandler.GetAll(c)
	})

	router.GET("/ping", func(c *gin.Context) {
		metricsHandler.Ping(c)
	})

	valueGroup := router.Group("/value")
	valueGroup.GET("/:metricType/:metricName", func(c *gin.Context) {
		metricsHandler.GetValue(c)
	})
	valueGroup.POST("/", func(c *gin.Context) {
		metricsHandler.GetValueJSON(c)
	})

	router.POST("/updates", func(c *gin.Context) {
		logger.Log.Debug("POST request for /updates")
		metricsHandler.PostUpdateJSON(c)
	})

	updateGroup := router.Group("/update")

	updateGroup.POST("/:metricType/:metricName/:metricValue", func(c *gin.Context) {
		metricsHandler.PostUpdate(c)
	})

	updateGroup.POST("/", func(c *gin.Context) {
		logger.Log.Debug("POST request for /update")
		metricsHandler.PostUpdateJSON(c)
	})

	return router
}
