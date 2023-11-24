package internal

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rizalarfiyan/be-revend/internal/handler"
	"github.com/rizalarfiyan/be-revend/middleware"
)

type router struct {
	app *fiber.App
}

func NewRouter(app *fiber.App) Router {
	return &router{
		app: app,
	}
}

func (r *router) BaseRoute(handler handler.BaseHandler) {
	r.app.Get("", handler.Home)
	r.app.Get("health", handler.Health)
}

func (r *router) AuthRoute(handler handler.AuthHandler) {
	auth := r.app.Group("auth")
	auth.Post("verification", handler.Verification)
	auth.Get("me", middleware.Auth(false), handler.Me)

	otp := auth.Group("otp")
	otp.Post("", handler.SendOTP)

	google := auth.Group("google")
	google.Get("", handler.Google)
	google.Get("callback", handler.GoogleCallback)
}
