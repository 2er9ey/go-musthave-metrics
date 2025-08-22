package handler

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/2er9ey/go-musthave-metrics/internal/logger"
	"github.com/2er9ey/go-musthave-metrics/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MetricHandler struct {
	service service.MetricServiceInterface
}

func NewMetricHandler(service service.MetricServiceInterface) *MetricHandler {
	return &MetricHandler{service: service}
}

type MetricRequest struct {
	ID    string `json:"id"`
	MType string `json:"type"`
	Value string `json:"value,omitempty"`
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

func (mh *MetricHandler) PostUpdateJSON(c *gin.Context) {
	if c.Request.Header.Get("Content-type") != "application/json" {
		logger.Log.Debug("got request with bad method", zap.String("method", c.Request.Method))
		c.String(http.StatusMethodNotAllowed, "Неверный тип данных")
		return
	}

	var req MetricRequest
	dec := json.NewDecoder(c.Request.Body)
	if err := dec.Decode(&req); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	mType := req.MType
	mName := req.ID
	mValue := req.Value

	if mType == "" || mName == "" || mValue == "" {
		c.String(http.StatusNotFound, "Неверный запрос {%s}, {%s}, {%s}", mType, mName, mValue)
		return
	}

	err := mh.service.Set(mName, mType, mValue)
	if err != nil {
		c.String(http.StatusBadRequest, "Неверное значение метрики")
		return
	}
	c.Header("Content-type", "application/json")
	c.String(http.StatusOK, "{}")
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

func (mh *MetricHandler) GetValueJSON(c *gin.Context) {
	if c.Request.Header.Get("Content-type") != "application/json" {
		logger.Log.Debug("got request with bad method", zap.String("method", c.Request.Method))
		c.String(http.StatusMethodNotAllowed, "Неверный тип данных")
		return
	}

	var req MetricRequest
	dec := json.NewDecoder(c.Request.Body)
	if err := dec.Decode(&req); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	mType := req.MType
	mName := req.ID

	if mType == "" || mName == "" {
		c.String(http.StatusNotFound, "Неверный запрос {%s}, {%s}", mType, mName)
		return
	}

	metric, err := mh.service.Get(mName, mType)
	c.Header("Content-type", "application/json")
	if err == nil {
		req.MType = mType
		req.ID = mName
		req.Value = metric

		enc := json.NewEncoder(c.Writer)
		if err := enc.Encode(req); err != nil {
			logger.Log.Debug("error encoding response", zap.Error(err))
			return
		}
		logger.Log.Debug("sending HTTP 200 response")
	} else {
		c.String(http.StatusNotFound, "{}")
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
	c.Writer.Header().Set("Content-type", "text/html; charset=utf-8")
	c.String(http.StatusOK, body)
}
