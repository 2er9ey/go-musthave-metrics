package handler

import (
	"net/http"
)

type MetricServiceInterface interface {
	Set(string, string, string) error
	Get(mID string) (string, error)
	GetAll() map[string]string
}

type MetricHandler struct {
	service MetricServiceInterface
}

func NewMetricHandler(service MetricServiceInterface) *MetricHandler {
	return &MetricHandler{service: service}
}

func (mh *MetricHandler) StatusBadRequest(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "", http.StatusBadRequest)
}

func (mh *MetricHandler) StatusNotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "", http.StatusNotFound)
}

func (mh *MetricHandler) PostUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-type") != "text/plain" {
		http.Error(w, "Неверный запрос", http.StatusNotFound)
		return
	}

	mType := r.PathValue("metricType")
	mName := r.PathValue("metricName")
	mValue := r.PathValue("metricValue")

	if mType == "" || mName == "" || mValue == "" {
		http.Error(w, "Неверный запрос {"+mType+"}, {"+mName+"}, {"+mValue+"}", http.StatusNotFound)
		return
	}

	err := mh.service.Set(mName, mType, mValue)
	if err != nil {
		http.Error(w, "Неверное значение метрики", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
