package cmd

import (
	"log"
	"os"
	"strconv"

	appConfig "github.com/sno6/gosane/config"
	"github.com/sno6/gosane/internal/config"

	"github.com/joho/godotenv"
	"github.com/sno6/gosane/internal/database"
	"github.com/sno6/gosane/internal/server"
	"github.com/spf13/cobra"
)

const (
	defaultConfigPath = "./config"
	appEnvVariable    = "APP_ENV"
)

func Run() error {
	rootCmd := &cobra.Command{
		Use:   "gosane",
		Short: "Run the Gosane API server.",
		Run: func(cmd *cobra.Command, args []string) {
			env := os.Getenv(appEnvVariable)
			if env == config.Local.String() {
				err := godotenv.Load()
				if err != nil {
					log.Fatalf("Error loading .env file: %v\n", err)
				}
			}

			var cfg appConfig.AppConfig
			err := config.NewFromFile(&config.Config{
				Path:        defaultConfigPath,
				Environment: config.EnvironmentFromString(env),
			}, &cfg)

			if err != nil {
				log.Fatalf("Unable to load config: %v\n", err)
			}

			db, err := database.New(&database.Config{
				Name:    cfg.PostgresDB,
				Host:    cfg.PostgresHost,
				Port:    cfg.PostgresPort,
				User:    cfg.PostgresUser,
				Pass:    cfg.PostgresPassword,
				LogMode: cfg.LogMode,
			})
			if err != nil {
				log.Fatalf("Error setting up database: %v\n", err)
			}

			srv, err := server.New(db, cfg, env)
			if err != nil {
				log.Fatalf("Error setting up server: %v\n", err)
			}

			srv.Run(":" + strconv.Itoa(cfg.Port))
		},
	}

	return rootCmd.Execute()
}
