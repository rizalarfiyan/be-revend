package service

type mqttService struct{}

func NewMQTTService() MQTTService {
	return &mqttService{}
}
