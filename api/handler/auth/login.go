package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/sno6/gosane/api/handler/auth/dto/request"
	"github.com/sno6/gosane/internal/http"
	"github.com/sno6/gosane/internal/validator"
	"github.com/sno6/gosane/service/auth"
)

type LoginHandler struct {
	authService *auth.Service
	validator   *validator.Validator
}

func NewLoginHandler(
	authService *auth.Service,
	validator *validator.Validator,
) *LoginHandler {
	return &LoginHandler{
		authService: authService,
		validator:   validator,
	}
}

func (*LoginHandler) Path() string {
	return "/login"
}

func (*LoginHandler) Method() string {
	return http.MethodPost
}

func (h *LoginHandler) HandleFunc(c *gin.Context) {
	var body request.LoginBody
	if err := h.validator.ValidateJSON(c.Request.Body, &body); err != nil {
		c.Error(http.BadRequest).SetMeta(err)
		return
	}

	tokens, err := h.authService.Login(c, body.Email, body.Password)
	if err != nil {
		c.Error(http.Unauthorized).SetMeta(err)
		return
	}

	c.JSON(200, tokens)
}
