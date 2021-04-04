package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/sno6/gosane/api/handler/auth/dto/request"
	"github.com/sno6/gosane/internal/http"
	"github.com/sno6/gosane/internal/validator"
	"github.com/sno6/gosane/service/auth"
)

type TokenRefreshHandler struct {
	authService *auth.Service
	validator   *validator.Validator
}

func NewTokenRefreshHandler(
	authService *auth.Service,
	validator *validator.Validator,
) *TokenRefreshHandler {
	return &TokenRefreshHandler{
		authService: authService,
		validator:   validator,
	}
}

func (*TokenRefreshHandler) Path() string {
	return "/refresh"
}

func (*TokenRefreshHandler) Method() string {
	return http.MethodPut
}

func (th *TokenRefreshHandler) HandleFunc(c *gin.Context) {
	var body request.RefreshTokenBody
	if err := th.validator.ValidateJSON(c.Request.Body, &body); err != nil {
		c.Error(http.BadRequest).SetMeta(err)
		return
	}

	info, err := th.authService.Refresh(c, body.RefreshToken)
	if err != nil {
		c.Error(http.Unauthorized).SetMeta(err)
		return
	}

	c.JSON(200, info)
}
