package auth

import (
	"errors"
	"fmt"
	"net/http"

	gologin "github.com/dghubble/gologin/v2"
	"github.com/dghubble/gologin/v2/google"
	"github.com/gin-gonic/gin"
	"github.com/sno6/gosane/config"
	"github.com/sno6/gosane/ent"
	"github.com/sno6/gosane/ent/schema"
	intHttp "github.com/sno6/gosane/internal/http"
	"github.com/sno6/gosane/service/auth"
	"github.com/sno6/gosane/service/user"
	"golang.org/x/oauth2"
)

type GoogleCallbackHandler struct {
	userService *user.Service
	authService *auth.Service
	appConfig   config.AppConfig
	oAuthCfg    *oauth2.Config
}

func NewGoogleCallbackHandler(userService *user.Service, authService *auth.Service, appConfig config.AppConfig, oAuthCfg *oauth2.Config) *GoogleCallbackHandler {
	return &GoogleCallbackHandler{
		userService: userService,
		authService: authService,
		appConfig:   appConfig,
		oAuthCfg:    oAuthCfg,
	}
}

func (*GoogleCallbackHandler) Path() string {
	return "/google/callback"
}

func (*GoogleCallbackHandler) Method() string {
	return http.MethodGet
}

func (gh *GoogleCallbackHandler) HandleFunc(c *gin.Context) {
	handler := google.StateHandler(gologin.DebugOnlyCookieConfig, google.CallbackHandler(gh.oAuthCfg, gh.success(c), gh.failure(c)))
	handler.ServeHTTP(c.Writer, c.Request)
}

func (gh *GoogleCallbackHandler) success(c *gin.Context) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		googleUser, err := google.UserFromContext(r.Context())
		if err != nil {
			c.Error(intHttp.Internal).SetMeta(err)
			return
		}

		u, err := gh.userService.FindByEmail(r.Context(), googleUser.Email)
		if err != nil {
			if ent.IsNotFound(err) {
				isVerified := googleUser.VerifiedEmail != nil && *googleUser.VerifiedEmail == true

				u, err = gh.authService.Register(r.Context(), &ent.User{
					Email:         googleUser.Email,
					EmailVerified: isVerified,
					ProviderID:    googleUser.Id,
					ProviderType:  schema.GoogleProvider,
					FirstName:     googleUser.GivenName,
					LastName:      googleUser.FamilyName,
				})
				if err != nil {
					c.Error(intHttp.Internal).SetMeta(err)
					return
				}
			} else {
				c.Error(intHttp.Internal).SetMeta(err)
				return
			}
		}

		tokens, err := gh.authService.CreateTokens(c, u)
		if err != nil {
			c.Error(intHttp.Internal).SetMeta(err)
			return
		}

		redirectURL := fmt.Sprintf("%s?token=%s&refresh=%s", gh.appConfig.OAuthSuccessRedirect, tokens.Access, tokens.Refresh)
		http.Redirect(w, r, redirectURL, 302)
	})
}

func (gh *GoogleCallbackHandler) failure(c *gin.Context) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Error(intHttp.Internal).SetMeta(errors.New("google callback failure"))
	})
}
