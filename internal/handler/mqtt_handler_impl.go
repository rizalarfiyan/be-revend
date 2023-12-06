package handler

import (
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/service"
)

type mqttHandler struct {
	service service.MQTTService
}

func NewMQTTHandler(service service.MQTTService) MQTTHandler {
	return &mqttHandler{
		service: service,
	}
}

func (h *mqttHandler) Trigger(client mqtt.Client, msg mqtt.Message) {
	req := request.MQTTTriggerRequest{
		BaseMQTTRequest: request.BaseMQTTRequest{
			Client:  client,
			Message: msg,
		},
	}
	err := json.Unmarshal([]byte(msg.Payload()), &req)
	if err != nil {
		return
	}

	if !req.IsValid() {
		return
	}

	h.service.Trigger(req)
}
