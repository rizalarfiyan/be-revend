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

type deviceHandler struct {
	service service.DeviceService
}

func NewDeviceHandler(service service.DeviceService) DeviceHandler {
	return &deviceHandler{
		service: service,
	}
}

// GetAllDevice godoc
//
//	@Summary		Get All Device based on parameter
//	@Description	All Device
//	@ID				get-all-device
//	@Tags			device
//	@Accept			json
//	@Produce		json
//	@Security		AccessToken
//	@Param			page		query		int		false	"Page"	default(1)
//	@Param			limit		query		int		false	"Limit"	default(10)
//	@Param			search		query		string	false	"Search"
//	@Param			order_by	query		string	false	"Order by"	Enums(token,name,location)
//	@Param			order		query		string	false	"Order"		Enums(asc, desc)
//	@Success		200			{object}	response.BaseResponse{data=response.BaseResponsePagination[response.Device]}
//	@Failure		500			{object}	response.BaseResponse
//	@Router			/device [get]
func (h *deviceHandler) GetAllDevice(ctx *fiber.Ctx) error {
	req := request.BasePagination{
		Page:    ctx.QueryInt("page", 1),
		Limit:   ctx.QueryInt("limit", constants.DefaultPageLimit),
		Search:  ctx.Query("search"),
		OrderBy: ctx.Query("order_by"),
		Order:   ctx.Query("order"),
		Status:  constants.FilterListStatus(ctx.Query("status")),
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

// AllDropdownDevice godoc
//
//	@Summary		Get All Dropdown Device based on parameter
//	@Description	All Dropdown Device
//	@ID				get-all-dropdown-device
//	@Tags			device
//	@Accept			json
//	@Produce		json
//	@Security		AccessToken
//	@Param			page	query		int		false	"Page"	default(1)
//	@Param			limit	query		int		false	"Limit"	default(10)
//	@Param			search	query		string	false	"Search"
//	@Success		200		{object}	response.BaseResponse{data=response.BaseResponsePagination[response.BaseDropdown]}
//	@Failure		500		{object}	response.BaseResponse
//	@Router			/device/dropdown [get]
func (h *deviceHandler) AllDropdownDevice(ctx *fiber.Ctx) error {
	req := request.AllDropdownDeviceRequest{
		BasePagination: request.BasePagination{
			Page:   ctx.QueryInt("page", 1),
			Limit:  ctx.QueryInt("limit", constants.DefaultPageLimit),
			Search: ctx.Query("search"),
		},
	}

	user := utils.GetUser(ctx)
	if user.Role != sql.RoleAdmin {
		req.HideDeleted = true
	}

	req.Normalize()

	res := h.service.GetAllDropdownDevice(ctx.Context(), req)
	return ctx.JSON(response.BaseResponse{
		Code:    http.StatusOK,
		Message: "Success!",
		Data:    res,
	})
}

// Create Device godoc
//
//	@Summary		Post Create Device based on parameter
//	@Description	Create Device
//	@ID				post-create-device
//	@Tags			device
//	@Accept			json
//	@Produce		json
//	@Security		AccessToken
//	@Param			data	body		request.CreateDeviceRequest	true	"Data"
//	@Success		200		{object}	response.BaseResponse
//	@Failure		500		{object}	response.BaseResponse
//	@Router			/device [post]
func (h *deviceHandler) CreateDevice(ctx *fiber.Ctx) error {
	req := new(request.CreateDeviceRequest)
	err := ctx.BodyParser(req)
	if err != nil {
		return err
	}

	exception.ValidateStruct(*req, false)

	h.service.CreateDevice(ctx.Context(), *req)
	return ctx.JSON(response.BaseResponse{
		Code:    http.StatusOK,
		Message: "Success!",
	})
}

// Update Device godoc
//
//	@Summary		Post Update Device based on parameter
//	@Description	Update Device
//	@ID				post-update-device
//	@Tags			device
//	@Accept			json
//	@Produce		json
//	@Security		AccessToken
//	@Param			data	body		request.UpdateDeviceRequest	true	"Data"
//	@Param			id		path		string						true	"Device ID"	example(550e8400-e29b-41d4-a716-446655440000)	Format(uuid)
//	@Success		200		{object}	response.BaseResponse
//	@Failure		500		{object}	response.BaseResponse
//	@Router			/device/{id} [put]
func (h *deviceHandler) UpdateDevice(ctx *fiber.Ctx) error {
	req := new(request.UpdateDeviceRequest)
	err := ctx.BodyParser(req)
	if err != nil {
		return err
	}

	deviceId, err := uuid.Parse(ctx.Params("id"))
	exception.IsNotProcessErrorMessage(err, "Path id is not a valid uuid format", false)
	req.Id = deviceId

	exception.ValidateStruct(*req, false)

	h.service.UpdateDevice(ctx.Context(), *req)
	return ctx.JSON(response.BaseResponse{
		Code:    http.StatusOK,
		Message: "Success!",
	})
}

// ToggleDeleteDevice godoc
//
//	@Summary		Toggle Delete Device based on parameter
//	@Description	Toggle Delete Device
//	@ID				toggle-delete-device
//	@Tags			device
//	@Accept			json
//	@Produce		json
//	@Security		AccessToken
//	@Param			id	path		string	true	"Device ID"	example(550e8400-e29b-41d4-a716-446655440000)	Format(uuid)
//	@Success		200	{object}	response.BaseResponse
//	@Failure		500	{object}	response.BaseResponse
//	@Router			/device/{id} [delete]
func (h *deviceHandler) ToggleDeleteDevice(ctx *fiber.Ctx) error {
	rawId := ctx.Params("id")
	deviceId, err := uuid.Parse(rawId)
	exception.IsNotProcessErrorMessage(err, "Path id is not a valid uuid format", false)

	user := utils.GetUser(ctx)
	h.service.ToggleDeleteDevice(ctx.Context(), deviceId, user.Id)
	return ctx.JSON(response.BaseResponse{
		Code:    http.StatusOK,
		Message: "Success!",
	})
}
