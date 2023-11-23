package handler

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/service"
)

type authHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) AuthHandler {
	return &authHandler{
		service,
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
	_ = ctx.QueryParser(&req)
	fmt.Println(req)

	mssage := h.service.GoogleCallback(ctx.Context(), req)
	return ctx.JSON(fiber.Map{
		"message": mssage,
	})
	// return ctx.Redirect(url, http.StatusTemporaryRedirect)
}
