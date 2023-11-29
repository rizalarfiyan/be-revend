package response

import (
	"github.com/rizalarfiyan/be-revend/constants"
)

type AuthVerification struct {
	Step          constants.AuthVerificationStep `json:"step"`
	Token         string                         `json:"token"`
	FirstName     string                         `json:"first_name"`
	LastName      string                         `json:"last_name"`
	PhoneNumber   string                         `json:"phone_number"`
	RemainingTime int64                          `json:"remaining_time"`
	Message       string                         `json:"-"`
}

type AuthOTPVerification struct {
	Token string `json:"token"`
}

type AuthSendOTP struct {
	Token string `json:"token"`
}
