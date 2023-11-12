package config

import (
	"strings"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rizalarfiyan/be-revend/constants"
	"github.com/rizalarfiyan/be-revend/exception"
	"github.com/rs/zerolog"

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

func FiberZerolog(logApi zerolog.Logger) fiberzerolog.Config {
	fields := []string{
		fiberzerolog.FieldIP,
		fiberzerolog.FieldMethod,
		fiberzerolog.FieldPath,
		fiberzerolog.FieldURL,
		fiberzerolog.FieldMethod,
		fiberzerolog.FieldPath,
		fiberzerolog.FieldLatency,
		fiberzerolog.FieldStatus,
		fiberzerolog.FieldBody,
		fiberzerolog.FieldError,
		fiberzerolog.FieldRequestID,
	}

	return fiberzerolog.Config{
		Logger: &logApi,
		Fields: fields,
		SkipBody: func(ctx *fiber.Ctx) bool {
			return strings.Contains(string(ctx.Request().Header.ContentType()), "multipart/form-data")
		},
	}
}
