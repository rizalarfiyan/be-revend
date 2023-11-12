package config

import (
	"os"
	"strconv"
	"time"

	"github.com/rizalarfiyan/be-revend/logger"
	"github.com/rs/zerolog"
)

type ConfigConvert interface {
	AsString(key string, defaultVal string) string
	AsInt(key string, defaultVal int) int
	AsBool(key string, defaultVal bool) bool
	AsTimeDuration(key string, defaultVal time.Duration) time.Duration
	AsZerologLevel(key string, defaultVal zerolog.Level) zerolog.Level
}

type configConvert struct {
	log logger.Logger
}

func NewConfigConvert() ConfigConvert {
	logs := logger.Get("config")
	return &configConvert{
		log: logs,
	}
}

func (cr configConvert) AsString(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	cr.log.Warnf("Environment variable %s not set", key)
	return defaultVal
}

func (cr configConvert) AsInt(key string, defaultVal int) int {
	valueStr := cr.AsString(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

func (cr configConvert) AsBool(key string, defaultVal bool) bool {
	valStr := cr.AsString(key, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

func (cr configConvert) AsTimeDuration(key string, defaultVal time.Duration) time.Duration {
	valStr := cr.AsString(key, "")
	val, err := time.ParseDuration(valStr)
	if err == nil {
		return val
	}

	return defaultVal
}

func (cr configConvert) AsZerologLevel(key string, defaultVal zerolog.Level) zerolog.Level {
	valStr := cr.AsString(key, "")
	if level, err := zerolog.ParseLevel(valStr); err == nil {
		return level
	}

	return defaultVal
}
