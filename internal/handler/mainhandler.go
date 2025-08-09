package handler

import (
	"net/http"
	"sort"

	"github.com/2er9ey/go-musthave-metrics/internal/service"
)

type MetricHandler struct {
	service service.MetricServiceInterface
}

func NewMetricHandler(service service.MetricServiceInterface) *MetricHandler {
	return &MetricHandler{service: service}
}

func (mh *MetricHandler) StatusBadRequest(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "", http.StatusBadRequest)
}

func (mh *MetricHandler) StatusNotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "", http.StatusNotFound)
}

func (mh *MetricHandler) PostUpdate(w http.ResponseWriter, r *http.Request) {
	// if r.Header.Get("Content-type") != "text/plain" {
	// 	http.Error(w, "Неверный запрос", http.StatusNotFound)
	// 	return
	// }

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
	w.Header().Set("Content-type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func (mh *MetricHandler) GetValue(w http.ResponseWriter, r *http.Request) {
	// if r.Header.Get("Content-type") != "text/plain" {
	// 	http.Error(w, "Неверный запрос", http.StatusNotFound)
	// 	return
	// }

	mType := r.PathValue("metricType")
	mName := r.PathValue("metricName")

	if mType == "" || mName == "" {
		http.Error(w, "Неверный запрос {"+mType+"}, {"+mName+"}", http.StatusNotFound)
		return
	}

	metric, err := mh.service.Get(mName, mType)
	w.Header().Set("Content-type", "text/plain")
	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(metric))
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(metric))
	}
}

func (mh *MetricHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// if r.Header.Get("Content-type") != "text/html" {
	// 	http.Error(w, "Неверный запрос", http.StatusNotFound)
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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}
