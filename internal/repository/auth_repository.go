package repository

import (
	"context"

	"github.com/rizalarfiyan/be-revend/internal/models"
)

type AuthRepository interface {
	CreateVerificationSession(ctx context.Context, idx string, payload models.VerificationSession) error
	GetVerificationSessionByToken(ctx context.Context, token string) (*models.VerificationSession, error)
	GetVerificationSessionByPhoneNumber(ctx context.Context, phoneNumber string) (*models.VerificationSession, error)
	GetVerificationSessionByIdentity(ctx context.Context, identity string) (*models.VerificationSession, error)
	DeleteVerificationSessionByGoogleId(ctx context.Context, googleId string) error
	DeleteVerificationSessionByToken(ctx context.Context, token string) error
	DeleteVerificationSessionByIdentity(ctx context.Context, identity string) error
	IncrementOTP(ctx context.Context, phoneNumber string) (int64, error)
	CreateOTP(ctx context.Context, phoneNumber, otp string) error
	OTPInformation(ctx context.Context, phoneNumber string) (*models.OTPInformation, error)
	GetOTP(ctx context.Context, phoneNumber string) (string, error)
	DeleteAllOTP(ctx context.Context, phoneNumber string) error
}
