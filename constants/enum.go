package constants

type AuthVerificationStep int

const (
	AuthVerificationRegister AuthVerificationStep = iota + 1
	AuthVerificationOtp
	AuthVerificationDone
)
