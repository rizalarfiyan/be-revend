package handler

import "github.com/gofiber/fiber/v2"

type BaseHandler interface {
	Home(ctx *fiber.Ctx) error
	Health(ctx *fiber.Ctx) error
	Test(ctx *fiber.Ctx) error
}
