package response

import (
	"github.com/rizalarfiyan/be-revend/constants"
)

type AuthVerification struct {
	Step    constants.AuthVerificationStep `json:"step"`
	Token   string                         `json:"token"`
	Message string                         `json:"-"`
}
