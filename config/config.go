package config

import (
	"log"
	"time"

	"github.com/rizalarfiyan/be-revend/models"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2/google"

	"github.com/joho/godotenv"
)

var conf *models.Config

func Init() {
	cc := NewConfigConvert()
	appEnv := cc.AsString("ENV", "development")

	err := godotenv.Load(".env")
	if err != nil && appEnv != "production" {
		log.Fatal(".env is not loaded properly. Err: ", err)
	}

	conf = new(models.Config)
	conf.Host = cc.AsString("HOST", "")
	conf.Port = cc.AsInt("PORT", 8080)
	conf.Name = cc.AsString("APP_NAME", "Revend")

	conf.Logger.Level = cc.AsZerologLevel("LOG_LEVEL", zerolog.InfoLevel)
	conf.Logger.Path = cc.AsString("LOG_PATH", "./log/logger.log")
	conf.Logger.IsCompressed = cc.AsBool("LOG_COMPRESSED", true)
	conf.Logger.IsDailyRotate = cc.AsBool("LOG_DAILY_ROTATE", true)
	conf.Logger.SleepDuration = cc.AsTimeDuration("LOG_SLEEP_DURATION", 5*time.Second)

	conf.Swagger.Username = cc.AsString("SWAGGER_USERNAME", "admin")
	conf.Swagger.Password = cc.AsString("SWAGGER_PASSWORD", "password")

	conf.DB.Name = cc.AsString("DB_NAME", "app")
	conf.DB.Host = cc.AsString("DB_HOST", "localhost")
	conf.DB.Port = cc.AsInt("DB_PORT", 3306)
	conf.DB.User = cc.AsString("DB_USER", "root")
	conf.DB.Password = cc.AsString("DB_PASSWORD", "password")
	conf.DB.ConnectionIdle = cc.AsTimeDuration("DB_CONNECTION_IDLE", 1*time.Minute)
	conf.DB.ConnectionLifetime = cc.AsTimeDuration("DB_CONNECTION_LIFETIME", 5*time.Minute)
	conf.DB.MaxIdle = cc.AsInt("DB_MAX_IDLE", 20)
	conf.DB.MaxOpen = cc.AsInt("DB_MAX_OPEN", 50)

	conf.Redis.Host = cc.AsString("REDIS_HOST", "")
	conf.Redis.Port = cc.AsInt("REDIS_PORT", 6379)
	conf.Redis.User = cc.AsString("REDIS_USER", "")
	conf.Redis.Password = cc.AsString("REDIS_PASSWORD", "")
	conf.Redis.ExpiredDuration = cc.AsTimeDuration("REDIS_EXPIRED_DURATION", 15*time.Minute)
	conf.Redis.DialTimeout = cc.AsTimeDuration("REDIS_DIAL_TIMEOUT", 5*time.Minute)

	conf.JWT.Secret = cc.AsString("JWT_SECRET", "secret")
	conf.JWT.Expire = cc.AsTimeDuration("JWT_EXPIRE", time.Hour*24*7) // 7 days

	conf.Auth.Google.Scopes = []string{"profile", "email"}
	conf.Auth.Google.RedirectURL = cc.AsString("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback")
	conf.Auth.Google.ClientID = cc.AsString("GOOGLE_CLIENT_ID", "")
	conf.Auth.Google.ClientSecret = cc.AsString("GOOGLE_CLIENT_SECRET", "")
	conf.Auth.Google.Endpoint = google.Endpoint

	conf.Auth.Verification.Callback = cc.AsString("VERIFICATION_CALLBACK", "")
	conf.Auth.Verification.Duration = cc.AsTimeDuration("VERIFICATION_DURATION", 15*time.Minute)
	conf.Auth.OTP.Duration = cc.AsTimeDuration("OTP_DURATION", 15*time.Minute)
	conf.Auth.OTP.MaxAttemp = cc.AsInt("OTP_MAX_ATTEMP", 3)

	conf.Whatsapp.ApiUrl = cc.AsString("WHATSAPP_API_URL", "http://localhost:3001")
}

func Get() *models.Config {
	return conf
}
