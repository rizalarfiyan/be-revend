package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rizalarfiyan/be-revend/database"
	"github.com/rizalarfiyan/be-revend/internal/response"
	"github.com/rizalarfiyan/be-revend/internal/sql"
	"github.com/rizalarfiyan/be-revend/utils"
)

type baseHandler struct{}

func NewBaseHandler() BaseHandler {
	return &baseHandler{}
}

// Base Home godoc
// @Summary      Get Base Home based on parameter
// @Description  Base Home
// @ID           get-base-home
// @Tags         home
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.BaseResponse
// @Failure      500  {object}  response.BaseResponse
// @Router       / [get]
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
// @Summary      Get Base Health based on parameter
// @Description  Base Health
// @ID           get-base-health
// @Tags         home
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.BaseResponse
// @Failure      500  {object}  response.BaseResponse
// @Router       /health [get]
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

func (h *baseHandler) Test(ctx *fiber.Ctx) error {
	db := database.GetPostgres()
	query := sql.New(utils.QueryWrap(db))

	users, err := query.GetAllUsers(utils.QueryBuild(ctx.Context(), func(b *utils.QueryBuilder) {
		b.Where("phone_number = $1", "62895377233002")
		b.Limit(10)
	}))

	if err != nil {
		log.Fatalln("ListAuthors", err)
	}

	data, _ := json.MarshalIndent(users, "", "  ")
	fmt.Println(string(data))

	return ctx.SendString("OKE!")
}
