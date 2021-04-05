package verification

import (
	"fmt"
	"time"

	"github.com/sno6/gosane/config"
	"github.com/sno6/gosane/internal/email"
	"github.com/sno6/gosane/internal/jwt"
)

const (
	// Note: This is the name of the email template you have set up
	// on your email provider.
	verificationTemplate = "email_verification"

	// All verification emails expire in a week.
	verificationEmailExpire = time.Hour * 24 * 7
)

type VerificationEmailData struct {
	VerificationURL string `json:"verification_url"`
}

type Verification struct {
	cfg     config.AppConfig
	emailer email.Emailer
	jwt     *jwt.Auth
}

func New(
	cfg config.AppConfig,
	emailer email.Emailer,
	jwt *jwt.Auth,
) *Verification {
	return &Verification{
		cfg:     cfg,
		emailer: emailer,
		jwt:     jwt,
	}
}

func (v *Verification) SendVerificationEmail(toEmail string) error {
	tokens, err := v.jwt.NewTokens(toEmail)
	if err != nil {
		return err
	}

	// Note: In most cases, this will need to point to your frontend..
	// to do that just add a dashboard url parameter to your JSON config and
	// pass it through to this service when initialising in server.go.
	content := fmt.Sprintf(
		"http://yourfrontend.com/email/verify?token=%s&refresh=%s",
		tokens.Access,
		tokens.Refresh,
	)

	return v.emailer.SendRawEmail(toEmail, &email.RawEmailData{
		Subject: "Please verify your account",
		Content: content,
	})
}

func (v *Verification) VerifyToken(token string) (*jwt.Claims, error) {
	return v.jwt.ParseToken(token)
}
