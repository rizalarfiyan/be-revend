package service

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/rizalarfiyan/be-revend/config"
	"github.com/rizalarfiyan/be-revend/constants"
	"github.com/rizalarfiyan/be-revend/internal/models"
	"github.com/rizalarfiyan/be-revend/internal/repository"
	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/response"
	baseModels "github.com/rizalarfiyan/be-revend/models"
	"github.com/rizalarfiyan/be-revend/utils"
	"github.com/segmentio/ksuid"
)

type mqttService struct {
	authRepo repository.AuthRepository
	conf     *baseModels.Config
	utils    utils.MQTTUtils
}

func NewMQTTService(authRepo repository.AuthRepository) MQTTService {
	return &mqttService{
		authRepo: authRepo,
		conf:     config.Get(),
		utils:    utils.NewMqttUtils(),
	}
}

func (s *mqttService) Trigger(req request.MQTTTriggerRequest) {
	defer s.utils.Recover()

	ctx := context.Background()
	switch req.Step {
	case constants.MQTTStepCancel:
		s.cancelRequest(ctx, req)
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
	user, err := s.authRepo.GetUserByIdentity(ctx, req.Data.Identity)
	s.utils.PanicIfErrorWithoutNoSqlResult(err)

	if !utils.IsEmpty(user) {
		s.sendTopic(req, response.MQTTActionResponse{
			Step: constants.MQTTStepCheckUser,
			Data: response.MQTTActionDataResponse{
				State: constants.MQTTCheckUserLogin,
			},
		})
		return
	}

	token := ksuid.New().String()
	payload := models.VerificationSession{
		Identity: req.Data.Identity,
		DeviceId: req.Data.DeviceId,
		IsNew:    true,
		Token:    token,
		Message:  "You must register first, please fill the form below",
	}

	err = s.authRepo.CreateVerificationSession(ctx, token, payload)
	s.utils.PanicIfError(err)

	redirectUrl, err := url.Parse(s.conf.Auth.Verification.Callback)
	if err != nil {
		return
	}

	parameters := redirectUrl.Query()
	parameters.Add("token", token)
	redirectUrl.RawQuery = parameters.Encode()

	s.sendTopic(req, response.MQTTActionResponse{
		Step: constants.MQTTStepCheckUser,
		Data: response.MQTTActionDataResponse{
			State: constants.MQTTCheckUserMustRegister,
			Link:  redirectUrl.String(),
		},
	})
}

func (s *mqttService) cancelRequest(ctx context.Context, req request.MQTTTriggerRequest) {
	data, err := s.authRepo.GetVerificationSessionByIdentity(ctx, req.Data.Identity)
	s.utils.PanicIfError(err)

	err = s.authRepo.DeleteVerificationSessionByIdentity(ctx, req.Data.Identity)
	s.utils.PanicIfError(err)

	if !utils.IsEmpty(data) && !utils.IsEmpty(data.PhoneNumber) {
		err = s.authRepo.DeleteAllOTP(ctx, data.PhoneNumber)
		s.utils.PanicIfError(err)
	}

	s.sendTopic(req, response.MQTTActionResponse{
		Step: constants.MQTTStepCancel,
	})
}

func (s *mqttService) sendTopic(req request.MQTTTriggerRequest, payload response.MQTTActionResponse) {
	bytePayload, err := json.Marshal(payload)
	s.utils.PanicIfError(err)

	topic := "revend/action/" + req.Data.DeviceId
	req.Client.Publish(topic, 0, false, bytePayload)
}
