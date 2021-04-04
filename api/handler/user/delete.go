package user

import (
	"github.com/gin-gonic/gin"
	"github.com/sno6/gosane/internal/http"
	"github.com/sno6/gosane/middleware"
	"github.com/sno6/gosane/service/auth"
)

type DeleteHandler struct {
	authService *auth.Service
}

func NewDeleteHandler(authService *auth.Service) *DeleteHandler {
	return &DeleteHandler{
		authService: authService,
	}
}

func (*DeleteHandler) Path() string {
	return ""
}

func (*DeleteHandler) Method() string {
	return http.MethodDelete
}

func (dh *DeleteHandler) HandleFunc(c *gin.Context) {
	u, err := middleware.UserFromContext(c)
	if err != nil {
		c.Error(http.Unauthorized).SetMeta(err)
		return
	}

	if err := dh.authService.DeleteUserByUuid(c, u.UUID); err != nil {
		c.Error(http.Internal).SetMeta(err)
	}
}
