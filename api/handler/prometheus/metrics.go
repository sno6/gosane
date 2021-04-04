package prometheus

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsHandler struct {
}

func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{}
}

func (*MetricsHandler) Path() string {
	return "/metrics"
}

func (*MetricsHandler) Method() string {
	return http.MethodGet
}

func (h *MetricsHandler) HandleFunc(c *gin.Context) {
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}
