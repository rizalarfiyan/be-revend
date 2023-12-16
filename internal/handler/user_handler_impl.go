package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-revend/config"
	"github.com/rizalarfiyan/be-revend/constants"
	"github.com/rizalarfiyan/be-revend/exception"
	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/response"
	"github.com/rizalarfiyan/be-revend/internal/service"
	"github.com/rizalarfiyan/be-revend/internal/sql"
	baseModels "github.com/rizalarfiyan/be-revend/models"
	"github.com/rizalarfiyan/be-revend/utils"
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

// GetAllUser godoc
//
//	@Summary		Get All User based on parameter
//	@Description	All User
//	@ID				get-all-user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		AccessToken
//	@Param			page		query		int		false	"Page"	default(1)
//	@Param			limit		query		int		false	"Limit"	default(10)
//	@Param			search		query		string	false	"Search"
//	@Param			order_by	query		string	false	"Order by"	Enums(first_name,last_name,phone_number)
//	@Param			order		query		string	false	"Order"		Enums(asc, desc)
//	@Param			status		query		string	false	"Status"	enum(,active,deleted)
//	@Success		200			{object}	response.BaseResponse{data=response.BaseResponsePagination[response.User]}
//	@Failure		500			{object}	response.BaseResponse
//	@Router			/user [get]
func (h *userHandler) GetAllUser(ctx *fiber.Ctx) error {
	req := request.GetAllUserRequest{
		BasePagination: request.BasePagination{
			Page:    ctx.QueryInt("page", 1),
			Limit:   ctx.QueryInt("limit", constants.DefaultPageLimit),
			Search:  ctx.Query("search"),
			OrderBy: ctx.Query("order_by"),
			Order:   ctx.Query("order"),
			Status:  constants.FilterListStatus(ctx.Query("status")),
		},
	}

	rawRole := ctx.Query("role")
	if rawRole != "" && utils.IsValidRole(rawRole) {
		req.Role = sql.Role(rawRole)
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

	res := h.service.GetAllUser(ctx.Context(), req)
	return ctx.JSON(response.BaseResponse{
		Code:    http.StatusOK,
		Message: "Success!",
		Data:    res,
	})
}

// GetUserById godoc
//
//	@Summary		Get User By Id based on parameter
//	@Description	Get User By Id
//	@ID				get-user-by-id
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		AccessToken
//	@Param			id	path		string	true	"User Id"
//	@Success		200	{object}	response.BaseResponse{data=response.User}
//	@Failure		500	{object}	response.BaseResponse
//	@Router			/user/{id} [get]
func (h *userHandler) GetUserById(ctx *fiber.Ctx) error {
	userId, err := uuid.Parse(ctx.Params("id", ""))
	if err != nil {
		exception.IsNotFound(nil, false)
	}

	res := h.service.GetUserById(ctx.Context(), userId)
	return ctx.JSON(response.BaseResponse{
		Code:    http.StatusOK,
		Message: "Success!",
		Data:    res,
	})
}

// AllDropdownUser godoc
//
//	@Summary		Get All Dropdown User based on parameter
//	@Description	All Dropdown User
//	@ID				get-all-dropdown-user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		AccessToken
//	@Param			page	query		int		false	"Page"	default(1)
//	@Param			limit	query		int		false	"Limit"	default(10)
//	@Param			search	query		string	false	"Search"
//	@Param			status	query		string	false	"Status"	enum(,active,deleted)
//	@Success		200		{object}	response.BaseResponse{data=response.BaseResponsePagination[response.BaseDropdown]}
//	@Failure		500		{object}	response.BaseResponse
//	@Router			/user/dropdown [get]
func (h *userHandler) AllDropdownUser(ctx *fiber.Ctx) error {
	req := request.BasePagination{
		Page:   ctx.QueryInt("page", 1),
		Limit:  ctx.QueryInt("limit", constants.DefaultPageLimit),
		Search: ctx.Query("search"),
		Status: constants.FilterListStatus(ctx.Query("status")),
	}

	req.Normalize()

	user := utils.GetUser(ctx)
	if user.Role != sql.RoleAdmin {
		req.Status = constants.FilterListStatusActive
	}

	res := h.service.GetAllDropdownUser(ctx.Context(), req)
	return ctx.JSON(response.BaseResponse{
		Code:    http.StatusOK,
		Message: "Success!",
		Data:    res,
	})
}

// Create User godoc
//
//	@Summary		Post Create User based on parameter
//	@Description	Create User
//	@ID				post-create-user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		AccessToken
//	@Param			data	body		request.CreateUserRequest	true	"Data"
//	@Success		200		{object}	response.BaseResponse
//	@Failure		500		{object}	response.BaseResponse
//	@Router			/user [post]
func (h *userHandler) CreateUser(ctx *fiber.Ctx) error {
	req := new(request.CreateUserRequest)
	err := ctx.BodyParser(req)
	if err != nil {
		return err
	}

	exception.ValidateStruct(*req, false)
	if !utils.IsValidRole(req.RawRole) {
		exception.ErrorManualValidation("role", "Role is not valid")
	}
	req.Role = sql.Role(req.RawRole)

	h.service.CreateUser(ctx.Context(), *req)
	return ctx.JSON(response.BaseResponse{
		Code:    http.StatusOK,
		Message: "Success!",
	})
}

// Update User godoc
//
//	@Summary		Post Update User based on parameter
//	@Description	Update User
//	@ID				post-update-user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		AccessToken
//	@Param			data	body		request.UpdateUserRequest	true	"Data"
//	@Param			id		path		string						true	"User ID"	example(550e8400-e29b-41d4-a716-446655440000)	Format(uuid)
//	@Success		200		{object}	response.BaseResponse
//	@Failure		500		{object}	response.BaseResponse
//	@Router			/user/{id} [put]
func (h *userHandler) UpdateUser(ctx *fiber.Ctx) error {
	req := new(request.UpdateUserRequest)
	err := ctx.BodyParser(req)
	if err != nil {
		return err
	}

	userId, err := uuid.Parse(ctx.Params("id"))
	exception.IsNotProcessErrorMessage(err, "Path id is not a valid uuid format", false)
	req.Id = userId

	exception.ValidateStruct(*req, false)
	if !utils.IsValidRole(req.RawRole) {
		exception.ErrorManualValidation("role", "Role is not valid")
	}
	req.Role = sql.Role(req.RawRole)

	h.service.UpdateUser(ctx.Context(), *req)
	return ctx.JSON(response.BaseResponse{
		Code:    http.StatusOK,
		Message: "Success!",
	})
}

// ToggleDeleteUser godoc
//
//	@Summary		Toggle Delete User based on parameter
//	@Description	Toggle Delete User
//	@ID				toggle-delete-user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		AccessToken
//	@Param			id	path		string	true	"User ID"	example(550e8400-e29b-41d4-a716-446655440000)	Format(uuid)
//	@Success		200	{object}	response.BaseResponse
//	@Failure		500	{object}	response.BaseResponse
//	@Router			/user/{id} [delete]
func (h *userHandler) ToggleDeleteUser(ctx *fiber.Ctx) error {
	rawId := ctx.Params("id")
	userId, err := uuid.Parse(rawId)
	exception.IsNotProcessErrorMessage(err, "Path id is not a valid uuid format", false)

	user := utils.GetUser(ctx)
	h.service.ToggleDeleteUser(ctx.Context(), userId, user.Id)
	return ctx.JSON(response.BaseResponse{
		Code:    http.StatusOK,
		Message: "Success!",
	})
}
