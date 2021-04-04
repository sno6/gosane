package user

import (
	"github.com/gin-gonic/gin"
	"github.com/sno6/gosane/api/dto"
	"github.com/sno6/gosane/internal/http"
	"github.com/sno6/gosane/middleware"
	"github.com/sno6/gosane/service/user"
)

type GetHandler struct {
	userService *user.Service
}

func NewGetHandler(userService *user.Service) *GetHandler {
	return &GetHandler{
		userService: userService,
	}
}

func (*GetHandler) Path() string {
	return ""
}

func (*GetHandler) Method() string {
	return http.MethodGet
}

func (uh *GetHandler) HandleFunc(c *gin.Context) {
	u, err := middleware.UserFromContext(c)
	if err != nil {
		c.Error(http.Unauthorized).SetMeta(err)
		return
	}

	c.JSON(200, dto.NewFromUser(u))
}
