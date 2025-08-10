package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/2er9ey/go-musthave-metrics/internal/repository"
	"github.com/2er9ey/go-musthave-metrics/internal/service"
	"github.com/gin-gonic/gin"
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
				contentType: "text/plain",
				statusCode:  404,
			},
			request: "/update",
		},
	}

	repo := repository.NewMemoryStorage()
	serv := service.NewMetricService(repo)
	metricsHandler := NewMetricHandler(serv)

	router := setupRouter(*metricsHandler)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			request.Header.Set("Content-Type", "text/plain")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)
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

func setupRouter(metricsHandler MetricHandler) *gin.Engine {
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
