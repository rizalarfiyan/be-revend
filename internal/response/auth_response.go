package response

import (
	"github.com/rizalarfiyan/be-revend/constants"
)

type AuthVerification struct {
	Step      constants.AuthVerificationStep `json:"step"`
	Token     string                         `json:"token"`
	FirstName string                         `json:"first_name"`
	LastName  string                         `json:"last_name"`
	Message   string                         `json:"-"`
}

type AuthOTPVerification struct {
	Token string `json:"token"`
}
