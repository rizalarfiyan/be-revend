package internal

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rizalarfiyan/be-revend/internal/handler"
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
