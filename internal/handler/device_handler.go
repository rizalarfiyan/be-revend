package handler

import "github.com/gofiber/fiber/v2"

type DeviceHandler interface {
	GetAllDevice(ctx *fiber.Ctx) error
	AllDropdownDevice(ctx *fiber.Ctx) error
	CreateDevice(ctx *fiber.Ctx) error
	UpdateDevice(ctx *fiber.Ctx) error
	ToggleDeleteDevice(ctx *fiber.Ctx) error
}
