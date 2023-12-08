package request

import (
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rizalarfiyan/be-revend/constants"
)

type BaseMQTTRequest struct {
	Client  mqtt.Client  `json:"-"`
	Message mqtt.Message `json:"-"`
}

type MQTTTriggerRequest struct {
	BaseMQTTRequest
	Step constants.BaseMQTTActionStep `json:"step"`
	Data MQTTTriggerDataRequest       `json:"data"`
}

type MQTTTriggerDataRequest struct {
	DeviceId string `json:"device_id"`
	Identity string `json:"identity"`
	Failed   int    `json:"failed"`
	Success  int    `json:"success"`
}

func (bs *MQTTTriggerRequest) IsValid() bool {
	bs.Data.DeviceId = strings.ToLower(strings.TrimSpace(bs.Data.DeviceId))
	if bs.Data.DeviceId == "" || len(bs.Data.DeviceId) != 42 {
		return false
	}

	if strings.TrimSpace(bs.Data.Identity) == "" {
		return false
	}

	return bs.Step.IsValid()
}
