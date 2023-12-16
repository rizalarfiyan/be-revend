package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rizalarfiyan/be-revend/database"
	"github.com/rizalarfiyan/be-revend/internal/response"
)

type baseHandler struct{}

func NewBaseHandler() BaseHandler {
	return &baseHandler{}
}

// Base Home godoc
//
//	@Summary		Get Base Home based on parameter
//	@Description	Base Home
//	@ID				get-base-home
//	@Tags			home
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.BaseResponse
//	@Failure		500	{object}	response.BaseResponse
//	@Router			/ [get]
func (h *baseHandler) Home(ctx *fiber.Ctx) error {
	return ctx.JSON(response.BaseResponse{
		Code:    http.StatusOK,
		Message: "Success!",
		Data: fiber.Map{
			"title": "BE Revend",
			"author": fiber.Map{
				"name":    "Muhamad Rizal Arfiyan",
				"email":   "rizal.arfiyan.23@gmail.com",
				"website": "https://rizalrfiyan.com",
				"github":  "https://github.com/rizalarfiyan",
			},
		},
	})
}

// Base Health godoc
//
//	@Summary		Get Base Health based on parameter
//	@Description	Base Health
//	@ID				get-base-health
//	@Tags			home
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.BaseResponse
//	@Failure		500	{object}	response.BaseResponse
//	@Router			/health [get]
func (h *baseHandler) Health(ctx *fiber.Ctx) error {
	return ctx.JSON(response.BaseResponse{
		Code:    http.StatusOK,
		Message: "Success!",
		Data: fiber.Map{
			"postgres": database.PostgresIsConnected(),
			"redis":    database.RedisIsConnected(),
		},
	})
}
