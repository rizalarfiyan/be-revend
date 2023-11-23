package service

import (
	"context"

	"github.com/rizalarfiyan/be-revend/internal/request"
)

type AuthService interface {
	Google() string
	GoogleCallback(ctx context.Context, req request.GoogleCallbackRequest) string
}
