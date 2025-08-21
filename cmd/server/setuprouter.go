package main

import (
	"io"
	"time"

	"github.com/2er9ey/go-musthave-metrics/internal/handler"
	"github.com/2er9ey/go-musthave-metrics/internal/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupRouter(metricsHandler handler.MetricHandler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	router := gin.New()
	//router := gin.Default()

	router.Use(gin.Recovery())
	router.Use(LoggerMiddleware())

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

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		//		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		//		c.Writer = blw
		c.Next() // Pass control to the next handler in the chain
		duration := time.Since(start)
		statusCode := c.Writer.Status()
		responseLen := c.Writer.Size()
		//		responseLen := len(blw.body.String())
		logger.Log.Info("Incoming Request:", zap.String("method", c.Request.Method), zap.String("URI", c.Request.URL.Path), zap.String("elapsedTime", duration.String()))
		logger.Log.Info("Outging reply:",
			zap.Int("statusCode", statusCode), zap.Int("responseLen", responseLen))
	}
}
