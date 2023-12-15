package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
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

func (r *userRepository) AllUser(ctx context.Context, req request.GetAllUserRequest) (*models.ContentPagination[sql.User], error) {
	var res models.ContentPagination[sql.User]

	baseBuilder := func(b *utils.QueryBuilder) {
		if req.Search != "" {
			b.Where("LOWER(CONCAT(first_name, ' ', last_name)) LIKE $1 OR LOWER(identity) LIKE $1 OR LOWER(phone_number) LIKE $1", fmt.Sprintf("%%%s%%", req.Search))
		}

		if req.Role != "" {
			b.Where("role = $1", req.Role)
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

func (r *userRepository) GetUserById(ctx context.Context, userId uuid.UUID) (sql.User, error) {
	return r.query.GetUserById(ctx, utils.PGUUID(userId))
}

func (r *userRepository) GetUserByGoogleId(ctx context.Context, googleID string) (sql.User, error) {
	return r.query.GetUserByGoogleId(ctx, utils.PGText(googleID))
}

func (r *userRepository) GetUserByPhoneNumber(ctx context.Context, googleID string) (sql.User, error) {
	return r.query.GetUserByPhoneNumber(ctx, googleID)
}

func (r *userRepository) GetUserByIdentity(ctx context.Context, identity string) (sql.User, error) {
	return r.query.GetUserByIdentity(ctx, identity)
}

func (r *userRepository) GetUserByGoogleIdOrPhoneNumber(ctx context.Context, googleID, phoneNumber string) (sql.User, error) {
	return r.query.GetUserByGoogleIdOrPhoneNumber(ctx, sql.GetUserByGoogleIdOrPhoneNumberParams{
		GoogleID:    utils.PGText(googleID),
		PhoneNumber: phoneNumber,
	})
}

func (r *userRepository) CreateUser(ctx context.Context, payload sql.CreateUserParams) error {
	return r.query.CreateUser(ctx, payload)
}

func (r *userRepository) AllDropdownUsers(ctx context.Context, req request.BasePagination) (*models.ContentPagination[sql.GetAllNameUsersRow], error) {
	var res models.ContentPagination[sql.GetAllNameUsersRow]

	baseBuilder := func(b *utils.QueryBuilder) {
		if req.Search != "" {
			b.Where("LOWER(CONCAT(first_name, ' ', last_name)) LIKE $1", fmt.Sprintf("%%%s%%", req.Search))
		}
	}

	users, err := r.queryBuilder.GetAllNameUsers(utils.QueryBuild(ctx, func(b *utils.QueryBuilder) {
		baseBuilder(b)
		b.Order("created_at DESC")
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

func (r *userRepository) ToggleDeleteUser(ctx context.Context, req sql.ToggleDeleteUserParams) error {
	return r.query.ToggleDeleteUser(ctx, req)
}
