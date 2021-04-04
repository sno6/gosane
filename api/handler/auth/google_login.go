package auth

import (
	"errors"
	"net/http"

	gologin "github.com/dghubble/gologin/v2"
	intHttp "github.com/sno6/gosane/internal/http"

	"github.com/dghubble/gologin/v2/google"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type GoogleLoginHandler struct {
	cfg *oauth2.Config
}

func NewGoogleLoginHandler(cfg *oauth2.Config) *GoogleLoginHandler {
	return &GoogleLoginHandler{
		cfg: cfg,
	}
}

func (*GoogleLoginHandler) Path() string {
	return "/google/login"
}

func (*GoogleLoginHandler) Method() string {
	return http.MethodGet
}

func (gh *GoogleLoginHandler) HandleFunc(c *gin.Context) {
	handler := google.StateHandler(gologin.DebugOnlyCookieConfig, google.LoginHandler(gh.cfg, gh.failure(c)))
	handler.ServeHTTP(c.Writer, c.Request)
}

func (gh *GoogleLoginHandler) failure(c *gin.Context) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Error(intHttp.Internal).SetMeta(errors.New("google login failure"))
	})
}
