package response

import "github.com/rizalarfiyan/be-revend/constants"

type MQTTActionResponse struct {
	Step constants.BaseMQTTActionStep `json:"step"`
	Data MQTTActionDataResponse       `json:"data"`
}

type MQTTActionDataResponse struct {
	State constants.MQTTCheckUserState `json:"state"`
	Link  string                       `json:"link"`
}
