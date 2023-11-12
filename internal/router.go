package internal

import "github.com/rizalarfiyan/be-revend/internal/handler"

type Router interface {
	BaseRoute(handler handler.BaseHandler)
}
