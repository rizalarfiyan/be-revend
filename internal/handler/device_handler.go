package handler

import "github.com/gofiber/fiber/v2"

type DeviceHandler interface {
	GetAllDevice(ctx *fiber.Ctx) error
}
