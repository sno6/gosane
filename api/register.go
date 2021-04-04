package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sno6/gosane/api/handler"
	"github.com/sno6/gosane/api/handler/health"
	"github.com/sno6/gosane/api/handler/oauth"
	"github.com/sno6/gosane/config"
	"github.com/sno6/gosane/internal/email"
	"github.com/sno6/gosane/internal/sentry"
	"github.com/sno6/gosane/internal/validator"
	"golang.org/x/oauth2"

	prometheusHandler "github.com/sno6/gosane/api/handler/prometheus"
	userHandler "github.com/sno6/gosane/api/handler/user"

	authService "github.com/sno6/gosane/service/auth"
	userService "github.com/sno6/gosane/service/user"
)

// Dependencies are everything needed by handlers to function.
type Dependencies struct {
	AppConfig    config.AppConfig
	Engine       *gin.Engine
	Logger       *log.Logger
	Emailer      *email.Email
	Sentry       *sentry.Sentry
	Validator    *validator.Validator
	UserService  *userService.Service
	AuthService  *authService.Service
	FBConfig     *oauth2.Config
	GoogleConfig *oauth2.Config
}

// Register all handlers with their dependencies.
func Register(deps *Dependencies) {
	handlerGroups := []handler.HandlerGroup{
		prometheusHandler.New(),
		health.New(),
		userHandler.New(
			deps.Validator,
			deps.Logger,
			deps.AuthService,
			deps.UserService,
		),
		oauth.New(
			deps.UserService,
			deps.AuthService,
			deps.Validator,
			deps.AppConfig,
			deps.FBConfig,
			deps.GoogleConfig,
		),
	}

	for _, hg := range handlerGroups {
		group := deps.Engine.Group(hg.RelativePath())

		// Attach middleware if the handler group implements the middleware chainer interface.
		chainer, hasMiddleware := hg.(handler.MiddlewareChainer)
		if hasMiddleware {
			group.Use(chainer.MiddlewareChain()...)
		}

		for _, handler := range hg.Handlers() {
			group.Handle(handler.Method(), handler.Path(), handler.HandleFunc)
		}
	}
}
