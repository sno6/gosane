package prometheus

import (
	"github.com/sno6/gosane/api/handler"
)

type PrometheusHandler struct{}

func New() *PrometheusHandler {
	return &PrometheusHandler{}
}

func (*PrometheusHandler) RelativePath() string {
	return "/prometheus"
}

func (s *PrometheusHandler) Handlers() []handler.Handler {
	return []handler.Handler{
		NewMetricsHandler(),
	}
}
