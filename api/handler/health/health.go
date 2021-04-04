package health

import (
	"github.com/sno6/gosane/api/handler"
)

type HealthHandler struct{}

func New() *HealthHandler {
	return &HealthHandler{}
}

func (*HealthHandler) RelativePath() string {
	return "/health"
}

func (*HealthHandler) Handlers() []handler.Handler {
	return []handler.Handler{
		NewGetHealthHandler(),
	}

}
