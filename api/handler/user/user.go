package user

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sno6/gosane/api/handler"
	"github.com/sno6/gosane/internal/validator"
	"github.com/sno6/gosane/middleware"
	"github.com/sno6/gosane/service/auth"
	"github.com/sno6/gosane/service/user"
)

type UserHandler struct {
	validator   *validator.Validator
	logger      *log.Logger
	authService *auth.Service
	userService *user.Service
}

func New(
	validator *validator.Validator,
	logger *log.Logger,
	authService *auth.Service,
	userService *user.Service,
) *UserHandler {
	return &UserHandler{
		validator:   validator,
		logger:      logger,
		userService: userService,
		authService: authService,
	}
}

func (*UserHandler) RelativePath() string {
	return "/user"
}

func (u *UserHandler) Handlers() []handler.Handler {
	return []handler.Handler{
		NewGetHandler(u.userService),
		NewUpdateHandler(u.validator, u.userService),
		NewDeleteHandler(u.authService),
	}
}

func (u *UserHandler) MiddlewareChain() gin.HandlersChain {
	return gin.HandlersChain{
		middleware.Auth(&middleware.AuthDependencies{
			AuthService: u.authService,
			Logger:      u.logger,
		}, nil),
	}
}
