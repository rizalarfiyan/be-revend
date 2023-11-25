package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rizalarfiyan/be-revend/config"
	"github.com/rizalarfiyan/be-revend/constants"
	"github.com/rizalarfiyan/be-revend/internal/models"
	"github.com/rizalarfiyan/be-revend/internal/repository"
	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/response"
	"github.com/rizalarfiyan/be-revend/libs"
	baseModels "github.com/rizalarfiyan/be-revend/models"
	"github.com/rizalarfiyan/be-revend/utils"
	"github.com/segmentio/ksuid"
)

type authService struct {
	repo repository.Repository
	conf *baseModels.Config
	wa   libs.Whatsapp
}

func NewAuthService(repo repository.Repository) AuthService {
	return &authService{
		repo: repo,
		conf: config.Get(),
		wa:   libs.NewWhatsapp(),
	}
}

func (s *authService) Google() string {
	googleConf := config.Get().Auth.Google
	redirectUrl, err := url.Parse(googleConf.Endpoint.AuthURL)
	if err != nil {
		panic(err)
	}
	parameters := redirectUrl.Query()
	parameters.Add("client_id", googleConf.ClientID)
	parameters.Add("scope", strings.Join(googleConf.Scopes, " "))
	parameters.Add("redirect_uri", googleConf.RedirectURL)
	parameters.Add("response_type", "code")
	redirectUrl.RawQuery = parameters.Encode()
	return redirectUrl.String()
}

func (s *authService) GoogleCallback(ctx context.Context, req request.GoogleCallbackRequest) (redirect string) {
	var data models.VerificationSession
	redirect = s.conf.Auth.Verification.Callback

	defer func() {
		token := ksuid.New().String()
		data.IsError = data.Message != ""
		data.Message = "Please wait, we are processing your request"
		err := s.repo.CreateVerificationSession(ctx, token, data)
		if err != nil {
			return
		}

		redirectUrl, err := url.Parse(redirect)
		if err != nil {
			return
		}

		parameters := redirectUrl.Query()
		parameters.Add("token", token)
		redirectUrl.RawQuery = parameters.Encode()
		redirect = redirectUrl.String()
	}()

	if strings.EqualFold(req.ErrorReason, "user_denied") {
		data.Message = "User has denied Permission"
		return
	}

	if utils.IsEmpty(req.Code) {
		data.Message = "Code Not Found to provide AccessToken"
		return
	}

	googleConf := config.Get().Auth.Google
	googleConf.AuthCodeURL(req.Code)
	token, err := googleConf.Exchange(ctx, req.Code)
	if err != nil {
		data.Message = err.Error()
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
	if err != nil {
		data.Message = err.Error()
		return
	}

	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		data.Message = err.Error()
		return
	}

	res := models.SocialGoogle{}
	err = json.Unmarshal(response, &res)
	if err != nil {
		data.Message = err.Error()
		return
	}

	data.GoogleId = res.ID
	data.FirstName = res.GivenName
	data.LastName = res.FamilyName
	return
}

func (s *authService) Verification(ctx context.Context, req request.AuthVerification) response.AuthVerification {
	data, err := s.repo.GetVerificationSessionByToken(ctx, req.Token)
	utils.PanicIfError(err, false)
	utils.IsNotProcessMessage(data, "Token is expired", false)

	if data.IsError {
		utils.IsNotProcessRawMessage(data.Message, false)
	}

	res := response.AuthVerification{
		Message: data.Message,
	}

	if data.GoogleId == "" {
		res.Step = constants.AuthVerificationRegister
		res.Message = "You must register first, please fill the form below"
		return res
	}

	user, err := s.repo.GetUserByGoogleId(ctx, data.GoogleId)
	fmt.Println(err)
	utils.PanicIfErrorWithoutNoSqlResult(err, false)

	if utils.IsEmpty(user) {
		//! fixme with auto fill from google
		res.Step = constants.AuthVerificationRegister
		res.Message = "Please fill the form below"
		res.FirstName = data.FirstName
		res.LastName = data.LastName
		return res
	}

	payload := baseModels.AuthToken{
		FirstName:   user.FirstName,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role,
	}

	if user.LastName.Valid {
		payload.LastName = user.LastName.String
	}

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

	token, err := utils.GenerateJwtToken(claims)
	utils.PanicIfError(err, false)
	return token
}

func (s *authService) SendOTP(ctx context.Context, phoneNumber string) {
	duration, inc, err := s.repo.OTPInformation(ctx, phoneNumber)
	utils.PanicIfError(err, false)

	if !(*duration <= 0 && *inc == 0) {
		nextOtp := time.Minute * time.Duration(math.Pow(2, float64(*inc)))
		currentOtp := time.Duration(s.conf.Auth.OTP.Duration - *duration)

		var res time.Duration
		if *inc > int64(s.conf.Auth.OTP.MaxAttemp) {
			utils.IsNotProcessData("OTP has been sent, please try again in next day", res)
		}

		if currentOtp < nextOtp {
			res = nextOtp - currentOtp
			utils.IsNotProcessData("OTP has been sent, please try again in "+res.String(), res)
		}
	}

	_, err = s.repo.IncrementOTP(ctx, phoneNumber)
	utils.PanicIfError(err, false)

	otp := utils.GenerateOtp(constants.OTPLength)
	err = s.repo.CreateOTP(ctx, phoneNumber, otp)
	utils.PanicIfError(err, false)

	data := map[string]string{
		"Name": s.conf.Name,
		"Code": otp,
	}

	err = s.wa.SendMessageTemplate(phoneNumber, constants.TemplateAuthOtp, data)
	utils.PanicIfError(err, false)
}

func (s *authService) OTPVerification(ctx context.Context, req request.AuthOTPVerification) response.AuthOTPVerification {
	otp, err := s.repo.GetOTP(ctx, req.PhoneNumber)
	utils.PanicIfError(err, false)
	utils.IsNotProcess(otp, false)

	if otp != req.OTP {
		utils.IsNotProcessRawMessage("OTP is not valid", false)
	}

	if req.Token != "" {
		return s.createNewUserFromOTP(ctx, req)
	}

	user, err := s.repo.GetUserByPhoneNumber(ctx, req.PhoneNumber)
	utils.PanicIfError(err, false)

	payload := baseModels.AuthToken{
		FirstName:   user.FirstName,
		LastName:    user.LastName.String,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role,
	}

	token := s.generateToken(ctx, payload)
	return response.AuthOTPVerification{
		Token: token,
	}
}

func (s *authService) createNewUserFromOTP(ctx context.Context, req request.AuthOTPVerification) response.AuthOTPVerification {
	return response.AuthOTPVerification{}
}

//? register
//? verification otp
