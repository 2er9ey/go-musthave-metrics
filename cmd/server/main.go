package main

import (
	"flag"
	"io"

	"github.com/2er9ey/go-musthave-metrics/internal/handler"
	"github.com/2er9ey/go-musthave-metrics/internal/repository"
	"github.com/2er9ey/go-musthave-metrics/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	repo := repository.NewMemoryStorage()
	service := service.NewMetricService(repo)
	metricsHadler := handler.NewMetricHandler(service)
	listenEndpoint := flag.String("a", "localhost:8080", "Адрес и порт для работы севрера")

	flag.Parse()

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/", func(c *gin.Context) {
		metricsHadler.GetAll(c.Writer, c.Request)
	})

	router.GET("/value/:metricType/:metricName", func(c *gin.Context) {
		c.Request.SetPathValue("metricType", c.Param("metricType"))
		c.Request.SetPathValue("metricName", c.Param("metricName"))
		metricsHadler.GetValue(c.Writer, c.Request)
	})

	router.POST("/update/:metricType/:metricName/:metricValue", func(c *gin.Context) {
		c.Request.SetPathValue("metricType", c.Param("metricType"))
		c.Request.SetPathValue("metricName", c.Param("metricName"))
		c.Request.SetPathValue("metricValue", c.Param("metricValue"))
		metricsHadler.PostUpdate(c.Writer, c.Request)
	})

	router.Run(*listenEndpoint)
}
