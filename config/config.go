package config

// Your server's configuration goes here..
//
// env:"MY_ENV"   tags for environment variables.
// json:"MY_JSON" tags for JSON file variables.
//
// Validation tags can be added, more info here: https://github.com/go-playground/validator
type AppConfig struct {
	Port                     int    `json:"http_port" validate:"required"`
	LogMode                  bool   `json:"log_mode"`
	SentryDSN                string `env:"SENTRY_DSN" validate:"required"`
	OAuthSuccessRedirect     string `json:"oauth_success_redirect" validate:"required"`
	FacebookOAuthAppID       string `json:"facebook_oauth_app_id" validate:"required"`
	FacebookOAuthAppSecret   string `env:"FACEBOOK_OAUTH_APP_SECRET" validate:"required"`
	FacebookOAuthCallbackURL string `json:"facebook_oauth_callback_url" validate:"required"`
	GoogleOAuthAppID         string `json:"google_oauth_app_id" validate:"required"`
	GoogleOAuthAppSecret     string `env:"GOOGLE_OAUTH_APP_SECRET" validate:"required"`
	GoogleOAuthCallbackURL   string `json:"google_oauth_callback_url" validate:"required"`
	JWTSecret                string `env:"JWT_SECRET" validate:"required"`
	PostgresDB               string `env:"POSTGRES_DB" validate:"required"`
	PostgresHost             string `env:"POSTGRES_HOST" validate:"required"`
	PostgresPassword         string `env:"POSTGRES_PASSWORD" validate:"required"`
	PostgresPort             string `env:"POSTGRES_PORT" validate:"required"`
	PostgresUser             string `env:"POSTGRES_USER" validate:"required"`
}
