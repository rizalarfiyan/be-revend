package models

import (
	"github.com/rs/zerolog"
)

type Config struct {
	Env    string
	Port   int
	Host   string
	Logger LoggerConfigs
}

type LoggerConfigs struct {
	Level zerolog.Level
}
