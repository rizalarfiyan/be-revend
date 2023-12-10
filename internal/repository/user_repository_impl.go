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
	db    *pgxpool.Pool
	query *sql.Queries
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{
		db:    db,
		query: sql.New(db),
	}
}

func (r *userRepository) AllUser(ctx context.Context, req request.BasePagination) (*models.ContentPagination[sql.User], error) {
	var res models.ContentPagination[sql.User]
	// var count int

	query := sql.New(utils.QueryWrap(r.db))
	users, err := query.GetAllUsers(utils.QueryBuild(ctx, func(b *utils.QueryBuilder) {
		if req.Search != "" {
			b.Where("CONCAT(first_name, ' ', last_name) LIKE $1", fmt.Sprintf("'%%%s%%'", req.Search))
		}
	}))

	if err != nil {
		return nil, err
	}

	res.Content = users
	return &res, nil
}
