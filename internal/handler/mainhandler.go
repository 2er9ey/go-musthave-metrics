package handler

import (
	"net/http"
	"sort"

	"github.com/2er9ey/go-musthave-metrics/internal/service"
	"github.com/gin-gonic/gin"
)

type MetricHandler struct {
	service service.MetricServiceInterface
}

func NewMetricHandler(service service.MetricServiceInterface) *MetricHandler {
	return &MetricHandler{service: service}
}

func (mh *MetricHandler) PostUpdate(c *gin.Context) {
	// if c.Request.Header.Get("Content-type") != "text/plain" {
	// 	c.String(http.StatusNotFound, "Неверный тип данных")
	// 	return
	// }

	mType := c.Param("metricType")
	mName := c.Param("metricName")
	mValue := c.Param("metricValue")

	if mType == "" || mName == "" || mValue == "" {
		c.String(http.StatusNotFound, "Неверный запрос {%s}, {%s}, {%s}", mType, mName, mValue)
		return
	}

	err := mh.service.Set(mName, mType, mValue)
	if err != nil {
		c.String(http.StatusBadRequest, "Неверное значение метрики")
		return
	}
	c.Header("Content-type", "text/plain")
	c.String(http.StatusOK, "")
}

func (mh *MetricHandler) GetValue(c *gin.Context) {
	// if  c.Request.Header.Get("Content-type") != "text/plain" {
	// 	c.String(http.StatusNotFound, "Неверный тип данных")
	// 	return
	// }

	mType := c.Param("metricType")
	mName := c.Param("metricName")

	if mType == "" || mName == "" {
		c.String(http.StatusNotFound, "Неверный запрос {%s}, {%s}", mType, mName)
		return
	}

	metric, err := mh.service.Get(mName, mType)
	c.Header("Content-type", "text/plain")
	if err == nil {
		c.String(http.StatusOK, metric)
	} else {
		c.String(http.StatusNotFound, metric)
	}
}

func (mh *MetricHandler) GetAll(c *gin.Context) {
	// if c.Request.Header.Get("Content-type") != "text/html" {
	//  c.String(http.StatusNotFound, "Неверный тип данных")
	// 	return
	// }
	metrics := mh.service.GetAll()
	body := "<html><head><title>Список известных метрик></title></head><body><table><tr><th>Имя</th><th>Тип</th><th>Значение</th></tr>"
	sort.Slice(metrics, func(i, j int) bool {
		return metrics[i].ID < metrics[j].ID
	})
	for _, metric := range metrics {
		body += "<tr><td>" + metric.ID + "</td><td>" + metric.MType + "</td><td align=right>" + metric.String() + "</td></td>"
	}
	body += "</table></body></html>"
	//	w.Header().Set("Content-type", "text/html; charset=utf-8")
	c.String(http.StatusOK, body)
}
