package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func GinMetricsMiddleware(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		HTTPRequestsInFlight.WithLabelValues(serviceName).Inc()
		defer HTTPRequestsInFlight.WithLabelValues(serviceName).Dec()

		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method
		path := c.FullPath()

		if path == "" {
			path = "unknown"
		}

		HTTPRequestsTotal.WithLabelValues(
			serviceName,
			method,
			path,
			status,
		).Inc()

		HTTPRequestDuration.WithLabelValues(
			serviceName,
			method,
			path,
			status,
		).Observe(duration)

		responseSize := c.Writer.Size()
		if responseSize < 0 {
			responseSize = 0
		}

		HTTPResponseSize.WithLabelValues(
			serviceName,
			method,
			path,
		).Observe(float64(responseSize))
	}
}

func MetricsHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func RegisterMetricsEndpoint(r *gin.Engine) {
	r.GET("/metrics", MetricsHandler())
}
