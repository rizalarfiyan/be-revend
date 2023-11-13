package handler

import "github.com/gofiber/fiber/v2"

type AuthHandler interface {
	Google(ctx *fiber.Ctx) error
}
