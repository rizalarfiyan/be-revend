package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rizalarfiyan/be-revend/constants"
	"github.com/rizalarfiyan/be-revend/database"
	"github.com/rizalarfiyan/be-revend/internal/models"
)

type mqttRepository struct {
	redis database.RedisInstance
}

func NewMqttRepository(redis database.RedisInstance) MqttRepository {
	return &mqttRepository{
		redis: redis,
	}
}

func (r *mqttRepository) CreateOrUpdateUserPoint(ctx context.Context, payload models.UserPoint) error {
	key := fmt.Sprintf(constants.KeyUserPoint, payload.Identity)
	strPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return r.redis.Set(key, string(strPayload))
}

func (r *mqttRepository) GetUserPoint(ctx context.Context, identity string) (*models.UserPoint, error) {
	var res models.UserPoint
	key := fmt.Sprintf(constants.KeyUserPoint, identity)
	err := r.redis.Get(key, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}


func (r *mqttRepository) DeleteUserPoint(ctx context.Context, identity string) error {
    key := fmt.Sprintf(constants.KeyUserPoint, identity)
    return r.redis.Del(key)
}