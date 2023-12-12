package service

import (
	"context"

	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/response"
)

type DeviceService interface {
	GetAllDevice(ctx context.Context, req request.BasePagination) response.BaseResponsePagination[response.Device]
}
