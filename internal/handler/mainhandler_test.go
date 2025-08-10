package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/2er9ey/go-musthave-metrics/internal/repository"
	"github.com/2er9ey/go-musthave-metrics/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainHandler(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name: "simple test #1",
			want: want{
				contentType: "text/plain",
				statusCode:  200,
			},
			request: "/update/gauge/xxx/1.234",
		},
		{
			name: "without metrit name",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
			},
			request: "/update",
		},
	}

	repo := repository.NewMemoryStorage()
	serv := service.NewMetricService(repo)
	mh := NewMetricHandler(serv)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /update/{metricType}/{metricName}/{metricValue}", mh.PostUpdate)
	mux.HandleFunc("POST /update/{metricType}/", mh.StatusNotFound)
	mux.HandleFunc("POST /update/{metricType}", mh.StatusNotFound)
	mux.HandleFunc("POST /update/", mh.StatusBadRequest)
	mux.HandleFunc("POST /update", mh.StatusBadRequest)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			request.Header.Set("Content-Type", "text/plain")
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-type"))

			_, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
		})
	}
}
