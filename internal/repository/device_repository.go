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
}
