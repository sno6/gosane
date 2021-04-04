package auth

import (
	"errors"
	"net/http"

	gologin "github.com/dghubble/gologin/v2"
	intHttp "github.com/sno6/gosane/internal/http"

	"github.com/dghubble/gologin/v2/facebook"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type FacebookLoginHandler struct {
	cfg *oauth2.Config
}

func NewFacebookLoginHandler(cfg *oauth2.Config) *FacebookLoginHandler {
	return &FacebookLoginHandler{
		cfg: cfg,
	}
}

func (*FacebookLoginHandler) Path() string {
	return "/facebook/login"
}

func (*FacebookLoginHandler) Method() string {
	return http.MethodGet
}

func (fh *FacebookLoginHandler) HandleFunc(c *gin.Context) {
	handler := facebook.StateHandler(gologin.DebugOnlyCookieConfig, facebook.LoginHandler(fh.cfg, fh.failure(c)))
	handler.ServeHTTP(c.Writer, c.Request)
}

func (fh *FacebookLoginHandler) failure(c *gin.Context) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Error(intHttp.Forbidden).SetMeta(errors.New("facebook login failure"))
	})
}
