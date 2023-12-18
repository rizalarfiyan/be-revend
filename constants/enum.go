package constants

type AuthVerificationStep int

const (
	AuthVerificationRegister AuthVerificationStep = iota + 1
	AuthVerificationOtp
	AuthVerificationDone
)

type BaseMQTTActionStep int

const (
	MQTTStepCancel BaseMQTTActionStep = iota + 1
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

type MQTTCheckUserState int

const (
	MQTTCheckUserLogin MQTTCheckUserState = iota + 1
	MQTTCheckUserMustRegister
	MQTTCheckUserSuccessRegister
)

type FilterListStatus string

const (
	FilterListStatusActive  FilterListStatus = "active"
	FilterListStatusDeleted FilterListStatus = "deleted"
)

func (v FilterListStatus) IsValid() bool {
	switch v {
	case FilterListStatusActive, FilterListStatusDeleted:
		return true
	}

	return false
}

type FilterTimeFrequency string

const (
	FilterTimeFrequencyToday   FilterTimeFrequency = "today"
	FilterTimeFrequencyWeek    FilterTimeFrequency = "week"
	FilterTimeFrequencyMonth   FilterTimeFrequency = "month"
	FilterTimeFrequencyQuarter FilterTimeFrequency = "quarter"
	FilterTimeFrequencyYear    FilterTimeFrequency = "year"
)

func (v FilterTimeFrequency) IsValid() bool {
	switch v {
	case FilterTimeFrequencyToday, FilterTimeFrequencyWeek, FilterTimeFrequencyMonth, FilterTimeFrequencyQuarter, FilterTimeFrequencyYear:
		return true
	}

	return false
}
