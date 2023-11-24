package libs

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/rizalarfiyan/be-revend/config"
	baseModels "github.com/rizalarfiyan/be-revend/models"
	"github.com/rizalarfiyan/be-revend/utils"
)

type Whatsapp interface {
	SendMessage(phoneNumber, message string) error
	SendMessageTemplate(phoneNumber, template string, data interface{}) error
}

type whatsapp struct {
	conf    *baseModels.Config
	timeout time.Duration
}

func NewWhatsapp() Whatsapp {
	return &whatsapp{
		conf:    config.Get(),
		timeout: 5 * time.Second,
	}
}

func (w *whatsapp) parseUrl(path ...string) (string, error) {
	apiUrl, err := url.Parse(w.conf.Whatsapp.ApiUrl)
	if err != nil {
		return "", err
	}

	return apiUrl.JoinPath(path...).String(), nil
}

func (w *whatsapp) send(method, apiUrl string, payload map[string]string) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, apiUrl, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(context.Background(), w.timeout)
	defer cancel()

	req = req.WithContext(ctx)
	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New("something wrong send whatsapp")
	}

	return nil
}

func (w *whatsapp) SendMessage(phoneNumber, message string) error {
	apiUrl, err := w.parseUrl("send", "message")
	if err != nil {
		return err
	}

	payload := map[string]string{
		"phone":   utils.WhatsappPhoneNumber(phoneNumber),
		"message": message,
	}

	return w.send(http.MethodPost, apiUrl, payload)
}

func (w *whatsapp) SendMessageTemplate(phoneNumber, template string, data interface{}) error {
	message, err := utils.ParseTextTemplate(template, data)
	if err != nil {
		return err
	}

	return w.SendMessage(phoneNumber, message)
}
