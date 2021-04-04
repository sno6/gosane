package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sno6/gosane/ent"
)

// UserFromContext is a helper method that casts the user context
// into a user struct, if it exists.
func UserFromContext(c *gin.Context) (*ent.User, error) {
	uCtx, exists := c.Get(UserContext)
	if !exists {
		return nil, ErrEmptyUserContext
	}

	u, ok := uCtx.(*ent.User)
	if !ok {
		return nil, ErrEmptyUserContext
	}

	return u, nil
}

func UserUUIDFromContext(c *gin.Context) (string, error) {
	u, err := UserFromContext(c)
	if err != nil {
		return "", nil
	}

	return u.UUID.String(), nil
}

func getAuthorizationTokenFromHeader(header http.Header) (string, error) {
	authHeader := header.Get("Authorization")
	if authHeader == "" {
		return "", ErrInvalidAuthRequest
	}

	spl := strings.Split(authHeader, " ")
	if len(spl) != 2 || spl[0] != "Bearer" || spl[1] == "" {
		return "", ErrInvalidAuthRequest
	}

	return spl[1], nil
}
