package internal

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rizalarfiyan/be-revend/internal/handler"
	"github.com/rizalarfiyan/be-revend/internal/sql"
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
	auth.Post("register", handler.Register)
	auth.Get("me", middleware.Auth(false), handler.Me)

	otp := auth.Group("otp")
	otp.Post("", handler.SendOTP)
	otp.Post("verification", handler.OTPVerification)

	google := auth.Group("google")
	google.Get("", handler.Google)
	google.Get("callback", handler.GoogleCallback)
}

func (r *router) UserRoute(handler handler.UserHandler) {
	user := r.app.Group("user")
	user.Get("", middleware.Auth(true), middleware.Role(sql.RoleAdmin, true), handler.GetAllUser)
	user.Post("", middleware.Auth(false), middleware.Role(sql.RoleAdmin, false), handler.CreateUser)
	user.Get("dropdown", middleware.Auth(true), middleware.Role(sql.RoleAdmin, true), handler.AllDropdownUser)
	user.Get(":id", middleware.Auth(false), middleware.Role(sql.RoleAdmin, false), handler.GetUserById)
	user.Delete(":id", middleware.Auth(false), middleware.Role(sql.RoleAdmin, false), handler.ToggleDeleteUser)
}

func (r *router) DeviceRoute(handler handler.DeviceHandler) {
	device := r.app.Group("device")
	device.Get("", middleware.Auth(true), middleware.Role(sql.RoleAdmin, true), handler.GetAllDevice)
	device.Get("dropdown", middleware.Auth(true), handler.AllDropdownDevice)
	device.Post("", middleware.Auth(false), middleware.Role(sql.RoleAdmin, false), handler.CreateDevice)
	device.Put(":id", middleware.Auth(false), middleware.Role(sql.RoleAdmin, false), handler.UpdateDevice)
	device.Delete(":id", middleware.Auth(false), middleware.Role(sql.RoleAdmin, false), handler.ToggleDeleteDevice)
}

func (r *router) HistoryRoute(handler handler.HistoryHandler) {
	history := r.app.Group("history")
	history.Get("", middleware.Auth(true), handler.GetAllHistory)
}
