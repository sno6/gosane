package server

import (
	"log"
	"os"
	"time"

	"github.com/sno6/gosane/internal/email/ses"
	"github.com/sno6/gosane/internal/prometheus"
	"github.com/sno6/gosane/internal/verification"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sno6/gosane/api"
	"github.com/sno6/gosane/internal/database"
	"github.com/sno6/gosane/internal/jwt"
	"github.com/sno6/gosane/internal/sentry"
	"github.com/sno6/gosane/internal/validator"
	"github.com/sno6/gosane/middleware"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"

	appCfg "github.com/sno6/gosane/config"

	authService "github.com/sno6/gosane/service/auth"
	userService "github.com/sno6/gosane/service/user"

	tokenStore "github.com/sno6/gosane/store/token"
	userStore "github.com/sno6/gosane/store/user"
)

type Server struct {
	Engine   *gin.Engine
	Database *database.Database
	Config   appCfg.AppConfig
}

func New(cfg appCfg.AppConfig, env string) (*Server, error) {
	db, err := database.New(&database.Config{
		Name:    cfg.PostgresDB,
		Host:    cfg.PostgresHost,
		Port:    cfg.PostgresPort,
		User:    cfg.PostgresUser,
		Pass:    cfg.PostgresPassword,
		LogMode: cfg.LogMode,
	})
	if err != nil {
		return nil, errors.Wrap(err, "error setting up database")
	}

	engine, err := initEngine(db, cfg, env)
	if err != nil {
		return nil, err
	}

	return &Server{
		Config:   cfg,
		Engine:   engine,
		Database: db,
	}, nil
}

func (s *Server) Run(addr string) error {
	return s.Engine.Run(addr)
}

func initEngine(db *database.Database, cfg appCfg.AppConfig, env string) (*gin.Engine, error) {
	// Note: Remove the following 8 lines if you don't want Sentry error logging.
	sentryClient, err := sentry.New(cfg.SentryDSN)
	if err != nil {
		return nil, errors.Wrap(err, "unable to initialise sentry")
	}

	defer func() {
		sentryClient.Flush(time.Second * 5)
	}()

	// Internal services.
	emailer, err := ses.New()
	if err != nil {
		return nil, errors.Wrap(err, "initEngine: error initialising server")
	}

	prometheus, err := prometheus.New()
	if err != nil {
		return nil, errors.Wrap(err, "initEngine: error initialising prometheus")
	}

	logger := log.New(os.Stdout, "[Gosane] ", log.LstdFlags)
	jwtAuth := jwt.New([]byte(cfg.JWTSecret))
	verification := verification.New(cfg, emailer, jwtAuth)
	validator := validator.New()

	// Application Stores.
	userStore := userStore.NewUserStore(db.Client)
	tokenStore := tokenStore.NewTokenStore(db.Client)

	// Application services.
	userService := userService.NewUserService(userStore)
	authService := authService.NewAuthService(jwtAuth, tokenStore, userService, verification)

	fbCfg := &oauth2.Config{
		ClientID:     cfg.FacebookOAuthAppID,
		ClientSecret: cfg.FacebookOAuthAppSecret,
		RedirectURL:  cfg.FacebookOAuthCallbackURL,
		Endpoint:     facebook.Endpoint,
		Scopes:       []string{"email"},
	}

	googleCfg := &oauth2.Config{
		ClientID:     cfg.GoogleOAuthAppID,
		ClientSecret: cfg.GoogleOAuthAppSecret,
		RedirectURL:  cfg.GoogleOAuthCallbackURL,
		Endpoint:     google.Endpoint,
		Scopes:       []string{"profile", "email"},
	}

	// Set up engine, add custom middleware.
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(middleware.Cors())
	engine.Use(middleware.RequestMetrics(prometheus))
	engine.Use(middleware.Recovery(sentryClient, logger))
	engine.Use(middleware.Errors(sentryClient))

	api.Register(&api.Dependencies{
		Engine:       engine,
		Logger:       logger,
		AppConfig:    cfg,
		Validator:    validator,
		AuthService:  authService,
		UserService:  userService,
		FBConfig:     fbCfg,
		GoogleConfig: googleCfg,
	})

	return engine, nil
}
