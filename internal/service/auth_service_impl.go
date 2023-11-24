package service

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/rizalarfiyan/be-revend/config"
	"github.com/rizalarfiyan/be-revend/internal/models"
	"github.com/rizalarfiyan/be-revend/internal/repository"
	"github.com/rizalarfiyan/be-revend/internal/request"
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
		idx := ksuid.New().String()
		err := s.repo.CreateSocialSession(ctx, idx, data)
		if err != nil {
			return
		}

		redirectUrl, err := url.Parse(redirect)
		if err != nil {
			return
		}

		parameters := redirectUrl.Query()
		parameters.Add("token", idx)
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
