package request

type GoogleCallbackRequest struct {
	Code        string `json:"code"`
	Scope       string `json:"scope"`
	Authuser    string `json:"authuser"`
	Prompt      string `json:"prompt"`
	ErrorReason string `json:"error_reason"`
}

type AuthVerification struct {
	Token string `json:"token" field:"Token" validate:"required" example:"2YbPyusF2G06BFQLamoKFXvGgPd"`
}

type AuthSendOTP struct {
	PhoneNumber string `json:"phone_number" field:"PhoneNumber" validate:"phone_number" example:"62895377233002"`
}

type AuthOTPVerification struct {
	PhoneNumber string `json:"phone_number" field:"PhoneNumber" validate:"phone_number" example:"62895377233002"`
	OTP         string `json:"otp" field:"OTP" example:"651721"`
	Token       string `json:"token" field:"Token" example:"2YbPyusF2G06BFQLamoKFXvGgPd"`
}
