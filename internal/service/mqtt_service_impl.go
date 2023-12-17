package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/rizalarfiyan/be-revend/config"
	"github.com/rizalarfiyan/be-revend/constants"
	"github.com/rizalarfiyan/be-revend/internal/models"
	"github.com/rizalarfiyan/be-revend/internal/repository"
	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/response"
	"github.com/rizalarfiyan/be-revend/internal/sql"
	"github.com/rizalarfiyan/be-revend/libs"
	"github.com/rizalarfiyan/be-revend/logger"
	baseModels "github.com/rizalarfiyan/be-revend/models"
	"github.com/rizalarfiyan/be-revend/utils"
	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"
)

type mqttService struct {
	repo        repository.MqttRepository
	authRepo    repository.AuthRepository
	userRepo    repository.UserRepository
	deviceRepo  repository.DeviceRepository
	historyRepo repository.HistoryRepository
	conf        *baseModels.Config
	utils       utils.MQTTUtils
	wa          libs.Whatsapp
	log         *zerolog.Logger
}

func NewMQTTService(repo repository.MqttRepository, authRepo repository.AuthRepository, userRepo repository.UserRepository, deviceRepo repository.DeviceRepository, historyRepo repository.HistoryRepository) MQTTService {
	return &mqttService{
		repo:        repo,
		authRepo:    authRepo,
		userRepo:    userRepo,
		deviceRepo:  deviceRepo,
		historyRepo: historyRepo,
		conf:        config.Get(),
		utils:       utils.NewMqttUtils(),
		wa:          libs.NewWhatsapp(),
		log:         logger.Get("mqtt-subscribe").Logs(),
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
		s.creatOrUpdateUserPoint(ctx, req)
		return
	default:
		return
	}
}

func (s *mqttService) checkUser(ctx context.Context, req request.MQTTTriggerRequest) {
	user, err := s.userRepo.GetUserByIdentity(ctx, req.Data.Identity)
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
	verification, err := s.authRepo.GetVerificationSessionByIdentity(ctx, req.Data.Identity)
	s.utils.PanicIfError(err)

	err = s.authRepo.DeleteVerificationSessionByIdentity(ctx, req.Data.Identity)
	s.utils.PanicIfError(err)

	if !utils.IsEmpty(verification) && !utils.IsEmpty(verification.PhoneNumber) {
		err = s.authRepo.DeleteAllOTP(ctx, verification.PhoneNumber)
		s.utils.PanicIfError(err)
	}

	userPoint, err := s.repo.GetUserPoint(ctx, req.Data.Identity)
	s.utils.PanicIfError(err)

	if userPoint == nil {
		s.sendTopic(req, response.MQTTActionResponse{
			Step: constants.MQTTStepCancel,
		})
        return
	}

	user, err := s.userRepo.GetUserByIdentity(ctx, req.Data.Identity)
	s.utils.PanicIfErrorWithoutNoSqlResult(err)

	if utils.IsEmpty(user) {
		s.sendTopic(req, response.MQTTActionResponse{
			Step: constants.MQTTStepCancel,
		})
        return
	}

	device, err := s.deviceRepo.GetDeviceByToken(ctx, req.Data.DeviceId)
	s.utils.PanicIfErrorWithoutNoSqlResult(err)

	if utils.IsEmpty(device) {
		s.sendTopic(req, response.MQTTActionResponse{
			Step: constants.MQTTStepCancel,
		})
        return
	}

	payload := sql.CreateHistoryParams{
		UserID:   user.ID,
		DeviceID: device.ID,
		Success:  int32(userPoint.Success),
		Failed:   int32(userPoint.Failed),
	}
	s.historyRepo.CreateHistory(ctx, payload)
	s.utils.PanicIfError(err)

	data := map[string]string{
		"Name":    user.FirstName,
		"Device":  device.Name,
		"Success": fmt.Sprint(payload.Success),
		"Failed":  fmt.Sprint(payload.Failed),
	}
	s.wa.SendMessageTemplate(user.PhoneNumber, constants.TemplateHistory, data)

	err = s.repo.DeleteUserPoint(ctx, req.Data.Identity)
	s.utils.PanicIfError(err)

	s.sendTopic(req, response.MQTTActionResponse{
		Step: constants.MQTTStepCancel,
	})
}

func (s *mqttService) creatOrUpdateUserPoint(ctx context.Context, req request.MQTTTriggerRequest) {
	err := s.repo.CreateOrUpdateUserPoint(ctx, models.UserPoint{
		Identity: req.Data.Identity,
		DeviceId: req.Data.DeviceId,
		Success:  req.Data.Success,
		Failed:   req.Data.Failed,
	})
	s.utils.PanicIfError(err)

	s.sendTopic(req, response.MQTTActionResponse{
		Step: constants.MQTTStepStatus,
	})
}

func (s *mqttService) sendTopic(req request.MQTTTriggerRequest, payload response.MQTTActionResponse) {
	bytePayload, err := json.Marshal(payload)
	s.utils.PanicIfError(err)

	topic := "revend/action/" + req.Data.DeviceId
	req.Client.Publish(topic, 0, false, bytePayload)
	s.log.Info().Str("topic", topic).RawJSON("payload", bytePayload).Msg("send topic")
}
