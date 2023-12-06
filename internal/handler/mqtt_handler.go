package handler

import mqtt "github.com/eclipse/paho.mqtt.golang"

type MQTTHandler interface {
	Trigger(client mqtt.Client, msg mqtt.Message)
}
