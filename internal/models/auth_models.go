package models

import "time"

type VerificationSession struct {
	GoogleId    string `json:"google_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Message     string `json:"message"`
	PhoneNumber string `json:"phone_number"`
	IsError     bool   `json:"is_error"`
	DeviceId    string `json:"device_id"`
	Identity    string `json:"identity"`
	Token       string `json:"token"`
	IsNew       bool   `json:"is_new"`
	IsVerified  bool   `json:"is_verified"`
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

type OTPInformation struct {
	Duration  time.Duration
	Increment int64
	Data      VerificationSession
}

type OTPDetailStatus struct {
	IsBlocked     bool
	IsVerified    bool
	RemainingTime time.Duration
	Token         string
}
