package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-revend/constants"
	"github.com/rizalarfiyan/be-revend/exception"
	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/response"
	"github.com/rizalarfiyan/be-revend/internal/service"
	"github.com/rizalarfiyan/be-revend/internal/sql"
	"github.com/rizalarfiyan/be-revend/utils"
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
//
//	@Summary		Get All History based on parameter
//	@Description	All History
//	@ID				get-all-history
//	@Tags			history
//	@Accept			json
//	@Produce		json
//	@Security		AccessToken
//	@Param			page		query		int		false	"Page"	default(1)
//	@Param			limit		query		int		false	"Limit"	default(10)
//	@Param			search		query		string	false	"Search"
//	@Param			order_by	query		string	false	"Order by"	Enums(success,failed,name,device)
//	@Param			order		query		string	false	"Order"		Enums(asc, desc)
//	@Param			device_id	query		string	false	"Device ID"	example(550e8400-e29b-41d4-a716-446655440000)	Format(uuid)
//	@Param			user_id		query		string	false	"User ID"	example(550e8400-e29b-41d4-a716-446655440000)	Format(uuid)
//	@Param			status		query		string	false	"Status"	Enums(active,deleted)
//	@Success		200			{object}	response.BaseResponse{data=response.BaseResponsePagination[response.History]}
//	@Failure		500			{object}	response.BaseResponse
//	@Router			/history [get]
func (h *historyHandler) GetAllHistory(ctx *fiber.Ctx) error {
	req := request.GetAllHistoryRequest{
		BasePagination: request.BasePagination{
			Page:    ctx.QueryInt("page", 1),
			Limit:   ctx.QueryInt("limit", constants.DefaultPageLimit),
			Search:  ctx.Query("search"),
			OrderBy: ctx.Query("order_by"),
			Order:   ctx.Query("order"),
			Status:  constants.FilterListStatus(ctx.Query("status")),
		},
	}

	var err error
	rawDeviceId := ctx.Query("device_id")
	if rawDeviceId != "" {
		req.DeviceId, err = uuid.Parse(rawDeviceId)
		exception.ErrorListPaginationValidation[response.History](err, "invalid UUID format for device ID", req.BasePagination)
	}

	rawUserId := ctx.Query("user_id")
	if rawUserId != "" {
		req.UserId, err = uuid.Parse(rawUserId)
		exception.ErrorListPaginationValidation[response.History](err, "invalid UUID format for user ID", req.BasePagination)
	}

	fieldOrder := map[string]string{
		"success": "h.success",
		"failed":  "h.failed",
		"name":    "CONCAT(u.first_name, ' ', u.last_name)",
		"device":  "d.name",
		"date":    "h.created_at",
	}

	user := utils.GetUser(ctx)
	if user.Role != sql.RoleAdmin {
		req.UserId = user.Id
		req.Status = constants.FilterListStatusActive
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

// GetAllHistoryStatistic godoc
//
//	@Summary		Get All History Statistic based on parameter
//	@Description	All History Statistic
//	@ID				get-all-history-statistic
//	@Tags			history
//	@Accept			json
//	@Produce		json
//	@Security		AccessToken
//	@Param			time_frequency	query		string	false	"Time Frequency"	Enums(today,week,month,quarter,year)
//	@Success		200				{object}	response.BaseResponse{data=[]response.HistoryStatistic}
//	@Failure		500				{object}	response.BaseResponse
//	@Router			/history/statistic [get]
func (h *historyHandler) GetAllHistoryStatistic(ctx *fiber.Ctx) error {
	req := request.GetAllHistoryStatisticRequest{
		WithTimeFrequency: request.WithTimeFrequency{
			TimeFrequency: constants.FilterTimeFrequency(ctx.Query("time_frequency")),
		},
		UserId: utils.GetUser(ctx).Id,
	}

	req.Normalize()

	res := h.service.GetAllHistoryStatistic(ctx.Context(), req)
	return ctx.JSON(response.BaseResponse{
		Code:    http.StatusOK,
		Message: "Success!",
		Data:    res,
	})
}
