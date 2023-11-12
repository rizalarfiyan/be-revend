package models

import (
	"github.com/rs/zerolog"
)

type Config struct {
	Env    string
	Port   int
	Host   string
	Logger LoggerConfigs
	Cors   CorsConfigs
}

type LoggerConfigs struct {
	Level zerolog.Level
}

type CorsConfigs struct {
	AllowOrigins     string
	AllowMethods     string
	AllowHeaders     string
	AllowCredentials bool
	ExposeHeaders    string
}
