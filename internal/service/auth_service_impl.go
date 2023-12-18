package service

import (
	"context"
	"encoding/json"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rizalarfiyan/be-revend/config"
	"github.com/rizalarfiyan/be-revend/constants"
	"github.com/rizalarfiyan/be-revend/exception"
	"github.com/rizalarfiyan/be-revend/internal/models"
	"github.com/rizalarfiyan/be-revend/internal/repository"
	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/response"
	"github.com/rizalarfiyan/be-revend/internal/sql"
	"github.com/rizalarfiyan/be-revend/libs"
	baseModels "github.com/rizalarfiyan/be-revend/models"
	"github.com/rizalarfiyan/be-revend/utils"
	"github.com/segmentio/ksuid"
)

type authService struct {
	repo     repository.AuthRepository
	userRepo repository.UserRepository
	conf     *baseModels.Config
	wa       libs.Whatsapp
	mqtt     mqtt.Client
}

func NewAuthService(repo repository.AuthRepository, userRepo repository.UserRepository, mqtt mqtt.Client) AuthService {
	return &authService{
		repo:     repo,
		userRepo: userRepo,
		conf:     config.Get(),
		wa:       libs.NewWhatsapp(),
		mqtt:     mqtt,
	}
}

func (s *authService) Google(state ...string) string {
	googleConf := config.Get().Auth.Google
	redirectUrl, err := url.Parse(googleConf.Endpoint.AuthURL)
	exception.PanicIfError(err, false)
	parameters := redirectUrl.Query()
	parameters.Add("client_id", googleConf.ClientID)
	parameters.Add("scope", strings.Join(googleConf.Scopes, " "))
	parameters.Add("redirect_uri", googleConf.RedirectURL)
	parameters.Add("response_type", "code")
	if len(state) > 0 {
		parameters.Add("state", state[0])
	}
	redirectUrl.RawQuery = parameters.Encode()
	return redirectUrl.String()
}

func (s *authService) GoogleCallback(ctx context.Context, req request.GoogleCallbackRequest) (redirect string) {
	isVerification := utils.IsEmpty(req.State)
	var res models.SocialGoogle
	var message string

	defer func() {
		if isVerification {
			redirect = s.googleCallbackVerificationSession(ctx, res, message)
			return
		}

		redirect = s.googleCallbackBindUser(ctx, res, req.State, message)
	}()

	if strings.EqualFold(req.ErrorReason, "user_denied") {
		message = "User has denied Permission"
		return
	}

	if utils.IsEmpty(req.Code) {
		message = "Code Not Found to provide AccessToken"
		return
	}

	googleConf := config.Get().Auth.Google
	googleConf.AuthCodeURL(req.Code)
	token, err := googleConf.Exchange(ctx, req.Code)
	if err != nil {
		message = err.Error()
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
	if err != nil {
		message = err.Error()
		return
	}

	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		message = err.Error()
		return
	}

	err = json.Unmarshal(response, &res)
	if err != nil {
		message = err.Error()
		return
	}

	return
}

func (s *authService) googleCallbackVerificationSession(ctx context.Context, res models.SocialGoogle, message string) (redirect string) {
	redirect = s.conf.Auth.Verification.Callback
	data := models.VerificationSession{
		Message: message,
	}

	if !utils.IsEmpty(res) {
		data.GoogleId = res.ID
		user, err := s.userRepo.GetUserByGoogleId(ctx, data.GoogleId)
		exception.PanicIfErrorWithoutNoSqlResult(err, false)

		if utils.IsEmpty(user) {
			data.Message = "You must register first or bind your account"
			return
		}

		data.FirstName = user.FirstName
		if user.LastName.Valid {
			data.LastName = user.LastName.String
		}

		data.Identity = user.Identity
		data.PhoneNumber = user.PhoneNumber
		data.IsVerified = true
	}

	token := ksuid.New().String()
	data.IsError = data.Message != ""

	if !data.IsError {
		data.Message = "Please wait, we are processing your request"
	}

	data.Token = token
	err := s.repo.CreateVerificationSession(ctx, token, data)
	if err != nil {
		return
	}

	redirectUrl, err := url.Parse(s.conf.Auth.Verification.Callback)
	if err != nil {
		return
	}

	parameters := redirectUrl.Query()
	parameters.Add("token", token)
	redirectUrl.RawQuery = parameters.Encode()
	redirect = redirectUrl.String()
	return
}

func (s *authService) googleCallbackBindUser(ctx context.Context, res models.SocialGoogle, token, message string) (redirect string) {
	redirect = s.conf.Auth.Verification.ProfileCallback
	var userId uuid.UUID

	defer func() {
		if userId != uuid.Nil {
			err := s.userRepo.DeleteBindGoogle(ctx, userId)
			if err != nil {
				return
			}
		}

		if message == "" {
			return
		}

		redirectUrl, err := url.Parse(redirect)
		if err != nil {
			return
		}

		parameters := redirectUrl.Query()
		parameters.Add("message", message)
		redirectUrl.RawQuery = parameters.Encode()
		redirect = redirectUrl.String()
	}()

	_, err := ksuid.Parse(token)
	if err != nil {
		message = "Token is not valid"
		return
	}

	rawUserId, err := s.userRepo.GetBindGoogle(ctx, token)
	if err != nil {
		return
	}

	if utils.IsEmpty(rawUserId) {
		message = "Token is not valid"
		return
	}

	userId, err = uuid.Parse(rawUserId)
	if err != nil {
		return
	}

	user, err := s.userRepo.GetUserById(ctx, userId)
	if err != nil {
		return
	}

	if !utils.IsEmpty(user) {
		message = "User already bind with google"
		return
	}

	payload := sql.UpdateGoogleUserProfileParams{
		GoogleID: utils.PGText(res.ID),
		ID:       utils.PGUUID(userId),
	}

	err = s.userRepo.UpdateGoogleUserProfile(ctx, payload)
	if err != nil {
		return
	}

	return
}

func (s *authService) GetSession(ctx context.Context, token string) models.VerificationSession {
	data, err := s.repo.GetVerificationSessionByToken(ctx, token)
	exception.PanicIfError(err, false)

	exception.IsNotProcessMessage(data, "Session is expired", false)
	return *data
}

func (s *authService) Verification(ctx context.Context, req request.AuthVerification) response.AuthVerification {
	data := s.GetSession(ctx, req.Token)
	if data.IsError {
		exception.IsNotProcessRawMessage(data.Message, false)
	}

	res := response.AuthVerification{
		Message: data.Message,
	}

	if data.IsError {
		return res
	}

	if utils.IsEmpty(data.Identity) {
		res.Message = "Something wrong for your request"
		return res
	}

	if utils.IsEmpty(data.PhoneNumber) {
		res.Step = constants.AuthVerificationRegister
		res.FirstName = data.FirstName
		res.LastName = data.LastName
		return res
	}

	otp := s.getOTPDetail(ctx, data.PhoneNumber)
	if (otp.IsBlocked || otp.RemainingTime > 0) && !data.IsVerified {
		res.Step = constants.AuthVerificationOtp
		res.PhoneNumber = data.PhoneNumber
		res.RemainingTime = int64(otp.RemainingTime.Seconds())
		return res
	}

	user, err := s.userRepo.GetUserByGoogleIdOrPhoneNumber(ctx, data.GoogleId, data.PhoneNumber)
	exception.PanicIfErrorWithoutNoSqlResult(err, false)

	if utils.IsEmpty(user) {
		res.Message = "Something wrong for your request"
		return res
	}

	payload := baseModels.AuthToken{
		Id:          user.ID.Bytes,
		FirstName:   user.FirstName,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role,
	}

	if user.LastName.Valid {
		payload.LastName = user.LastName.String
	}

	err = s.repo.DeleteVerificationSessionByToken(ctx, req.Token)
	exception.PanicIfError(err, false)

	err = s.repo.DeleteAllOTP(ctx, data.PhoneNumber)
	exception.PanicIfError(err, false)

	token := s.generateToken(ctx, payload)
	res.Step = constants.AuthVerificationDone
	res.Message = "Please wait, you are being redirected"
	res.Token = token

	return res
}

func (s *authService) generateToken(ctx context.Context, payload baseModels.AuthToken) string {
	conf := config.Get()
	claims := baseModels.AuthTokenClaims{
		AuthToken: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(conf.JWT.Expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token, err := utils.GenerateJwtToken(conf.JWT.Secret, claims)
	exception.PanicIfError(err, false)
	return token
}

func (s *authService) createAndSendOTP(ctx context.Context, phoneNumber, token string, payload models.VerificationSession) {
	otp := utils.GenerateOtp(constants.OTPLength)
	data := map[string]string{
		"Name": s.conf.Name,
		"Code": otp,
	}

	err := s.wa.SendMessageTemplate(phoneNumber, constants.TemplateAuthOtp, data)
	exception.PanicIfError(err, false)

	err = s.repo.CreateVerificationSession(ctx, token, payload)
	exception.PanicIfError(err, false)

	_, err = s.repo.IncrementOTP(ctx, phoneNumber)
	exception.PanicIfError(err, false)

	err = s.repo.CreateOTP(ctx, phoneNumber, otp)
	exception.PanicIfError(err, false)
}

func (s *authService) getOTPDetail(ctx context.Context, phoneNumber string) models.OTPDetailStatus {
	otp, err := s.repo.OTPInformation(ctx, phoneNumber)
	exception.PanicIfError(err, false)

	res := models.OTPDetailStatus{
		Token: otp.Data.Token,
	}

	if !(otp.Duration <= 0 && otp.Increment == 0) {
		nextOtp := time.Minute * time.Duration(math.Pow(2, float64(otp.Increment)))
		currentOtp := time.Duration(s.conf.Auth.OTP.Duration - otp.Duration)

		if otp.Increment > int64(s.conf.Auth.OTP.MaxAttemp) {
			res.IsBlocked = true
		}

		if currentOtp < nextOtp {
			res.RemainingTime = nextOtp - currentOtp
		}
	}

	return res
}

func (s *authService) SendOTP(ctx context.Context, req request.AuthSendOTP) response.AuthSendOTP {
	otp := s.getOTPDetail(ctx, req.PhoneNumber)
	if otp.IsBlocked {
		exception.IsNotProcessData("OTP has been sent, please try again in next day", response.AuthSendOTP{
			Token: otp.Token,
		})
	}

	if otp.RemainingTime > 0 {
		exception.IsNotProcessData("OTP has been sent, please try again in "+otp.RemainingTime.String(), response.AuthSendOTP{
			Token: otp.Token,
		})
	}

	user, err := s.userRepo.GetUserByPhoneNumber(ctx, req.PhoneNumber)
	exception.PanicIfErrorWithoutNoSqlResult(err, false)

	if utils.IsEmpty(user) {
		exception.ErrorManualValidation("phone_number", "Unable to process phone number")
	}

	token := ksuid.New().String()
	if !utils.IsEmpty(req.Token) {
		token = req.Token
	}
	payload := models.VerificationSession{
		PhoneNumber: req.PhoneNumber,
		FirstName:   user.FirstName,
		Identity:    user.Identity,
		Token:       token,
	}

	if user.GoogleID.Valid {
		payload.GoogleId = user.GoogleID.String
	}

	if user.LastName.Valid {
		payload.LastName = user.LastName.String
	}

	s.createAndSendOTP(ctx, req.PhoneNumber, token, payload)
	return response.AuthSendOTP{
		Token: token,
	}
}

func (s *authService) OTPVerification(ctx context.Context, req request.AuthOTPVerification) response.AuthOTPVerification {
	data := s.GetSession(ctx, req.Token)

	otp, err := s.repo.GetOTP(ctx, data.PhoneNumber)
	exception.PanicIfError(err, false)
	exception.IsNotProcess(otp, false)

	if otp != req.OTP {
		exception.IsNotProcessRawMessage("OTP is not valid", false)
	}

	user, err := s.userRepo.GetUserByPhoneNumber(ctx, data.PhoneNumber)
	exception.PanicIfErrorWithoutNoSqlResult(err, false)

	if utils.IsEmpty(user) {
		payload := sql.CreateUserParams{
			FirstName:   data.FirstName,
			LastName:    pgtype.Text{String: data.LastName, Valid: data.LastName != ""},
			PhoneNumber: data.PhoneNumber,
			GoogleID:    pgtype.Text{String: data.GoogleId, Valid: data.GoogleId != ""},
			Identity:    data.Identity,
		}
		err = s.userRepo.CreateUser(ctx, payload)
		exception.PanicIfError(err, false)

		user, err = s.userRepo.GetUserByPhoneNumber(ctx, data.PhoneNumber)
		exception.PanicIfErrorWithoutNoSqlResult(err, false)
	}

	payload := baseModels.AuthToken{
		Id:          user.ID.Bytes,
		FirstName:   user.FirstName,
		LastName:    user.LastName.String,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role,
	}

	err = s.repo.DeleteVerificationSessionByToken(ctx, req.Token)
	exception.PanicIfError(err, false)

	err = s.repo.DeleteAllOTP(ctx, data.PhoneNumber)
	exception.PanicIfError(err, false)

	if data.IsNew {
		mqttPayload := response.MQTTActionResponse{
			Step: constants.MQTTStepCheckUser,
			Data: response.MQTTActionDataResponse{
				State: constants.MQTTCheckUserSuccessRegister,
			},
		}

		mqttBytePayload, err := json.Marshal(mqttPayload)
		exception.PanicIfError(err, false)

		topic := "revend/action/" + data.DeviceId
		s.mqtt.Publish(topic, 0, false, mqttBytePayload)
	}

	token := s.generateToken(ctx, payload)
	return response.AuthOTPVerification{
		Token: token,
	}
}

func (s *authService) Register(ctx context.Context, req request.AuthRegister) {
	user, err := s.userRepo.GetUserByPhoneNumber(ctx, req.PhoneNumber)
	exception.PanicIfErrorWithoutNoSqlResult(err, false)

	if !utils.IsEmpty(user) {
		exception.ErrorManualValidation("phone_number", "Phone Number already exist.")
	}

	data := s.GetSession(ctx, req.Token)
	data.PhoneNumber = req.PhoneNumber
	data.FirstName = req.FirstName
	data.LastName = req.LastName
	err = s.repo.CreateVerificationSession(ctx, req.Token, data)
	exception.PanicIfError(err, false)
	s.createAndSendOTP(ctx, req.PhoneNumber, req.Token, data)
}
