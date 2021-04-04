package oauth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/sno6/gosane/config"
	intHttp "github.com/sno6/gosane/internal/http"

	gologin "github.com/dghubble/gologin/v2"
	"github.com/dghubble/gologin/v2/facebook"
	"github.com/gin-gonic/gin"
	"github.com/sno6/gosane/ent"
	"github.com/sno6/gosane/ent/schema"
	"github.com/sno6/gosane/service/auth"
	"github.com/sno6/gosane/service/user"
	"golang.org/x/oauth2"
)

type FacebookCallbackHandler struct {
	userService *user.Service
	authService *auth.Service
	appConfig   config.AppConfig
	oAuthCfg    *oauth2.Config
}

func NewFacebookCallbackHandler(userService *user.Service, authService *auth.Service, appConfig config.AppConfig, oAuthCfg *oauth2.Config) *FacebookCallbackHandler {
	return &FacebookCallbackHandler{
		userService: userService,
		authService: authService,
		appConfig:   appConfig,
		oAuthCfg:    oAuthCfg,
	}
}

func (*FacebookCallbackHandler) Path() string {
	return "/facebook/callback"
}

func (*FacebookCallbackHandler) Method() string {
	return http.MethodGet
}

func (fh *FacebookCallbackHandler) HandleFunc(c *gin.Context) {
	handler := facebook.StateHandler(gologin.DebugOnlyCookieConfig, facebook.CallbackHandler(fh.oAuthCfg, fh.success(c), fh.failure(c)))
	handler.ServeHTTP(c.Writer, c.Request)
}

func (fh *FacebookCallbackHandler) success(c *gin.Context) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fbUser, err := facebook.UserFromContext(r.Context())
		if err != nil {
			c.Error(intHttp.Internal).SetMeta(err)
			return
		}

		names := strings.Split(fbUser.Name, " ")
		if len(names) < 2 {
			c.Error(intHttp.Internal).SetMeta(errors.New("partial user information"))
			return
		}

		u, err := fh.userService.FindByEmail(r.Context(), fbUser.Email)
		if err != nil {
			if ent.IsNotFound(err) {
				u, err = fh.authService.Register(r.Context(), &ent.User{
					Email:         fbUser.Email,
					ProviderID:    fbUser.ID,
					ProviderType:  schema.FacebookProvider,
					EmailVerified: true,
					FirstName:     names[0],
					LastName:      strings.Join(names[1:], " "),
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

		tokens, err := fh.authService.CreateTokens(c, u)
		if err != nil {
			c.Error(intHttp.Internal).SetMeta(err)
			return
		}

		redirectURL := fmt.Sprintf("%s?token=%s&refresh=%s", fh.appConfig.OAuthSuccessRedirect, tokens.Access, tokens.Refresh)
		http.Redirect(w, r, redirectURL, 302)
	})
}

func (fh *FacebookCallbackHandler) failure(c *gin.Context) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Error(intHttp.Internal).SetMeta(errors.New("facebook callback failure"))
	})
}
