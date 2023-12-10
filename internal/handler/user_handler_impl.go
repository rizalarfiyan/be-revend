package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rizalarfiyan/be-revend/config"
	"github.com/rizalarfiyan/be-revend/constants"
	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/response"
	"github.com/rizalarfiyan/be-revend/internal/service"
	baseModels "github.com/rizalarfiyan/be-revend/models"
)

type userHandler struct {
	service service.UserService
	conf    *baseModels.Config
}

func NewUserHandler(service service.UserService) UserHandler {
	return &userHandler{
		service: service,
		conf:    config.Get(),
	}
}

// AllUser godoc
// @Summary      Get All User based on parameter
// @Description  All User
// @ID           get-all-user
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     AccessToken
// @Param        page query int false "Page" default(1)
// @Param        limit query int false "Limit" default(10)
// @Param        search query string false "Search"
// @Param        order_by query string false "Order by" Enums(first_name,last_name,phone_number)
// @Param        order query string false "Order" Enums(asc, desc)
// @Success      200  {object}  response.BaseResponse{data=response.BaseResponsePagination[response.User]}
// @Failure      500  {object}  response.BaseResponse
// @Router       /user [get]
func (h *userHandler) AllUser(ctx *fiber.Ctx) error {
	req := request.BasePagination{
		Page:    ctx.QueryInt("page", 1),
		Limit:   ctx.QueryInt("limit", constants.DefaultPageLimit),
		Search:  ctx.Query("search"),
		OrderBy: ctx.Query("order_by"),
		Order:   ctx.Query("order"),
	}

	fieldOrder := map[string]string{
		"first_name":   "first_name",
		"last_name":    "last_name",
		"phone_number": "phone_number",
		"identity":     "identity",
		"role":         "role",
	}

	req.ValidateAndUpdateOrderBy(fieldOrder)
	req.Normalize()

	res := h.service.AllUser(ctx.Context(), req)
	return ctx.JSON(response.BaseResponse{
		Code:    http.StatusOK,
		Message: "Success!",
		Data:    res,
	})
}
