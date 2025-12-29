package redisclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/s-turchinskiy/keeper/internal/server/models"
	"github.com/s-turchinskiy/keeper/internal/utils/errorsutils"
	"time"
)

type RedisClient struct {
	rdb        *redis.Client
	expiration time.Duration
}

func NewRedisClient(rdb *redis.Client, expiration time.Duration) *RedisClient {
	return &RedisClient{
		rdb:        rdb,
		expiration: expiration,
	}

}
func (r *RedisClient) Set(ctx context.Context, secret *models.Secret) error {

	if r.rdb == nil || secret == nil {
		return nil
	}

	key := r.key(secret.UserID, secret.ID)
	bytes, err := json.Marshal(*secret)
	if err != nil {
		fmt.Println(errorsutils.WrapError(err))
		return err
	}

	err = r.rdb.Set(ctx, key, bytes, r.expiration).Err()
	if err != nil {
		fmt.Println(errorsutils.WrapError(err))
		return err
	}

	return nil

}

func (r *RedisClient) Get(ctx context.Context, userID, secretID string) (secret *models.Secret, err error) {

	if r.rdb == nil {
		return nil, nil
	}

	key := r.key(userID, secretID)
	val, err := r.rdb.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		fmt.Println(errorsutils.WrapError(err))
		return nil, err
	}

	if val == "" {
		return nil, nil
	}

	err = json.Unmarshal([]byte(val), &secret)
	if err != nil {
		fmt.Println(errorsutils.WrapError(err))
		return nil, err
	}
	return secret, nil
}

func (r *RedisClient) key(userID, secretID string) string {
	return fmt.Sprintf("secret_%s_%s", userID, secretID)
}
