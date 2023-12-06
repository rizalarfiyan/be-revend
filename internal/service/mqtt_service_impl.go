package service

import (
	"context"
	"encoding/json"

	"github.com/rizalarfiyan/be-revend/constants"
	"github.com/rizalarfiyan/be-revend/internal/repository"
	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/response"
	"github.com/rizalarfiyan/be-revend/utils"
)

type mqttService struct {
	authRepo repository.Repository
	utils    utils.MQTTUtils
}

func NewMQTTService(authRepo repository.Repository) MQTTService {
	return &mqttService{
		authRepo: authRepo,
		utils:    utils.NewMqttUtils(),
	}
}

func (s *mqttService) Trigger(req request.MQTTTriggerRequest) {
	defer s.utils.Recover()

	ctx := context.Background()
	switch req.Step {
	case constants.MQTTStepCancel:
		return
	case constants.MQTTStepCheckUser:
		s.checkUser(ctx, req)
		return
	case constants.MQTTStepStatus:
		return
	default:
		return
	}
}

func (s *mqttService) checkUser(ctx context.Context, req request.MQTTTriggerRequest) {
	user, err := s.authRepo.GetUserByIdentity(ctx, req.Data.Identifier)
	s.utils.PanicIfErrorWithoutNoSqlResult(err)

	if utils.IsEmpty(user) {
		//! create session user, send mqtt must register
		return
	}

	s.sendTopic(req, response.MQTTActionResponse{
		Step: constants.MQTTStepCheckUser,
		Data: response.MQTTActionDataResponse{
			State: constants.MQTTCheckUserLogin,
		},
	})
}

func (s *mqttService) sendTopic(req request.MQTTTriggerRequest, payload response.MQTTActionResponse) {
	bytePayload, err := json.Marshal(payload)
	s.utils.PanicIfError(err)

	topic := "revend/action/" + req.Data.DeviceId
	req.Client.Publish(topic, 0, false, bytePayload)
}
