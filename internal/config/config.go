package config

import (
	"fmt"
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Config struct {
	RestServerPort string `env:"REST_SERVER_PORT" env-default:"8080"`

	PostgresUser     string `env:"POSTGRES_USER" env-default:"postgres"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresDB       string `env:"POSTGRES_DB"`
	PostgresHost     string `env:"POSTGRES_HOST" env-default:"postgres"`
	PostgresPort     string `env:"POSTGRES_PORT" env-default:"5432"`
}

func New(filename string) (*Config, error) {
	cfg := Config{}
	err := cleanenv.ReadConfig(fmt.Sprintf("./config/%s", filename), &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func LoggerConfig() middleware.RequestLoggerConfig {
	return middleware.RequestLoggerConfig{
		LogURI:     true,
		LogStatus:  true,
		LogMethod:  true,
		LogLatency: true,
		LogError:   true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			slog.Info("request",
				"method", v.Method,
				"uri", v.URI,
				"status", v.Status,
				"latency", v.Latency,
				"error", v.Error,
			)
			return nil
		},
	}
}
