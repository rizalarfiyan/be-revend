package config

import (
	"log"

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
}

func Get() *models.Config {
	return conf
}
