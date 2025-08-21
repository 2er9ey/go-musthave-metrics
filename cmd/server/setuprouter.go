package main

import (
	"io"

	"github.com/2er9ey/go-musthave-metrics/internal/handler"
	"github.com/gin-gonic/gin"
)

func SetupRouter(metricsHandler handler.MetricHandler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/", func(c *gin.Context) {
		metricsHandler.GetAll(c)
	})

	router.GET("/value/:metricType/:metricName", func(c *gin.Context) {
		metricsHandler.GetValue(c)
	})

	router.POST("/update/:metricType/:metricName/:metricValue", func(c *gin.Context) {
		metricsHandler.PostUpdate(c)
	})
	return router
}
