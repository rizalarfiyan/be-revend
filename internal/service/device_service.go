package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/response"
)

type DeviceService interface {
	GetAllDevice(ctx context.Context, req request.BasePagination) response.BaseResponsePagination[response.Device]
	GetAllDropdownDevice(ctx context.Context, req request.BasePagination) response.BaseResponsePagination[response.BaseDropdown]
	CreateDevice(ctx context.Context, req request.CreateDeviceRequest)
	UpdateDevice(ctx context.Context, req request.UpdateDeviceRequest)
	ToggleDeleteDevice(ctx context.Context, deviceId, userId uuid.UUID)
}
