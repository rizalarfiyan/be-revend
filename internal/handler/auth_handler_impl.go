package handler

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rizalarfiyan/be-revend/config"
	"github.com/rizalarfiyan/be-revend/exception"
	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/response"
	"github.com/rizalarfiyan/be-revend/internal/service"
	baseModels "github.com/rizalarfiyan/be-revend/models"
	"github.com/rizalarfiyan/be-revend/utils"
	"github.com/segmentio/ksuid"
)

type authHandler struct {
	service service.AuthService
	conf    *baseModels.Config
}

func NewAuthHandler(service service.AuthService) AuthHandler {
	return &authHandler{
		service: service,
		conf:    config.Get(),
	}
}

// Auth Me godoc
//
//	@Summary		Get Auth Me based on parameter
//	@Description	Auth Me
//	@ID				get-auth-me
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Security		AccessToken
//	@Success		200	{object}	response.BaseResponse{data=models.AuthToken}
//	@Failure		500	{object}	response.BaseResponse
//	@Router			/auth/me [get]
func (h *authHandler) Me(ctx *fiber.Ctx) error {
	user := utils.GetUser(ctx)
	return ctx.JSON(response.BaseResponse{
		Code:    http.StatusOK,
		Message: "Success!",
		Data:    user,
	})
}

// Auth Google Redirection godoc
//
//	@Summary		Get Auth Google Redirection based on parameter
//	@Description	Auth Google Redirection
//	@ID				get-auth-google
//	@Tags			auth
//	@Success		307
//	@Failure		500
//	@Router			/auth/google [get]
func (h *authHandler) Google(ctx *fiber.Ctx) error {
	url := h.service.Google()
	return ctx.Redirect(url, http.StatusTemporaryRedirect)
}

// Auth Google Redirection Callback godoc
//
//	@Summary		Get Auth Google Callback Redirection based on parameter
//	@Description	Auth Google Callback Redirection
//	@ID				get-auth-google-callback
//	@Tags			auth
//	@Success		307
//	@Failure		500
//	@Router			/auth/google/callback [get]
func (h *authHandler) GoogleCallback(ctx *fiber.Ctx) error {
	req := new(request.GoogleCallbackRequest)
	err := ctx.QueryParser(req)
	if err != nil {
		return ctx.Redirect(h.conf.Auth.Verification.Callback, http.StatusTemporaryRedirect)
	}

	url := h.service.GoogleCallback(ctx.Context(), *req)
	return ctx.Redirect(url, http.StatusTemporaryRedirect)
}

// Auth Register godoc
//
//	@Summary		Post Auth Register based on parameter
//	@Description	Auth Register
//	@ID				post-auth-register
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.AuthRegister	true	"Data"
//	@Success		200		{object}	response.BaseResponse
//	@Failure		500		{object}	response.BaseResponse
//	@Router			/auth/register [post]
func (h *authHandler) Register(ctx *fiber.Ctx) error {
	req := new(request.AuthRegister)
	err := ctx.BodyParser(req)
	if err != nil {
		return err
	}

	exception.ValidateStruct(*req, false)

	_, err = ksuid.Parse(req.Token)
	exception.IsNotProcessErrorMessage(err, "Token is not valid", false)

	h.service.Register(ctx.Context(), *req)
	return ctx.JSON(response.BaseResponse{
		Code:    http.StatusOK,
		Message: "Success!",
	})
}

// Auth Verification godoc
//
//	@Summary		Post Auth Verification based on parameter
//	@Description	Auth Verification
//	@ID				post-auth-verification
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.AuthVerification	true	"Data"
//	@Success		200		{object}	response.BaseResponse{data=response.AuthVerification}
//	@Failure		500		{object}	response.BaseResponse
//	@Router			/auth/verification [post]
func (h *authHandler) Verification(ctx *fiber.Ctx) error {
	req := new(request.AuthVerification)
	err := ctx.BodyParser(req)
	if err != nil {
		return err
	}

	exception.ValidateStruct(*req, false)

	token, err := ksuid.Parse(req.Token)
	exception.IsNotProcessErrorMessage(err, "Token is not valid", false)

	if !token.Time().Add(h.conf.Auth.Verification.Duration).After(time.Now()) {
		exception.IsNotProcessRawMessage("Token is not valid", false)
	}

	res := h.service.Verification(ctx.Context(), *req)
	return ctx.JSON(response.BaseResponse{
		Code:    http.StatusOK,
		Message: res.Message,
		Data:    res,
	})
}

// Auth Send OTP godoc
//
//	@Summary		Post Auth Send OTP based on parameter
//	@Description	Auth Send OTP
//	@ID				post-auth-send-otp
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.AuthSendOTP	true	"Data"
//	@Success		200		{object}	response.BaseResponse
//	@Failure		500		{object}	response.BaseResponse
//	@Router			/auth/otp [post]
func (h *authHandler) SendOTP(ctx *fiber.Ctx) error {
	req := new(request.AuthSendOTP)
	err := ctx.BodyParser(req)
	if err != nil {
		return err
	}

	exception.ValidateStruct(*req, false)

	res := h.service.SendOTP(ctx.Context(), *req)
	return ctx.JSON(response.BaseResponse{
		Code:    http.StatusOK,
		Message: "Success!",
		Data:    res,
	})
}

// Auth OTP Verification godoc
//
//	@Summary		Post Auth Verification based on parameter
//	@Description	Auth Verification
//	@ID				post-auth-otp-verification
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.AuthOTPVerification	true	"Data"
//	@Success		200		{object}	response.BaseResponse
//	@Failure		500		{object}	response.BaseResponse
//	@Router			/auth/otp/verification [post]
func (h *authHandler) OTPVerification(ctx *fiber.Ctx) error {
	req := new(request.AuthOTPVerification)
	err := ctx.BodyParser(req)
	if err != nil {
		return err
	}

	exception.ValidateStruct(*req, false)

	res := h.service.OTPVerification(ctx.Context(), *req)
	return ctx.JSON(response.BaseResponse{
		Code:    http.StatusOK,
		Message: "Success!",
		Data:    res,
	})
}
