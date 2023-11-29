package repository

import (
	"context"
	"time"

	"github.com/rizalarfiyan/be-revend/internal/models"
	"github.com/rizalarfiyan/be-revend/internal/sql"
)

type Repository interface {
	GetUserByGoogleId(ctx context.Context, googleID string) (sql.User, error)
	GetUserByPhoneNumber(ctx context.Context, googleID string) (sql.User, error)
	GetUserByGoogleIdOrPhoneNumber(ctx context.Context, googleID, phoneNumber string) (sql.User, error)
	CreateUser(ctx context.Context, payload sql.CreateUserParams) error
	CreateVerificationSession(ctx context.Context, idx string, payload models.VerificationSession) error
	GetVerificationSessionByToken(ctx context.Context, token string) (*models.VerificationSession, error)
	DeleteVerificationSessionByGoogleId(ctx context.Context, googleId string) error
	DeleteVerificationSessionByToken(ctx context.Context, token string) error
	IncrementOTP(ctx context.Context, phoneNumber string) (int64, error)
	CreateOTP(ctx context.Context, phoneNumber, otp string) error
	OTPInformation(ctx context.Context, phoneNumber string) (*time.Duration, *int64, error)
	GetOTP(ctx context.Context, phoneNumber string) (string, error)
}
