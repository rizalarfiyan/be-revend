package handler

import "github.com/gofiber/fiber/v2"

type HistoryHandler interface {
	GetAllHistory(ctx *fiber.Ctx) error
	GetAllHistoryStatistic(ctx *fiber.Ctx) error
}
