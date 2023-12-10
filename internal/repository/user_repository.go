package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-revend/internal/models"
	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/sql"
)

type UserRepository interface {
	AllUser(ctx context.Context, req request.BasePagination) (*models.ContentPagination[sql.User], error)
	GetUserById(ctx context.Context, userId uuid.UUID) (sql.User, error)
}
