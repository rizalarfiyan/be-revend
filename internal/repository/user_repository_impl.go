package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizalarfiyan/be-revend/config"
	"github.com/rizalarfiyan/be-revend/constants"
	"github.com/rizalarfiyan/be-revend/database"
	"github.com/rizalarfiyan/be-revend/internal/models"
	"github.com/rizalarfiyan/be-revend/internal/request"
	"github.com/rizalarfiyan/be-revend/internal/sql"
	baseModels "github.com/rizalarfiyan/be-revend/models"
	"github.com/rizalarfiyan/be-revend/utils"
)

type userRepository struct {
	db           *pgxpool.Pool
	query        *sql.Queries
	queryBuilder *sql.Queries
	redis        database.RedisInstance
	conf         *baseModels.Config
}

func NewUserRepository(db *pgxpool.Pool, redis database.RedisInstance) UserRepository {
	return &userRepository{
		db:           db,
		query:        sql.New(db),
		queryBuilder: sql.New(utils.QueryWrap(db)),
		redis:        redis,
		conf:         config.Get(),
	}
}

func (r *userRepository) AllUser(ctx context.Context, req request.GetAllUserRequest) (*models.ContentPagination[sql.User], error) {
	var res models.ContentPagination[sql.User]

	baseBuilder := func(b *utils.QueryBuilder) {
		if req.Status != "" {
			switch req.Status {
			case constants.FilterListStatusDeleted:
				b.Where("deleted_at IS NOT NULL")
			case constants.FilterListStatusActive:
				b.Where("deleted_at IS NULL")
			}
		}

		if req.Role != "" {
			b.Where("role = $1", req.Role)
		}

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

func (r *userRepository) UpdateUser(ctx context.Context, payload sql.UpdateUserParams) error {
	return r.query.UpdateUser(ctx, payload)
}

func (r *userRepository) UpdateUserProfile(ctx context.Context, payload sql.UpdateUserProfileParams) error {
	return r.query.UpdateUserProfile(ctx, payload)
}

func (r *userRepository) AllDropdownUsers(ctx context.Context, req request.BasePagination) (*models.ContentPagination[sql.GetAllNameUsersRow], error) {
	var res models.ContentPagination[sql.GetAllNameUsersRow]

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

func (r *userRepository) DeleteGoogleUserProfile(ctx context.Context, userId uuid.UUID) error {
	return r.query.DeleteGoogleUserProfile(ctx, utils.PGUUID(userId))
}

func (r *userRepository) UpdateGoogleUserProfile(ctx context.Context, payload sql.UpdateGoogleUserProfileParams) error {
	return r.query.UpdateGoogleUserProfile(ctx, payload)
}

func (r *userRepository) CreateBindGoogle(ctx context.Context, token string, userId uuid.UUID) error {
	key := fmt.Sprintf(constants.KeyBindGoogle, token, userId)
	return r.redis.Setxc(key, r.conf.Auth.OTP.Duration, userId.String())
}

func (r *userRepository) CreateBindGoogleFresh(ctx context.Context, token string, userId uuid.UUID) error {
	err := r.DeleteBindGoogle(ctx, userId)
	if err != nil {
		return err
	}

	return r.CreateBindGoogle(ctx, token, userId)
}

func (r *userRepository) GetBindGoogle(ctx context.Context, token string) (string, error) {
	keySearch := fmt.Sprintf(constants.KeyBindGoogle, token, "*")

	key, err := r.redis.Keys(keySearch)
	if err != nil {
		return "", err
	}

	if len(key) <= 0 {
		return "", nil
	}

	return r.redis.GetString(key[0])
}

func (r *userRepository) DeleteBindGoogle(ctx context.Context, userId uuid.UUID) error {
	keySearch := fmt.Sprintf(constants.KeyBindGoogle, "*", userId.String())
	return r.redis.DelKeysByPatern(keySearch)
}
