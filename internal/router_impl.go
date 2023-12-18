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
	profile := user.Group("profile")
	profile.Get("", middleware.Auth(false), handler.GetUserProfile)
	profile.Put("", middleware.Auth(false), handler.UpdateUserProfile)
	profile.Delete("google", middleware.Auth(false), handler.DeleteGoogleUserProfile)

	user.Get("", middleware.Auth(true), middleware.Role(sql.RoleAdmin, true), handler.GetAllUser)
	user.Post("", middleware.Auth(false), middleware.Role(sql.RoleAdmin, false), handler.CreateUser)
	user.Get("dropdown", middleware.Auth(true), middleware.Role(sql.RoleAdmin, true), handler.AllDropdownUser)
	user.Get(":id", middleware.Auth(false), middleware.Role(sql.RoleAdmin, false), handler.GetUserById)
	user.Put(":id", middleware.Auth(false), middleware.Role(sql.RoleAdmin, false), handler.UpdateUser)
	user.Delete(":id", middleware.Auth(false), middleware.Role(sql.RoleAdmin, false), handler.ToggleDeleteUser)
}

func (r *router) DeviceRoute(handler handler.DeviceHandler) {
	device := r.app.Group("device")
	device.Get("", middleware.Auth(true), middleware.Role(sql.RoleAdmin, true), handler.GetAllDevice)
	device.Post("", middleware.Auth(false), middleware.Role(sql.RoleAdmin, false), handler.CreateDevice)
	device.Get("dropdown", middleware.Auth(true), handler.AllDropdownDevice)
	device.Put(":id", middleware.Auth(false), middleware.Role(sql.RoleAdmin, false), handler.UpdateDevice)
	device.Delete(":id", middleware.Auth(false), middleware.Role(sql.RoleAdmin, false), handler.ToggleDeleteDevice)
}

func (r *router) HistoryRoute(handler handler.HistoryHandler) {
	history := r.app.Group("history")
	history.Get("", middleware.Auth(true), handler.GetAllHistory)
	history.Get("statistic", middleware.Auth(true), handler.GetAllHistoryStatistic)
	history.Get("top-performance", middleware.Auth(true), handler.GetAllHistoryTopPerformance)
}
