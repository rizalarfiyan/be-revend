package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
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

// Auth Google godoc
// @Summary      Get Auth Google based on parameter
// @Description  Auth Google
// @ID           get-auth-google
// @Tags         auth
// @Success      307
// @Failure      500
// @Router       /auth/google [get]
func (h *authHandler) Google(ctx *fiber.Ctx) error {
	url := h.service.Google()
	return ctx.Redirect(url, http.StatusTemporaryRedirect)
}
