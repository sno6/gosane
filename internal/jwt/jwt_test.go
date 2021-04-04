package jwt

import (
	"testing"

	"github.com/pkg/errors"
)

var (
	secret   = []byte("test-secret")
	testUUID = "test-uuid"
)

func TestNewToken(t *testing.T) {
	auther := New(secret)

	tok, err := auther.NewTokens(testUUID)
	if err != nil {
		t.Errorf("error creating a token: %v", err)
	}
	if tok.Access == "" || tok.Refresh == "" {
		t.Error("generated token is empty")
	}
}

func TestParseToken(t *testing.T) {
	auther := New(secret)

	tok, err := auther.NewTokens(testUUID)
	if err != nil {
		t.Errorf("error creating a token: %v", err)
	}

	// Check a correctly parsed token.
	if err := checkToken(auther, tok.Access); err != nil {
		t.Error(err)
	}

	// Check an empty token.
	if err := checkToken(auther, ""); err == nil {
		t.Error("empty token was successfully parsed")
	}

	// Check a valid token signed with an incorrect signature.
	errAuther := New([]byte("another-secret"))
	if err := checkToken(errAuther, tok.Access); err == nil {
		t.Error("token passed with invalid signature")
	}
}

func checkToken(auther *Auth, tok string) error {
	claims, err := auther.ParseToken(tok)
	if err != nil {
		return err
	}
	if claims.Identifier != testUUID {
		return errors.New("incorrect token claims: wrong identifier")
	}
	if claims.StandardClaims.ExpiresAt == 0 {
		return errors.New("incorrent token claims: no expiry")
	}
	return nil
}
