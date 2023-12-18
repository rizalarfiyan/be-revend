package request

type GoogleCallbackRequest struct {
	Code        string `json:"code"`
	Scope       string `json:"scope"`
	Authuser    string `json:"authuser"`
	Prompt      string `json:"prompt"`
	State       string `json:"state"`
	ErrorReason string `json:"error_reason"`
}

type AuthVerification struct {
	Token string `json:"token" field:"Token" validate:"required" example:"2YbPyusF2G06BFQLamoKFXvGgPd"`
}

type AuthSendOTP struct {
	PhoneNumber string `json:"phone_number" field:"PhoneNumber" validate:"phone_number" example:"62895377233002"`
	Token       string `json:"token" field:"Token" validate:"omitempty,required" example:"2YbPyusF2G06BFQLamoKFXvGgPd"`
}

type AuthOTPVerification struct {
	OTP   string `json:"otp" field:"OTP" example:"651721"`
	Token string `json:"token" field:"Token" example:"2YbPyusF2G06BFQLamoKFXvGgPd"`
}

type AuthRegister struct {
	Token       string `json:"token" field:"Token" validate:"required" example:"2YbPyusF2G06BFQLamoKFXvGgPd"`
	PhoneNumber string `json:"phone_number" field:"PhoneNumber" validate:"phone_number" example:"62895377233002"`
	FirstName   string `json:"first_name" field:"FirstName" validate:"required,min=3" example:"Rizal"`
	LastName    string `json:"last_name" field:"LastName" validate:"omitempty,min=3" example:"Arfiyan"`
}
