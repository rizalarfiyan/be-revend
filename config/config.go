package config

import (
	"log"
	"time"

	"github.com/rizalarfiyan/be-revend/models"
	"github.com/rs/zerolog"

	"github.com/joho/godotenv"
)

var conf *models.Config

func Init() {
	cc := NewConfigConvert()
	appEnv := cc.AsString("APP_ENV", "development")

	err := godotenv.Load(".env")
	if err != nil && appEnv != "production" {
		log.Fatal(".env is not loaded properly. Err: ", err)
	}

	conf = new(models.Config)
	conf.Host = cc.AsString("HOST", "")
	conf.Port = cc.AsInt("PORT", 8080)

	conf.Logger.Level = cc.AsZerologLevel("LOG_LEVEL", zerolog.InfoLevel)
	conf.Logger.Path = cc.AsString("LOG_PATH", "./log/logger.log")
	conf.Logger.IsCompressed = cc.AsBool("LOG_COMPRESSED", true)
	conf.Logger.IsDailyRotate = cc.AsBool("LOG_DAILY_ROTATE", true)
	conf.Logger.SleepDuration = cc.AsTimeDuration("LOG_SLEEP_DURATION", 5*time.Second)

	conf.Swagger.Username = cc.AsString("SWAGGER_USERNAME", "admin")
	conf.Swagger.Password = cc.AsString("SWAGGER_PASSWORD", "password")
}

func Get() *models.Config {
	return conf
}
