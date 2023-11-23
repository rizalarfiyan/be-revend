package service

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/rizalarfiyan/be-revend/config"
	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/utils"
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

func (s *authService) GoogleCallback(ctx context.Context, req request.GoogleCallbackRequest) string {
	if strings.EqualFold(req.ErrorReason, "user_denied") {
		return "User has denied Permission"
	}

	if utils.IsEmpty(req.Code) {
		return "Code Not Found to provide AccessToken"
	}

	googleConf := config.Get().Auth.Google
	googleConf.AuthCodeURL(req.Code)

	token, err := googleConf.Exchange(ctx, req.Code)
	if err != nil {
		return err.Error()
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
	if err != nil {
		// fmt.Println("Get: " + err.Error() + "\n")
		return err.Error()
	}

	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		// fmt.Println("ReadAll: " + err.Error() + "\n")
		return err.Error()
	}

	_ = response

	// fmt.Println("=========================================")
	// fmt.Println("parseResponseBody: " + string(response) + "\n")
	// fmt.Println("=========================================")

	return ""
}
