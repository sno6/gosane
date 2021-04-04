package middleware

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sno6/gosane/internal/http"
	"github.com/sno6/gosane/service/auth"
)

const UserContext = "user_ctx"

var (
	ErrInvalidAuthRequest = errors.New("invalid authentication request")
	ErrEmptyUserContext   = errors.New("user is not found in context")
)

type AuthDependencies struct {
	AuthService *auth.Service
	Logger      *log.Logger
}

func Auth(deps *AuthDependencies, whitelist []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(whitelist) > 0 && isRequestAllowedAuthless(c.FullPath(), whitelist) {
			c.Next()
			return
		}

		token, err := getAuthorizationTokenFromHeader(c.Request.Header)
		if err != nil {
			c.Error(http.Unauthorized).SetMeta(err)
			c.Abort()
			return
		}

		user, err := deps.AuthService.FindUserByToken(c, token)
		if err != nil {
			deps.Logger.Println("Unable to find user by provided token in authorization header")
			c.Error(http.Unauthorized).SetMeta(ErrInvalidAuthRequest)
			c.Abort()
			return
		}

		c.Set(UserContext, user)
		c.Next()
	}
}

func isRequestAllowedAuthless(p string, whitelist []string) bool {
	for _, v := range whitelist {
		if v == p {
			return true
		}
	}
	return false
}
