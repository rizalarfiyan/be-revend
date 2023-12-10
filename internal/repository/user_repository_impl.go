package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizalarfiyan/be-revend/internal/models"
	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/sql"
	"github.com/rizalarfiyan/be-revend/utils"
)

type userRepository struct {
	db           *pgxpool.Pool
	query        *sql.Queries
	queryBuilder *sql.Queries
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{
		db:           db,
		query:        sql.New(db),
		queryBuilder: sql.New(utils.QueryWrap(db)),
	}
}

func (r *userRepository) AllUser(ctx context.Context, req request.BasePagination) (*models.ContentPagination[sql.User], error) {
	var res models.ContentPagination[sql.User]

	baseBuilder := func(b *utils.QueryBuilder) {
		if req.Search != "" {
			b.Where("LOWER(CONCAT(first_name, ' ', last_name)) LIKE $1 OR LOWER(identity) LIKE $1 OR LOWER(phone_number) LIKE $1", fmt.Sprintf("%%%s%%", req.Search))
		}
	}

	users, err := r.queryBuilder.GetAllUsers(utils.QueryBuild(ctx, func(b *utils.QueryBuilder) {
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

	count, err := r.queryBuilder.CountAllUsers(utils.QueryBuild(ctx, baseBuilder))
	if err != nil {
		return nil, err
	}

	res.Content = users
	res.Count = count
	return &res, nil
}
