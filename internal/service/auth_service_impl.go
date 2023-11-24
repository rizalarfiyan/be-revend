package service

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

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

type authService struct {
	repo repository.Repository
	conf *baseModels.Config
}

func NewAuthService(repo repository.Repository) AuthService {
	return &authService{
		repo: repo,
		conf: config.Get(),
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
	var data models.SocialSession
	redirect = s.conf.Auth.Callback

	defer func() {
		token := ksuid.New().String()
		data.IsError = data.Message != ""
		data.Message = "Please wait, we are processing your request"
		err := s.repo.CreateSocialSession(ctx, token, data)
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
	data, err := s.repo.GetSocialSessionByToken(ctx, req.Token)
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
	utils.PanicIfError(err, false)

	if utils.IsEmpty(user) {
		res.Step = constants.AuthVerificationOtp
		res.Message = "Please fill the form below"
		return res
	}

	//! generate token

	return res
}

//? register
//? send otp
//? verification otp
