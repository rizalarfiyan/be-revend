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

type authRepository struct {
	db    *pgxpool.Pool
	query *sql.Queries
	redis database.RedisInstance
	conf  *baseModels.Config
}

func NewAuthRepository(db *pgxpool.Pool, redis database.RedisInstance) AuthRepository {
	return &authRepository{
		db:    db,
		query: sql.New(db),
		redis: redis,
		conf:  config.Get(),
	}
}

func (r *authRepository) GetUserByGoogleId(ctx context.Context, googleID string) (sql.User, error) {
	return r.query.GetUserByGoogleId(ctx, googleID)
}

func (r *authRepository) GetUserByPhoneNumber(ctx context.Context, googleID string) (sql.User, error) {
	return r.query.GetUserByPhoneNumber(ctx, googleID)
}

func (r *authRepository) GetUserByIdentity(ctx context.Context, identity string) (sql.User, error) {
	return r.query.GetUserByIdentity(ctx, identity)
}

func (r *authRepository) GetUserByGoogleIdOrPhoneNumber(ctx context.Context, googleID, phoneNumber string) (sql.User, error) {
	return r.query.GetUserByGoogleIdOrPhoneNumber(ctx, sql.GetUserByGoogleIdOrPhoneNumberParams{
		GoogleID:    googleID,
		PhoneNumber: phoneNumber,
	})
}

func (r *authRepository) CreateUser(ctx context.Context, payload sql.CreateUserParams) error {
	return r.query.CreateUser(ctx, payload)
}

func (r *authRepository) CreateVerificationSession(ctx context.Context, token string, payload models.VerificationSession) error {
	if !utils.IsEmpty(payload.GoogleId) {
		err := r.DeleteVerificationSessionByGoogleId(ctx, payload.GoogleId)
		if err != nil {
			return err
		}
	}

	if !utils.IsEmpty(payload.Identity) {
		err := r.DeleteVerificationSessionByIdentity(ctx, payload.Identity)
		if err != nil {
			return err
		}
	}

	strPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	duration := r.conf.Auth.Verification.Duration
	if payload.IsError {
		duration = 1 * time.Minute
	}

	key := fmt.Sprintf(constants.KeyVerificationSession, token, payload.GoogleId, payload.PhoneNumber, payload.Identity)
	return r.redis.Setxc(key, duration, string(strPayload))
}

func (r *authRepository) getVerificationSession(ctx context.Context, keySearch string) (*models.VerificationSession, error) {
	key, err := r.redis.Keys(keySearch)
	if err != nil {
		return nil, err
	}

	if len(key) <= 0 {
		return nil, nil
	}

	var res models.VerificationSession
	err = r.redis.Get(key[0], &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *authRepository) GetVerificationSessionByToken(ctx context.Context, token string) (*models.VerificationSession, error) {
	keySearch := fmt.Sprintf(constants.KeyVerificationSession, token, "*", "*", "*")
	return r.getVerificationSession(ctx, keySearch)
}

func (r *authRepository) GetVerificationSessionByPhoneNumber(ctx context.Context, phoneNumber string) (*models.VerificationSession, error) {
	keySearch := fmt.Sprintf(constants.KeyVerificationSession, "*", "*", phoneNumber, "*")
	return r.getVerificationSession(ctx, keySearch)
}

func (r *authRepository) GetVerificationSessionByIdentity(ctx context.Context, identity string) (*models.VerificationSession, error) {
	keySearch := fmt.Sprintf(constants.KeyVerificationSession, "*", "*", "*", identity)
	return r.getVerificationSession(ctx, keySearch)
}

func (r *authRepository) DeleteVerificationSessionByGoogleId(ctx context.Context, googleId string) error {
	keySearch := fmt.Sprintf(constants.KeyVerificationSession, "*", googleId, "*", "*")
	return r.redis.DelKeysByPatern(keySearch)
}

func (r *authRepository) DeleteVerificationSessionByIdentity(ctx context.Context, identity string) error {
	keySearch := fmt.Sprintf(constants.KeyVerificationSession, "*", "*", "*", identity)
	return r.redis.DelKeysByPatern(keySearch)
}

func (r *authRepository) DeleteVerificationSessionByToken(ctx context.Context, token string) error {
	keySearch := fmt.Sprintf(constants.KeyVerificationSession, token, "*", "*", "*")
	return r.redis.DelKeysByPatern(keySearch)
}

func (r *authRepository) IncrementOTP(ctx context.Context, phoneNumber string) (int64, error) {
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

func (r *authRepository) CreateOTP(ctx context.Context, phoneNumber, otp string) error {
	key := fmt.Sprintf(constants.KeyOTP, phoneNumber)
	return r.redis.Setxc(key, r.conf.Auth.OTP.Duration, otp)
}

func (r *authRepository) OTPInformation(ctx context.Context, phoneNumber string) (*models.OTPInformation, error) {
	res := models.OTPInformation{}
	key := fmt.Sprintf(constants.KeyOTP, phoneNumber)
	keyInc := fmt.Sprintf(constants.KeyOTPIncrement, phoneNumber)
	duration, err := r.redis.Duration(key)
	if err != nil {
		return nil, err
	}

	incStr, err := r.redis.GetString(keyInc)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	var inc int64
	if incStr != "" {
		inc, err = strconv.ParseInt(incStr, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	otp, err := r.GetVerificationSessionByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return nil, err
	}

	res.Increment = inc

	if otp != nil {
		res.Data = *otp
	}

	if duration != nil {
		res.Duration = *duration
	}

	return &res, nil
}

func (r *authRepository) GetOTP(ctx context.Context, phoneNumber string) (string, error) {
	key := fmt.Sprintf(constants.KeyOTP, phoneNumber)
	return r.redis.GetString(key)
}

func (r *authRepository) DeleteAllOTP(ctx context.Context, phoneNumber string) error {
	keySearch := fmt.Sprintf(constants.KeyOTP+"*", phoneNumber)
	return r.redis.DelKeysByPatern(keySearch)
}
