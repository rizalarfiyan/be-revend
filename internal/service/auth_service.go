package service

import (
	"context"

	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/response"
)

type AuthService interface {
	Google() string
	GoogleCallback(ctx context.Context, req request.GoogleCallbackRequest) string
	Verification(ctx context.Context, req request.AuthVerification) response.AuthVerification
	SendOTP(ctx context.Context, req request.AuthSendOTP) response.AuthSendOTP
	OTPVerification(ctx context.Context, req request.AuthOTPVerification) response.AuthOTPVerification
	Register(ctx context.Context, req request.AuthRegister)
}
