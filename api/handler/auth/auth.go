package auth

import (
	"github.com/sno6/gosane/api/handler"
	"github.com/sno6/gosane/config"
	"github.com/sno6/gosane/internal/validator"
	"github.com/sno6/gosane/service/auth"
	"github.com/sno6/gosane/service/user"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	userService  *user.Service
	authService  *auth.Service
	validator    *validator.Validator
	appConfig    config.AppConfig
	fbConfig     *oauth2.Config
	googleConfig *oauth2.Config
}

func New(
	userService *user.Service,
	authService *auth.Service,
	validator *validator.Validator,
	appConfig config.AppConfig,
	fbConfig *oauth2.Config,
	googleConfig *oauth2.Config,
) *AuthHandler {
	return &AuthHandler{
		userService:  userService,
		authService:  authService,
		validator:    validator,
		appConfig:    appConfig,
		fbConfig:     fbConfig,
		googleConfig: googleConfig,
	}
}

func (*AuthHandler) RelativePath() string {
	return "/auth"
}

func (oh *AuthHandler) Handlers() []handler.Handler {
	return []handler.Handler{
		NewFacebookLoginHandler(oh.fbConfig),
		NewFacebookCallbackHandler(oh.userService, oh.authService, oh.appConfig, oh.fbConfig),
		NewGoogleLoginHandler(oh.googleConfig),
		NewGoogleCallbackHandler(oh.userService, oh.authService, oh.appConfig, oh.googleConfig),
		NewTokenRefreshHandler(oh.authService, oh.validator),
		NewRegisterHandler(oh.userService, oh.authService, oh.validator),
		NewVerifyEmailHandler(oh.authService, oh.validator),
		NewLoginHandler(oh.authService, oh.validator),
	}
}
