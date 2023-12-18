package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizalarfiyan/be-revend/constants"
	"github.com/rizalarfiyan/be-revend/internal/models"
	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/sql"
	"github.com/rizalarfiyan/be-revend/utils"
)

type historyRepository struct {
	db           *pgxpool.Pool
	query        *sql.Queries
	queryBuilder *sql.Queries
}

func NewHistoryRepository(db *pgxpool.Pool) HistoryRepository {
	return &historyRepository{
		db:           db,
		query:        sql.New(db),
		queryBuilder: sql.New(utils.QueryWrap(db)),
	}
}

func (r *historyRepository) AllHistory(ctx context.Context, req request.GetAllHistoryRequest) (*models.ContentPagination[sql.GetAllHistoryRow], error) {
	var res models.ContentPagination[sql.GetAllHistoryRow]

	baseBuilder := func(b *utils.QueryBuilder) {
		if req.Status != "" {
			switch req.Status {
			case constants.FilterListStatusDeleted:
				b.Where("u.deleted_at IS NOT NULL AND d.deleted_at IS NOT NULL")
			case constants.FilterListStatusActive:
				b.Where("u.deleted_at IS NULL AND d.deleted_at IS NULL")
			}
		}

		if req.DeviceId != uuid.Nil {
			b.Where("h.device_id = $1", req.DeviceId)
		}

		if req.UserId != uuid.Nil {
			b.Where("h.user_id = $1", req.UserId)
		}

		if req.Search != "" {
			b.Where("LOWER(d.name) LIKE $1 OR LOWER(CONCAT(u.first_name, ' ', u.last_name)) LIKE $1", fmt.Sprintf("%%%s%%", req.Search))
		}
	}

	histories, err := r.queryBuilder.GetAllHistory(utils.QueryBuild(ctx, func(b *utils.QueryBuilder) {
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

	count, err := r.queryBuilder.CountAllHistory(utils.QueryBuild(ctx, baseBuilder))
	if err != nil {
		return nil, err
	}

	res.Content = histories
	res.Count = count
	return &res, nil
}

func (r *historyRepository) CreateHistory(ctx context.Context, payload sql.CreateHistoryParams) error {
	return r.query.CreateHistory(ctx, payload)
}

func (r *historyRepository) AllHistoryStatistic(ctx context.Context, req models.AllHistoryStatistic) ([]sql.GetAllHistoryStatisticRow, error) {
	histories, err := r.queryBuilder.GetAllHistoryStatistic(utils.QueryBuild(ctx, func(b *utils.QueryBuilder) {
		b.Where("created_at BETWEEN $1 AND $2", req.StartDate, req.EndDate)
		b.Where("user_id = $1", req.UserId)

		switch req.TimeFrequency {
		case constants.FilterTimeFrequencyToday:
			b.GroupBy("DATE_TRUNC('hour', created_at)")
		case constants.FilterTimeFrequencyYear:
			b.GroupBy("DATE_TRUNC('month', created_at)")
		default:
			b.GroupBy("DATE_TRUNC('day', created_at)")
		}

		b.Order("date DESC")
	}))

	if err != nil {
		return nil, err
	}

	return histories, nil
}
