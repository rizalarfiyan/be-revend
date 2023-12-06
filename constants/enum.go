package constants

type AuthVerificationStep int

const (
	AuthVerificationRegister AuthVerificationStep = iota + 1
	AuthVerificationOtp
	AuthVerificationDone
)

type BaseMQTTActionStep int

const (
	MQTTStepCancel BaseMQTTActionStep = iota
	MQTTStepCheckUser
	MQTTStepStatus
)

func (v BaseMQTTActionStep) IsValid() bool {
	val := BaseMQTTActionStep(v)
	if val >= MQTTStepCancel && val <= MQTTStepStatus {
		return true
	}

	return false
}
