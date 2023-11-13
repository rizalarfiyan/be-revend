package service

import (
	"net/url"
	"strings"

	"github.com/rizalarfiyan/be-revend/config"
)

type authService struct {
}

func NewAuthService() AuthService {
	return &authService{}
}

func (s *authService) Google() string {
	googleConf := config.Get().Auth.Google
	URL, err := url.Parse(googleConf.Endpoint.AuthURL)
	if err != nil {
		panic(err)
	}
	parameters := url.Values{}
	parameters.Add("client_id", googleConf.ClientID)
	parameters.Add("scope", strings.Join(googleConf.Scopes, " "))
	parameters.Add("redirect_uri", googleConf.RedirectURL)
	parameters.Add("response_type", "code")
	URL.RawQuery = parameters.Encode()
	return URL.String()
}
