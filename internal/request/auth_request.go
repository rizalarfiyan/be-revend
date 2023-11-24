package request

type GoogleCallbackRequest struct {
	Code        string `json:"code"`
	Scope       string `json:"scope"`
	Authuser    string `json:"authuser"`
	Prompt      string `json:"prompt"`
	ErrorReason string `json:"error_reason"`
}

type AuthVerification struct {
	Token string `json:"token" field:"Token" example:"2YbPyusF2G06BFQLamoKFXvGgPd"`
}
