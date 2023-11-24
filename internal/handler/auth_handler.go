package handler

import "github.com/gofiber/fiber/v2"

type AuthHandler interface {
	Me(ctx *fiber.Ctx) error
	Google(ctx *fiber.Ctx) error
	GoogleCallback(ctx *fiber.Ctx) error
	Verification(ctx *fiber.Ctx) error
	SendOTP(ctx *fiber.Ctx) error
}
