package repository

import (
	"context"

	"github.com/rizalarfiyan/be-revend/internal/models"
	"github.com/rizalarfiyan/be-revend/internal/sql"
)

type Repository interface {
	GetUserByGoogleId(ctx context.Context, googleID string) (sql.User, error)
	GetUserByPhoneNumber(ctx context.Context, googleID string) (sql.User, error)
	CreateSocialSession(ctx context.Context, idx string, payload models.SocialSession) error
	GetSocialSessionByToken(ctx context.Context, token string) (*models.SocialSession, error)
	DeleteSocialSessionByGoogleId(ctx context.Context, googleId string) error
	DeleteSocialSessionByToken(ctx context.Context, token string) error
}
