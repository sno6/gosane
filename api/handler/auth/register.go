package auth

import (
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sno6/gosane/api/handler/auth/dto/request"
	"github.com/sno6/gosane/ent"
	"github.com/sno6/gosane/internal/http"
	"github.com/sno6/gosane/internal/validator"
	"github.com/sno6/gosane/service/auth"
	"github.com/sno6/gosane/service/user"
)

type RegisterHandler struct {
	userService *user.Service
	authService *auth.Service
	validator   *validator.Validator
}

func NewRegisterHandler(
	userService *user.Service,
	authService *auth.Service,
	validator *validator.Validator,
) *RegisterHandler {
	return &RegisterHandler{
		userService: userService,
		authService: authService,
		validator:   validator,
	}
}

func (*RegisterHandler) Path() string {
	return "/register"
}

func (*RegisterHandler) Method() string {
	return http.MethodPost
}

func (h *RegisterHandler) HandleFunc(c *gin.Context) {
	var body request.RegisterBody
	if err := h.validator.ValidateJSON(c.Request.Body, &body); err != nil {
		c.Error(http.BadRequest).SetMeta(err)
		return
	}

	existing, err := h.userService.FindByEmail(c, body.Email)
	if err != nil && !ent.IsNotFound(err) {
		c.Error(http.Unauthorized).SetMeta(err)
		return
	}

	// A user already exists by that email.
	// Sleep for a random amount of time so malicious users can't aggregate user emails.
	if existing != nil {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
		c.JSON(200, nil)
		return
	}

	_, err = h.authService.Register(c, &ent.User{
		Email:         body.Email,
		EmailVerified: false,
		Password:      body.Password,
		FirstName:     body.FirstName,
		LastName:      body.LastName,
	})
	if err != nil {
		c.Error(http.Unauthorized).SetMeta(err)
		return
	}

	c.JSON(200, nil)
}
