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

func (r *repository) CreateSocialSession(ctx context.Context, token string, payload models.SocialSession) error {
	err := r.DeleteSocialSessionByGoogleId(ctx, payload.GoogleId)
	if err != nil {
		return err
	}

	strPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	key := fmt.Sprintf(constants.KeySocialSession, token, payload.GoogleId)
	return r.redis.Setxc(key, r.conf.Auth.SocialSessionDuration, string(strPayload))
}

func (r *repository) GetSocialSessionByToken(ctx context.Context, token string) (*models.SocialSession, error) {
	keySearch := fmt.Sprintf(constants.KeySocialSession, token, "*")
	key, err := r.redis.Keys(keySearch)
	if err != nil {
		return nil, err
	}

	if len(key) <= 0 {
		return nil, nil
	}

	var res models.SocialSession
	err = r.redis.Get(key[0], &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *repository) DeleteSocialSessionByGoogleId(ctx context.Context, googleId string) error {
	keySearch := fmt.Sprintf(constants.KeySocialSession, "*", googleId)
	return r.redis.DelKeysByPatern(keySearch)
}

func (r *repository) DeleteSocialSessionByToken(ctx context.Context, token string) error {
	keySearch := fmt.Sprintf(constants.KeySocialSession, token, "*")
	return r.redis.DelKeysByPatern(keySearch)
}
