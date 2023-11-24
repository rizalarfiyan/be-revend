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
