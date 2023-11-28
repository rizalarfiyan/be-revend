package models

type VerificationSession struct {
	GoogleId    string `json:"google_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Message     string `json:"message"`
	PhoneNumber string `json:"phone_number"`
	IsError     bool   `json:"is_error"`
}

type SocialGoogle struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}
