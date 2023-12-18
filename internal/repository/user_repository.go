package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-revend/internal/models"
	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/sql"
)

type UserRepository interface {
	AllUser(ctx context.Context, req request.GetAllUserRequest) (*models.ContentPagination[sql.User], error)
	GetUserById(ctx context.Context, userId uuid.UUID) (sql.User, error)
	GetUserByGoogleId(ctx context.Context, googleID string) (sql.User, error)
	GetUserByPhoneNumber(ctx context.Context, googleID string) (sql.User, error)
	GetUserByIdentity(ctx context.Context, identity string) (sql.User, error)
	GetUserByGoogleIdOrPhoneNumber(ctx context.Context, googleID, phoneNumber string) (sql.User, error)
	CreateUser(ctx context.Context, payload sql.CreateUserParams) error
	UpdateUser(ctx context.Context, payload sql.UpdateUserParams) error
	AllDropdownUsers(ctx context.Context, req request.BasePagination) (*models.ContentPagination[sql.GetAllNameUsersRow], error)
	ToggleDeleteUser(ctx context.Context, req sql.ToggleDeleteUserParams) error
	UpdateUserProfile(ctx context.Context, payload sql.UpdateUserProfileParams) error
	DeleteGoogleUserProfile(ctx context.Context, userId uuid.UUID) error
	UpdateGoogleUserProfile(ctx context.Context, payload sql.UpdateGoogleUserProfileParams) error
	CreateBindGoogle(ctx context.Context, token string, userId uuid.UUID) error
	CreateBindGoogleFresh(ctx context.Context, token string, userId uuid.UUID) error
	GetBindGoogle(ctx context.Context, token string) (string, error)
	DeleteBindGoogle(ctx context.Context, userId uuid.UUID) error
}
