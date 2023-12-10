package handler

import "github.com/gofiber/fiber/v2"

type UserHandler interface {
	GetAllUser(ctx *fiber.Ctx) error
	GetUserById(ctx *fiber.Ctx) error
}
