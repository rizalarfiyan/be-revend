package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rizalarfiyan/be-revend/constants"
	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/response"
	"github.com/rizalarfiyan/be-revend/internal/service"
)

type deviceHandler struct {
	service service.DeviceService
}

func NewDeviceHandler(service service.DeviceService) DeviceHandler {
	return &deviceHandler{
		service: service,
	}
}

// GetAllDevice godoc
// @Summary      Get All Device based on parameter
// @Description  All Device
// @ID           get-all-device
// @Tags         device
// @Accept       json
// @Produce      json
// @Security     AccessToken
// @Param        page query int false "Page" default(1)
// @Param        limit query int false "Limit" default(10)
// @Param        search query string false "Search"
// @Param        order_by query string false "Order by" Enums(token,name,location)
// @Param        order query string false "Order" Enums(asc, desc)
// @Success      200  {object}  response.BaseResponse{data=response.BaseResponsePagination[response.Device]}
// @Failure      500  {object}  response.BaseResponse
// @Router       /device [get]
func (h *deviceHandler) GetAllDevice(ctx *fiber.Ctx) error {
	req := request.BasePagination{
		Page:    ctx.QueryInt("page", 1),
		Limit:   ctx.QueryInt("limit", constants.DefaultPageLimit),
		Search:  ctx.Query("search"),
		OrderBy: ctx.Query("order_by"),
		Order:   ctx.Query("order"),
	}

	fieldOrder := map[string]string{
		"token":    "token",
		"name":     "name",
		"location": "location",
	}

	req.ValidateAndUpdateOrderBy(fieldOrder)
	req.Normalize()

	res := h.service.GetAllDevice(ctx.Context(), req)
	return ctx.JSON(response.BaseResponse{
		Code:    http.StatusOK,
		Message: "Success!",
		Data:    res,
	})
}
