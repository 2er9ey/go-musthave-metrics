package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/2er9ey/go-musthave-metrics/internal/logger"
	"github.com/2er9ey/go-musthave-metrics/internal/models"
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

// type MetricRequest struct {
// 	ID    string `json:"id"`
// 	MType string `json:"type"`
// 	Value string `json:"value,omitempty"`
// }

type MetricRequestBunch []models.Metrics

// UnmarshalJSON implements the json.Unmarshaler interface for Items.
func (i *MetricRequestBunch) UnmarshalJSON(data []byte) error {
	logger.Log.Debug("JSONRequest = ", zap.String("data", string(data)))
	// Try to unmarshal as an array first.
	var arr []models.Metrics
	if err := json.Unmarshal(data, &arr); err == nil {
		*i = arr
		return nil
	}

	// If it's not an array, try to unmarshal as a single object.
	var single models.Metrics
	if err := json.Unmarshal(data, &single); err == nil {
		*i = []models.Metrics{single} // Convert single item to a slice.
		return nil
	}

	return fmt.Errorf("cannot unmarshal JSON into Items: %s", string(data))
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
		logger.Log.Debug("cannot set metric", zap.Error(err))
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

	dec := json.NewDecoder(c.Request.Body)

	var req MetricRequestBunch
	if err := dec.Decode(&req); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err := mh.service.SetBunch(req)
	if err != nil {
		logger.Log.Debug("cannot set of metric", zap.Error(err))
		c.String(http.StatusBadRequest, "Неверное значение метрики")
		return
	}

	// var metrics models.Metrics
	// for _, item := range req {
	// 	err := mh.service.Set(item.ID, item.MType, item.Value)
	// 	if err != nil {
	// 		logger.Log.Debug("cannot set of metric", zap.Error(err))
	// 		c.String(http.StatusBadRequest, "Неверное значение метрики")
	// 		return
	// 	}
	// 	// 	// metrics = append(metrics, metric)
	// }

	// err := mh.service.SetBunch(metrics)
	// if err != nil {
	// 	logger.Log.Debug("cannot set bunch of metric", zap.Error(err))
	// 	c.String(http.StatusBadRequest, "Неверное значение метрики")
	// 	return
	// }
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

	var req models.Metrics
	dec := json.NewDecoder(c.Request.Body)
	if err := dec.Decode(&req); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if req.MType == "" || req.ID == "" {
		logger.Log.Debug("metric is incorrect", zap.String("req.String()", req.String()),
			zap.String("req.MType", req.MType), zap.String("req.ID", req.ID))
		c.String(http.StatusNotFound, "Неверный запрос {%s}, {%s}", req.MType, req.ID)
		return
	}

	metric, err := mh.service.GetMetric(req.ID, req.MType)
	c.Header("Content-type", "application/json")
	if err == nil {
		enc := json.NewEncoder(c.Writer)
		if err := enc.Encode(metric); err != nil {
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

func (mh *MetricHandler) Ping(c *gin.Context) {
	res, err := mh.service.Ping()
	// fmt.Println("Check result = ", res, err)
	if !res {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, "OK")
}
