package repository

import (
	"context"

	"github.com/rizalarfiyan/be-revend/internal/models"
	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/sql"
)

type DeviceRepository interface {
	AllDevice(ctx context.Context, req request.BasePagination) (*models.ContentPagination[sql.Device], error)
	AllDropdownDevice(ctx context.Context, req request.BasePagination) (*models.ContentPagination[sql.GetAllNameDeviceRow], error)
	GetDeviceByToken(ctx context.Context, token string) (sql.Device, error)
	CreateDevice(ctx context.Context, payload sql.CreateDeviceParams) error
	UpdateDevice(ctx context.Context, payload sql.UpdateDeviceParams) error
	ToggleDeleteDevice(ctx context.Context, req sql.ToggleDeleteDeviceParams) error
}
