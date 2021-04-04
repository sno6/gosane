package database

import (
	"context"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/sno6/gosane/ent"
)

type Database struct {
	Client *ent.Client
	Config *Config
}

type Config struct {
	Name         string
	Host         string
	Port         string
	User         string
	Pass         string
	MigrationDir string
	LogMode      bool
}

func New(cfg *Config) (*Database, error) {
	client, err := initDB(cfg)
	if err != nil {
		return nil, err
	}
	return &Database{
		Client: client,
		Config: cfg,
	}, nil
}

func initDB(cfg *Config) (*ent.Client, error) {
	for _, s := range []string{cfg.Name, cfg.Host, cfg.Port, cfg.User, cfg.Pass} {
		if s == "" {
			return nil, errors.New("incomplete environment for database")
		}
	}

	connStr := fmt.Sprintf("dbname=%s host=%s port=%s user=%s password=%s sslmode=disable",
		cfg.Name, cfg.Host, cfg.Port, cfg.User, cfg.Pass)

	client, err := ent.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if cfg.LogMode {
		client = client.Debug()
	}

	// TODO (sno6): This will not be auto generated in the future.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	return client, nil
}
