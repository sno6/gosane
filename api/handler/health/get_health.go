package health

import (
	"github.com/gin-gonic/gin"
	"github.com/sno6/gosane/internal/http"
)

type GetHealthHandler struct{}

func NewGetHealthHandler() *GetHealthHandler {
	return &GetHealthHandler{}
}

func (*GetHealthHandler) Path() string {
	return ""
}

func (*GetHealthHandler) Method() string {
	return http.MethodGet
}

func (*GetHealthHandler) HandleFunc(c *gin.Context) {
	c.Status(200)
}
