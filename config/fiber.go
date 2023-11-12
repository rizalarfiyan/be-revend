package config

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rizalarfiyan/be-revend/constants"
	"github.com/rizalarfiyan/be-revend/exception"

	"github.com/gofiber/fiber/v2"
)

func FiberConfig() fiber.Config {
	return fiber.Config{
		ErrorHandler: exception.ErrorHandler,
		BodyLimit:    constants.FiberBodyLimit,
	}
}

func CorsConfig() cors.Config {
	return cors.Config{
		AllowOrigins:     conf.Cors.AllowOrigins,
		AllowMethods:     conf.Cors.AllowMethods,
		AllowHeaders:     conf.Cors.AllowHeaders,
		AllowCredentials: conf.Cors.AllowCredentials,
		ExposeHeaders:    conf.Cors.ExposeHeaders,
	}
}
