package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizalarfiyan/be-revend/constants"
	"github.com/rizalarfiyan/be-revend/database"
	"github.com/rizalarfiyan/be-revend/internal/models"
	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/sql"
	"github.com/rizalarfiyan/be-revend/utils"
)

type deviceRepository struct {
	db           *pgxpool.Pool
	query        *sql.Queries
	queryBuilder *sql.Queries
	redis        database.RedisInstance
}

func NewDeviceRepository(db *pgxpool.Pool, redis database.RedisInstance) DeviceRepository {
	return &deviceRepository{
		db:           db,
		query:        sql.New(db),
		queryBuilder: sql.New(utils.QueryWrap(db)),
		redis:        redis,
	}
}

func (r *deviceRepository) AllDevice(ctx context.Context, req request.BasePagination) (*models.ContentPagination[sql.Device], error) {
	var res models.ContentPagination[sql.Device]

	baseBuilder := func(b *utils.QueryBuilder) {
		if req.Status != "" {
			switch req.Status {
			case constants.FilterListStatusDeleted:
				b.Where("deleted_at IS NOT NULL")
			case constants.FilterListStatusActive:
				b.Where("deleted_at IS NULL")
			}
		}

		if req.Search != "" {
			b.Where("LOWER(name) LIKE $1 OR LOWER(location) LIKE $1", fmt.Sprintf("%%%s%%", req.Search))
		}
	}

	devices, err := r.queryBuilder.GetAllDevice(utils.QueryBuild(ctx, func(b *utils.QueryBuilder) {
		baseBuilder(b)
		if req.OrderBy != "" && req.Order != "" {
			b.Ordering(req.OrderBy, req.Order)
		} else {
			b.Order("created_at DESC")
		}
		b.Pagination(req.Page, req.Limit)
	}))

	if err != nil {
		return nil, err
	}

	count, err := r.queryBuilder.CountAllDevice(utils.QueryBuild(ctx, baseBuilder))
	if err != nil {
		return nil, err
	}

	res.Content = devices
	res.Count = count
	return &res, nil
}

func (r *deviceRepository) AllDropdownDevice(ctx context.Context, req request.BasePagination) (*models.ContentPagination[sql.GetAllNameDeviceRow], error) {
	var res models.ContentPagination[sql.GetAllNameDeviceRow]

	baseBuilder := func(b *utils.QueryBuilder) {
		if req.Status != "" {
			switch req.Status {
			case constants.FilterListStatusDeleted:
				b.Where("deleted_at IS NOT NULL")
			case constants.FilterListStatusActive:
				b.Where("deleted_at IS NULL")
			}
		}

		if req.Search != "" {
			b.Where("LOWER(name) LIKE $1 OR LOWER(location) LIKE $1", fmt.Sprintf("%%%s%%", req.Search))
		}
	}

	devices, err := r.queryBuilder.GetAllNameDevice(utils.QueryBuild(ctx, func(b *utils.QueryBuilder) {
		baseBuilder(b)
		b.Order("created_at DESC")
		b.Pagination(req.Page, req.Limit)
	}))

	if err != nil {
		return nil, err
	}

	count, err := r.queryBuilder.CountAllDevice(utils.QueryBuild(ctx, baseBuilder))
	if err != nil {
		return nil, err
	}

	res.Content = devices
	res.Count = count
	return &res, nil
}

func (r *deviceRepository) GetDeviceByToken(ctx context.Context, token string) (sql.Device, error) {
	return r.query.GetDeviceByToken(ctx, token)
}

func (r *deviceRepository) CreateDevice(ctx context.Context, payload sql.CreateDeviceParams) error {
	return r.query.CreateDevice(ctx, payload)
}

func (r *deviceRepository) UpdateDevice(ctx context.Context, payload sql.UpdateDeviceParams) error {
	return r.query.UpdateDevice(ctx, payload)
}

func (r *deviceRepository) ToggleDeleteDevice(ctx context.Context, req sql.ToggleDeleteDeviceParams) error {
	return r.query.ToggleDeleteDevice(ctx, req)
}
