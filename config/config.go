package config

import (
	"errors"
	"fmt"
	"strings"

	db "github.com/gocomerse/internal/db"
	loggerModel "github.com/gocomerse/internal/logger/model"

	"github.com/jinzhu/configor"
	"github.com/joho/godotenv"
)

var ErrInvalidFileExtension = errors.New("file extension not supported")

type AppConfig struct {
	APPName string
	Env     string
	Server  struct {
		Port     int
		Host     string
		GRPCPort int
	}

	Logger   *loggerModel.Config
	Database *db.Config
}

func LoadConfig(fileNames ...string) (*AppConfig, error) {
	loadFiles := make([]string, 0, len(fileNames))
	envFiles := make([]string, 0, len(fileNames))

	for _, file := range fileNames {
		fileParts := strings.Split(file, ".")
		ext := fileParts[len(fileParts)-1]

		switch ext {
		case "yml", "json", "yaml", "toml":
			loadFiles = append(loadFiles, file)
		case "env":
			envFiles = append(envFiles, file)
		default:
			return nil, ErrInvalidFileExtension
		}

		if len(envFiles) > 0 {

			err := godotenv.Load(envFiles...)
			if err != nil {
				return nil, fmt.Errorf("error while loading env files(%s): %w", strings.Join(envFiles, ","), err)
			}
		}

	}
	_cfg, err := loadConfig(loadFiles...)
	if err != nil {
		return _cfg, err
	}
	return _cfg, err
}

func loadConfig(fileName ...string) (*AppConfig, error) {
	var appConfig AppConfig
	conf := newConf()
	if err := conf.Load(&appConfig, fileName...); err != nil {
		return nil, fmt.Errorf("failed to load config file: %w", err)
	}
	return &appConfig, nil
}

func newConf() *configor.Configor {
	conf := configor.Config{ENVPrefix: "GOCOMERSE"}
	config := configor.New(&conf)
	return config
}
