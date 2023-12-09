package repository

import (
	"context"

	"github.com/rizalarfiyan/be-revend/internal/models"
	"github.com/rizalarfiyan/be-revend/internal/sql"
)

type AuthRepository interface {
	GetUserByGoogleId(ctx context.Context, googleID string) (sql.User, error)
	GetUserByPhoneNumber(ctx context.Context, googleID string) (sql.User, error)
	GetUserByIdentity(ctx context.Context, identity string) (sql.User, error)
	GetUserByGoogleIdOrPhoneNumber(ctx context.Context, googleID, phoneNumber string) (sql.User, error)
	CreateUser(ctx context.Context, payload sql.CreateUserParams) error
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
