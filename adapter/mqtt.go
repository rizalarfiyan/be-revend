package adapter

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rizalarfiyan/be-revend/config"
	"github.com/rizalarfiyan/be-revend/logger"
)

var mqttConn *mqtt.Client
var mqttLog logger.Logger

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	mqttLog.Error(err, "MQTT connection lost")
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	mqttLog.Info("MQTT connected")
}

func MQTTInit() {
	mqttLog = logger.Get("mqtt")
	mqttLog.Info("Connect MQTT server...")
	conf := config.Get()
	dsn := fmt.Sprintf("%s:%d", conf.MQTT.Server, conf.MQTT.Port)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(dsn)
	opts.SetClientID(conf.MQTT.ClientId)
	opts.AutoReconnect = true
	opts.OnConnectionLost = connectLostHandler
	opts.OnConnect = connectHandler

	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		mqttLog.Fatal(token.Error(), "MQTT connection problem")
	}

	mqttConn = new(mqtt.Client)
	mqttConn = &client
}

func MQTTConnection() *mqtt.Client {
	return mqttConn
}

func MQTTIsConnected() bool {
	if mqttConn == nil {
		return false
	}
	connected := (*mqttConn).IsConnected()
	if !connected {
		mqttLog.Warn("MQTT connection lost")
	}
	return connected
}
