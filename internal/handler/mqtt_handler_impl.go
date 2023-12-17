package handler

import (
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/service"
	"github.com/rizalarfiyan/be-revend/logger"
	"github.com/rs/zerolog"
)

type mqttHandler struct {
	service service.MQTTService
	log     *zerolog.Logger
}

func NewMQTTHandler(service service.MQTTService) MQTTHandler {
	return &mqttHandler{
		service: service,
		log:     logger.Get("mqtt-subscribe").Logs(),
	}
}

func (h *mqttHandler) Trigger(client mqtt.Client, msg mqtt.Message) {
	h.log.Info().Str("topic", msg.Topic()).RawJSON("payload", msg.Payload()).Msg("subscribe trigger handler")

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
