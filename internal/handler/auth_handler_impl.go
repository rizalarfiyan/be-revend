package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rizalarfiyan/be-revend/config"
	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/service"
	baseModels "github.com/rizalarfiyan/be-revend/models"
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

// Auth Google Redirection godoc
// @Summary      Get Auth Google Redirection based on parameter
// @Description  Auth Google Redirection
// @ID           get-auth-google
// @Tags         auth
// @Success      307
// @Failure      500
// @Router       /auth/google [get]
func (h *authHandler) Google(ctx *fiber.Ctx) error {
	url := h.service.Google()
	return ctx.Redirect(url, http.StatusTemporaryRedirect)
}

// Auth Google Redirection Callback godoc
// @Summary      Get Auth Google Callback Redirection based on parameter
// @Description  Auth Google Callback Redirection
// @ID           get-auth-google-callback
// @Tags         auth
// @Success      307
// @Failure      500
// @Router       /auth/google/callback [get]
func (h *authHandler) GoogleCallback(ctx *fiber.Ctx) error {
	req := request.GoogleCallbackRequest{}
	err := ctx.QueryParser(&req)
	if err != nil {
		return ctx.Redirect(h.conf.Auth.Callback, http.StatusTemporaryRedirect)
	}

	url := h.service.GoogleCallback(ctx.Context(), req)
	return ctx.Redirect(url, http.StatusTemporaryRedirect)
}
