package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rizalarfiyan/be-revend/config"
	"github.com/rizalarfiyan/be-revend/constants"
	"github.com/rizalarfiyan/be-revend/database"
	"github.com/rizalarfiyan/be-revend/internal/models"
	"github.com/rizalarfiyan/be-revend/internal/sql"
	baseModels "github.com/rizalarfiyan/be-revend/models"
	"github.com/rizalarfiyan/be-revend/utils"
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
	return r.redis.Setxc(key, r.conf.Auth.Social.SessionDuration, string(strPayload))
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

func (r *repository) IncrementOTP(ctx context.Context, phoneNumber string) (int64, error) {
	key := fmt.Sprintf(constants.KeyOTPIncrement, phoneNumber)
	inc, err := r.redis.Increment(key)
	if err != nil {
		return 0, err
	}

	if inc != 1 {
		return inc, nil
	}

	return inc, r.redis.SetDuration(key, utils.RemaniningToday())
}

func (r *repository) CreateOTP(ctx context.Context, phoneNumber, otp string) error {
	key := fmt.Sprintf(constants.KeyOTP, phoneNumber)
	return r.redis.Setxc(key, r.conf.Auth.OTP.Duration, otp)
}

func (r *repository) OTPInformation(ctx context.Context, phoneNumber string) (*time.Duration, *int64, error) {
	key := fmt.Sprintf(constants.KeyOTP, phoneNumber)
	keyInc := fmt.Sprintf(constants.KeyOTPIncrement, phoneNumber)
	duration, err := r.redis.Duration(key)
	if err != nil {
		return nil, nil, err
	}

	incStr, err := r.redis.GetString(keyInc)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, nil, err
	}

	var inc int64
	if incStr != "" {
		inc, err = strconv.ParseInt(incStr, 10, 64)
		if err != nil {
			return nil, nil, err
		}
	}

	return duration, &inc, nil
}
