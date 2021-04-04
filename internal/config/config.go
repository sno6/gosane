package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	env "github.com/Netflix/go-env"
	"github.com/go-playground/validator/v10"
)

const defaultCfgFile = "config.local.json"

type Config struct {
	Path        string
	Environment Environment
}

func NewFromFile(cfg *Config, s interface{}) error {
	file, err := getFile(cfg)
	if err != nil {
		return err
	}
	defer file.Close()

	err = readJSONIntoStruct(file, s)
	if err != nil {
		return err
	}

	_, err = env.UnmarshalFromEnviron(s)
	if err != nil {
		return err
	}

	validator := validator.New()
	err = validator.Struct(s)
	if err != nil {
		return err
	}

	return nil
}

func getFile(cfg *Config) (io.ReadCloser, error) {
	stat, err := os.Stat(cfg.Path)
	if err != nil {
		return nil, err
	}

	if !stat.IsDir() {
		cfg.Path, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	var cfgPath string
	if cfg.Environment.String() == "" {
		cfgPath = path.Join(cfg.Path, defaultCfgFile)
	} else {
		cfgPath = path.Join(cfg.Path, fmt.Sprintf("config.%s.json", cfg.Environment.String()))
	}

	return os.Open(cfgPath)
}

func readJSONIntoStruct(file io.Reader, s interface{}) error {
	return json.NewDecoder(file).Decode(s)
}
