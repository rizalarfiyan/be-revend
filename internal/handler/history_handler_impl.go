package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rizalarfiyan/be-revend/constants"
	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/response"
	"github.com/rizalarfiyan/be-revend/internal/service"
)

type historyHandler struct {
	service service.HistoryService
}

func NewHistoryHandler(service service.HistoryService) HistoryHandler {
	return &historyHandler{
		service: service,
	}
}

// GetAllHistory godoc
// @Summary      Get All History based on parameter
// @Description  All History
// @ID           get-all-history
// @Tags         history
// @Accept       json
// @Produce      json
// @Security     AccessToken
// @Param        page query int false "Page" default(1)
// @Param        limit query int false "Limit" default(10)
// @Param        search query string false "Search"
// @Param        order_by query string false "Order by" Enums(success,failed,name,device)
// @Param        order query string false "Order" Enums(asc, desc)
// @Success      200  {object}  response.BaseResponse{data=response.BaseResponsePagination[response.History]}
// @Failure      500  {object}  response.BaseResponse
// @Router       /history [get]
func (h *historyHandler) GetAllHistory(ctx *fiber.Ctx) error {
	req := request.BasePagination{
		Page:    ctx.QueryInt("page", 1),
		Limit:   ctx.QueryInt("limit", constants.DefaultPageLimit),
		Search:  ctx.Query("search"),
		OrderBy: ctx.Query("order_by"),
		Order:   ctx.Query("order"),
	}

	fieldOrder := map[string]string{
		"success": "h.success",
		"failed":  "h.failed",
		"name":    "CONCAT(u.first_name, ' ', u.last_name)",
		"device":  "d.name",
	}

	req.ValidateAndUpdateOrderBy(fieldOrder)
	req.Normalize()

	res := h.service.GetAllHistory(ctx.Context(), req)
	return ctx.JSON(response.BaseResponse{
		Code:    http.StatusOK,
		Message: "Success!",
		Data:    res,
	})
}
