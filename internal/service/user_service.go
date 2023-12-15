package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/response"
)

type UserService interface {
	GetAllUser(ctx context.Context, req request.GetAllUserRequest) response.BaseResponsePagination[response.User]
	GetUserById(ctx context.Context, userId uuid.UUID) response.User
	GetAllDropdownUser(ctx context.Context, req request.BasePagination) response.BaseResponsePagination[response.BaseDropdown]
	ToggleDeleteUser(ctx context.Context, userId, currentUserId uuid.UUID)
}
