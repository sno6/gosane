package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/sno6/gosane/api/handler/auth/dto/request"
	"github.com/sno6/gosane/internal/http"
	"github.com/sno6/gosane/internal/validator"
	"github.com/sno6/gosane/service/auth"
)

type VerifyEmailHandler struct {
	authService *auth.Service
	validator   *validator.Validator
}

func NewVerifyEmailHandler(
	authService *auth.Service,
	validator *validator.Validator,
) *VerifyEmailHandler {
	return &VerifyEmailHandler{
		authService: authService,
		validator:   validator,
	}
}

func (*VerifyEmailHandler) Path() string {
	return "/verify/email"
}

func (*VerifyEmailHandler) Method() string {
	return http.MethodPost
}

func (h *VerifyEmailHandler) HandleFunc(c *gin.Context) {
	var body request.VerifyEmailBody
	if err := h.validator.ValidateJSON(c.Request.Body, &body); err != nil {
		c.Error(http.BadRequest).SetMeta(err)
		return
	}

	err := h.authService.VerifyEmail(c, body.Token)
	if err != nil {
		c.Error(http.Unauthorized).SetMeta(err)
		return
	}

	c.JSON(200, nil)
}
