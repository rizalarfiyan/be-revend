package service

import "github.com/rizalarfiyan/be-revend/internal/request"

type MQTTService interface {
	Trigger(req request.MQTTTriggerRequest)
}
