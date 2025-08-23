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

	router.GET("/value/:metricType/:metricName", func(c *gin.Context) {
		metricsHandler.GetValue(c)
	})

	router.POST("/update/:metricType/:metricName/:metricValue", func(c *gin.Context) {
		metricsHandler.PostUpdate(c)
	})

	router.POST("/update", func(c *gin.Context) {
		metricsHandler.PostUpdateJSON(c)
	})

	router.POST("/value", func(c *gin.Context) {
		metricsHandler.GetValueJSON(c)
	})
	return router
}
