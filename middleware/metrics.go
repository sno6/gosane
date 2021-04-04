package middleware

import (
	"time"

	"github.com/sno6/gosane/internal/prometheus"

	"github.com/gin-gonic/gin"
)

func RequestMetrics(prometheus *prometheus.Prometheus) gin.HandlerFunc {
	return func(c *gin.Context) {
		labels := getGenericRequestLabelsFromContext(c)

		now := time.Now()
		prometheus.HTTPRequestCounter.With(labels).Add(1)

		// Run the handler and any other middleware.
		c.Next()

		diff := float64(time.Since(now).Milliseconds())
		prometheus.HTTPResponseDuration.With(labels).Add(diff)
	}
}

func getGenericRequestLabelsFromContext(c *gin.Context) map[string]string {
	return map[string]string{
		prometheus.CtxNameHandler: c.FullPath(),
		prometheus.CtxNameMethod:  c.Request.Method,
	}
}
