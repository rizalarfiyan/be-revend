package repository

import (
	"context"

	"github.com/rizalarfiyan/be-revend/internal/models"
)

type MqttRepository interface {
	CreateOrUpdateUserPoint(ctx context.Context, payload models.UserPoint) error
	GetUserPoint(ctx context.Context, identity string) (*models.UserPoint, error)
	DeleteUserPoint(ctx context.Context, identity string) error
}
