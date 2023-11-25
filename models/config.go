package models

import (
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
)

type Config struct {
	Env      string
	Port     int
	Host     string
	Name     string
	Logger   LoggerConfigs
	Cors     CorsConfigs
	Swagger  SwaggerConfigs
	DB       DBConfigs
	Redis    RedisConfigs
	JWT      JWTConfigs
	Auth     AuthConfigs
	Whatsapp WhatsappConfigs
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

type DBConfigs struct {
	Name               string
	Host               string
	Port               int
	User               string
	Password           string
	ConnectionIdle     time.Duration
	ConnectionLifetime time.Duration
	MaxIdle            int
	MaxOpen            int
}

type RedisConfigs struct {
	Host            string
	Port            int
	User            string
	Password        string
	ExpiredDuration time.Duration
	DialTimeout     time.Duration
}

type JWTConfigs struct {
	Secret string
	Expire time.Duration
}

type AuthConfigs struct {
	Google       oauth2.Config
	Verification VerificationAuthConfigs
	OTP          OTPConfigs
}

type VerificationAuthConfigs struct {
	Callback string
	Duration time.Duration
}

type OTPConfigs struct {
	Duration  time.Duration
	MaxAttemp int
}

type WhatsappConfigs struct {
	ApiUrl string
}
