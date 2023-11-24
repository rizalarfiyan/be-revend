package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizalarfiyan/be-revend/config"
	"github.com/rizalarfiyan/be-revend/constants"
	"github.com/rizalarfiyan/be-revend/database"
	"github.com/rizalarfiyan/be-revend/internal/models"
	"github.com/rizalarfiyan/be-revend/internal/sql"
	baseModels "github.com/rizalarfiyan/be-revend/models"
)

type repository struct {
	db    *pgxpool.Pool
	query *sql.Queries
	redis database.RedisInstance
	conf  *baseModels.Config
}

func NewAuthRepository(db *pgxpool.Pool, redis database.RedisInstance) Repository {
	return &repository{
		db:    db,
		query: sql.New(db),
		redis: redis,
		conf:  config.Get(),
	}
}

func (r *repository) GetUserByGoogleId(ctx context.Context, googleID string) (sql.User, error) {
	return r.query.GetUserByGoogleId(ctx, googleID)
}

func (r *repository) GetUserByPhoneNumber(ctx context.Context, googleID string) (sql.User, error) {
	return r.query.GetUserByPhoneNumber(ctx, googleID)
}

func (r *repository) CreateSocialSession(ctx context.Context, idx string, payload models.SocialSession) error {
	keySearch := fmt.Sprintf(constants.KeySocialSession, "*", payload.GoogleId)
	err := r.redis.DelKeysByPatern(keySearch)
	if err != nil {
		return err
	}

	strPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	key := fmt.Sprintf(constants.KeySocialSession, idx, payload.GoogleId)
	return r.redis.Setxc(key, r.conf.Auth.SocialSessionDuration, string(strPayload))
}
