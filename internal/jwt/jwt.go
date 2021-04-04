package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/pkg/errors"

	lib "github.com/dgrijalva/jwt-go"
)

var ErrInvalidToken = errors.New("invalid claims")

var (
	accessTokenExpire  = time.Hour * 3
	refreshTokenExpire = time.Hour * 24
)

const RefreshTokenLength = 32

type Auth struct {
	secret []byte
}

type Claims struct {
	Identifier string
	lib.StandardClaims
}

type TokenInfo struct {
	Access           string    `json:"access,omitempty"`
	AccessExpiresAt  time.Time `json:"accessExpiresAt"`
	Refresh          string    `json:"refresh"`
	RefreshExpiresAt time.Time `json:"refreshExpiresAt"`
}

func New(secret []byte) *Auth {
	return &Auth{
		secret: secret,
	}
}

// NewTokens is used to create auth tokens (access & refresh) for API authentication.
func (a *Auth) NewTokens(identifier string) (*TokenInfo, error) {
	now := time.Now()
	accessExpiresAt := now.Add(accessTokenExpire)
	refreshExpiresAt := now.Add(refreshTokenExpire)

	claims := &Claims{
		Identifier: identifier,
		StandardClaims: lib.StandardClaims{
			ExpiresAt: accessExpiresAt.Unix(),
		},
	}

	accessToken, err := lib.NewWithClaims(lib.SigningMethodHS256, claims).SignedString(a.secret)
	if err != nil {
		return nil, err
	}

	refreshToken, err := a.generateRefreshToken()
	if err != nil {
		return nil, err
	}

	return &TokenInfo{
		Access:           accessToken,
		AccessExpiresAt:  accessExpiresAt,
		Refresh:          refreshToken,
		RefreshExpiresAt: refreshExpiresAt,
	}, nil
}

// NewAccessToken is a generic method to generate a JWT token with a given timeout.
func (a *Auth) NewAccessToken(identifier string, expiresIn time.Duration) (string, error) {
	expiresAt := time.Now().Add(expiresIn)

	claims := &Claims{
		Identifier: identifier,
		StandardClaims: lib.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
		},
	}

	return lib.NewWithClaims(lib.SigningMethodHS256, claims).SignedString(a.secret)
}

// ParseToken parses and validates a token is legitimate.
func (a *Auth) ParseToken(token string) (*Claims, error) {
	var claims Claims
	t, err := lib.ParseWithClaims(token, &claims, func(token *lib.Token) (interface{}, error) {
		return a.secret, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "error parsing token")
	}
	if !t.Valid {
		return nil, ErrInvalidToken
	}
	return &claims, nil
}

func (*Auth) generateRefreshToken() (string, error) {
	token := make([]byte, RefreshTokenLength)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(token), nil
}
