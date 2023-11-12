package models

import (
	"time"

	"github.com/rs/zerolog"
)

type Config struct {
	Env     string
	Port    int
	Host    string
	Logger  LoggerConfigs
	Cors    CorsConfigs
	Swagger SwaggerConfigs
}

type LoggerConfigs struct {
	Level         zerolog.Level
	Path          string
	IsCompressed  bool
	IsDailyRotate bool
	SleepDuration time.Duration
}

type CorsConfigs struct {
	AllowOrigins     string
	AllowMethods     string
	AllowHeaders     string
	AllowCredentials bool
	ExposeHeaders    string
}

type SwaggerConfigs struct {
	Username string
	Password string
}
