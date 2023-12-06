package handler

import (
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rizalarfiyan/be-revend/internal/request"
)

type mqttHandler struct{}

func NewMQTTHandler() MQTTHandler {
	return &mqttHandler{}
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

	// client.Publish("revend/action/"+req.Data.DeviceId, 0, false, "Hello World")
}
