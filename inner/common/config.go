package common

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
)

type Config struct {
	DbDriverName   string `validate:"required"`
	Dsn            string `validate:"required"`
	AppName        string `validate:"required"`
	AppVersion     string `validate:"required"`
	LogLevel       string
	LogDevelopMode bool
	SslSert        string `validate:"required"`
	SslKey         string `validate:"required"`
}

func GetConfig(envFile string) Config {
	var err = godotenv.Load(envFile)
	if err != nil {
		log.Infof(fmt.Sprintf("Error loading .env file: %v\n", zap.Error(err)))
	}
	var cfg = Config{
		DbDriverName:   os.Getenv("DB_DRIVER_NAME"),
		Dsn:            os.Getenv("DB_DSN"),
		AppName:        os.Getenv("APP_NAME"),
		AppVersion:     os.Getenv("APP_VERSION"),
		LogLevel:       os.Getenv("LOG_LEVEL"),
		LogDevelopMode: os.Getenv("LOG_DEVELOP_MODE") == "true",
		SslSert:        os.Getenv("SSL_SERT"),
		SslKey:         os.Getenv("SSL_KEY"),
	}
	err = validator.New().Struct(cfg)
	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			panic(fmt.Sprintf("config validation error: %v", err))
		}
	}
	return cfg
}
