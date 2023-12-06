package internal

import "github.com/rizalarfiyan/be-revend/internal/handler"

type Subscribe interface {
	BaseSubscribe(handler handler.MQTTHandler)
}
