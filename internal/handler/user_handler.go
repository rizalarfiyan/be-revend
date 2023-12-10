package handler

import "github.com/gofiber/fiber/v2"

type UserHandler interface {
	AllUser(ctx *fiber.Ctx) error
}
