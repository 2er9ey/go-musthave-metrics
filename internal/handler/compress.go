package handler

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/2er9ey/go-musthave-metrics/internal/logger"
	"github.com/gin-gonic/gin"
)

// compressReader реализует интерфейс io.ReadCloser и позволяет прозрачно для сервера
// декомпрессировать получаемые от клиента данные
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

type bufferWriter struct {
	gin.ResponseWriter
	buffer *bytes.Buffer
}

func newBufferWriter(w gin.ResponseWriter) *bufferWriter {
	return &bufferWriter{
		ResponseWriter: w,
		buffer:         new(bytes.Buffer),
	}
}

func (c *bufferWriter) Header() http.Header {
	return c.ResponseWriter.Header()
}

func (c *bufferWriter) Write(p []byte) (int, error) {
	return c.buffer.Write(p)
}

func (c *bufferWriter) WriteHeader(statusCode int) {
	// if statusCode < 300 {
	// 	c.w.Header().Set("Content-Encoding", "gzip")
	// }
	c.ResponseWriter.WriteHeader(statusCode)
}

func GzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		contentEncoding := c.Request.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			logger.Log.Debug("Request with compression")
			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
			cr, err := newCompressReader(c.Request.Body)
			if err != nil {
				c.Writer.WriteHeader(http.StatusInternalServerError)
				return
			}
			// меняем тело запроса на новое
			c.Request.Body = cr
			defer cr.Close()
		} else {
			logger.Log.Debug("Request without compression")
		}

		acceptEncoding := c.Request.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			logger.Log.Debug("Accepts compression")
			originalWriter := c.Writer
			bw := newBufferWriter(c.Writer)
			c.Writer = bw
			c.Next()
			contentType := c.Writer.Header().Get("Content-Type")
			if strings.Contains(contentType, "text/html") || strings.Contains(contentType, "application/json") {
				logger.Log.Debug("Content should be compressed")
				originalWriter.Header().Add("content-encoding", "gzip")
				gz := gzip.NewWriter(originalWriter)
				gz.Write(bw.buffer.Bytes())
				defer gz.Close()
			} else {
				logger.Log.Debug("Content should not be compressed")
				c.Next()
			}
		} else {
			logger.Log.Debug("No accepter compression")
			c.Next()
		}
	}
}
