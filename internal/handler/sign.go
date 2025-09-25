package handler

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/2er9ey/go-musthave-metrics/internal/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Signer struct {
	key string
}

func NewSigner(k string) *Signer {
	return &Signer{key: k}
}

func (s *Signer) SignMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		contentSign := c.Request.Header.Get("HashSHA256")
		if len(contentSign) > 0 {
			// decodedBytes, err := hex.DecodeString(contentSign)
			// if err != nil {
			// 	c.Writer.WriteHeader(http.StatusInternalServerError)
			// 	return
			// }
			logger.Log.Debug("Request with sign", zap.String("signature", contentSign))
			cr, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, "{}")
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(cr))

			h := hmac.New(sha256.New, []byte(s.key))
			h.Write(cr)
			dst := hex.EncodeToString(h.Sum(nil))
			logger.Log.Debug("Calculates signature: ", zap.String("signature", dst))
			if contentSign != dst {
				logger.Log.Debug("Signature does not equal.")
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		} else {
			logger.Log.Debug("Request without sign")
		}

		originalWriter := c.Writer
		bw := newBufferWriter(c.Writer)
		c.Writer = bw

		c.Next()

		h2 := hmac.New(sha256.New, []byte(s.key))
		h2.Write(bw.buffer.Bytes())
		dst2 := hex.EncodeToString(h2.Sum(nil))

		logger.Log.Debug("Signing reply", zap.String("signature", dst2))

		originalWriter.Header().Add("HashSHA256", dst2)
		originalWriter.Write(bw.buffer.Bytes())
	}
}
