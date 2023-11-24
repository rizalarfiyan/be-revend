package middleware

import (
	"errors"

	jwtMiddleware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/rizalarfiyan/be-revend/config"
	"github.com/rizalarfiyan/be-revend/constants"
	"github.com/rizalarfiyan/be-revend/internal/response"
	"github.com/rizalarfiyan/be-revend/utils"
)

var errJwt = errors.New("unauthorize, invalid or expired jwt")

func baseAuth(isList bool, isMandatory bool) fiber.Handler {
	conf := config.Get()
	data := utils.DefaultErrorData(isList)
	return jwtMiddleware.New(jwtMiddleware.Config{
		ContextKey: constants.KeyLocalsUser,
		SigningKey: jwtMiddleware.SigningKey{
			Key: []byte(conf.JWT.Secret),
		},
		SuccessHandler: func(ctx *fiber.Ctx) error {
			_, err := utils.ValidateUser(ctx)
			if err != nil {
				if !isMandatory {
					return ctx.Next()
				}

				return ctx.Status(fiber.StatusUnauthorized).JSON(response.BaseResponse{
					Code: fiber.StatusUnauthorized,
					// force to error message "unauthorize, invalid or expired jwt", because not showing detail of error
					Message: errJwt.Error(),
					Data:    data,
				})
			}

			return ctx.Next()
		},
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			if !isMandatory {
				return ctx.Next()
			}

			return ctx.Status(fiber.StatusUnauthorized).JSON(response.BaseResponse{
				Code:    fiber.StatusUnauthorized,
				Message: errJwt.Error(),
				Data:    data,
			})
		},
	})
}

func Auth(isList bool) fiber.Handler {
	return baseAuth(isList, true)
}

func OptinalAuth(isList bool) fiber.Handler {
	return baseAuth(isList, false)
}
