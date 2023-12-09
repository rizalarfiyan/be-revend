package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/rizalarfiyan/be-revend/internal/response"
	"github.com/rizalarfiyan/be-revend/internal/sql"
	"github.com/rizalarfiyan/be-revend/utils"
)

var (
	errAuthNotImplement = errors.New("auth not implement in this api")
	errForbidden        = errors.New("you don't have access to this api")
)

func baseRoles(roles []sql.Role, isList bool) fiber.Handler {
	data := utils.DefaultErrorData(isList)
	mapRoles := make(map[sql.Role]bool)

	for _, role := range roles {
		if _, ok := mapRoles[role]; !ok {
			mapRoles[role] = true
		}
	}

	return func(ctx *fiber.Ctx) error {
		current := utils.GetUser(ctx)
		if utils.IsEmpty(current) {
			return ctx.Status(fiber.StatusInternalServerError).JSON(response.BaseResponse{
				Code:    fiber.StatusInternalServerError,
				Message: errAuthNotImplement.Error(),
				Data:    data,
			})
		}

		if _, ok := mapRoles[current.Role]; !ok {
			return ctx.Status(fiber.StatusForbidden).JSON(response.BaseResponse{
				Code:    fiber.StatusForbidden,
				Message: errForbidden.Error(),
				Data:    data,
			})
		}

		return ctx.Next()
	}
}

func Role(role sql.Role, isList bool) fiber.Handler {
	return baseRoles([]sql.Role{role}, isList)
}

func Roles(roles []sql.Role, isList bool) fiber.Handler {
	return baseRoles(roles, isList)
}
