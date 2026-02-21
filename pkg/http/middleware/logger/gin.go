package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/logger"
	"github.com/gin-gonic/gin"
)

var sensitiveFields = map[string]bool{
	"password": true,
	"email":    true,
}

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func mask(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for k, val := range v {
			if sensitiveFields[k] {
				result[k] = "***"
			} else {
				result[k] = mask(val)
			}
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, val := range v {
			result[i] = mask(val)
		}
		return result
	default:
		return v
	}
}

func parseJSON(body []byte) interface{} {
	var data interface{}
	json.Unmarshal(body, &data)
	return data
}

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		var reqBody []byte
		if c.Request.Body != nil {
			reqBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}
		reqData := mask(parseJSON(reqBody))

		resBuf := &bytes.Buffer{}
		c.Writer = &responseWriter{body: resBuf, ResponseWriter: c.Writer}

		c.Next()

		resData := mask(parseJSON(resBuf.Bytes()))
		status := c.Writer.Status()
		log := logger.Get()
		latency := time.Since(start).Milliseconds()

		fields := []interface{}{
			"status", status,
			"method", c.Request.Method,
			"uri", c.Request.URL.Path,
			"latency_ms", latency,
			"ip", c.ClientIP(),
		}

		switch {
		case status >= 500:
			log.Error("http_request", fields...)
		case status >= 400:
			log.Warn("http_request", fields...)
		default:
			log.Info("http_request", fields...)
		}

		if reqData != nil || resData != nil {
			log.Debug("http_body", "request", reqData, "response", resData)
		}
	}
}
