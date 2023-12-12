package internal

import "github.com/rizalarfiyan/be-revend/internal/handler"

type Router interface {
	BaseRoute(handler handler.BaseHandler)
	AuthRoute(handler handler.AuthHandler)
	UserRoute(handler handler.UserHandler)
	DeviceRoute(handler handler.DeviceHandler)
	HistoryRoute(handler handler.HistoryHandler)
}
