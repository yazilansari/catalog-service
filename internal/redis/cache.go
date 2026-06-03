package redis

import (
	"catalog-service/internal/logger"
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func GetCache[T any](
	ctx context.Context,
	key string,
) (*T, error) {

	result, err :=
		Client.Get(
			ctx,
			key,
		).Result()

	if err != nil {

		logger.Log.Warn(
			"redis cache miss",

			zap.String(
				"cache_key",
				key,
			),

			zap.Error(err),
		)

		if err == redis.Nil {
			return nil, nil
		}

		return nil, err
	}

	var data T

	err =
		json.Unmarshal(
			[]byte(result),
			&data,
		)

	if err != nil {

		logger.Log.Error(
			"redis unmarshal failed",

			zap.String(
				"cache_key",
				key,
			),

			zap.Error(err),
		)

		return nil, err
	}

	logger.Log.Info(
		"redis unmarshal success",

		zap.String(
			"cache_key",
			key,
		),
	)

	return &data, nil
}

func SetCache(
	ctx context.Context,
	key string,
	value interface{},
	ttl time.Duration,
) error {

	data, err :=
		json.Marshal(
			value,
		)

	if err != nil {
		return err
	}

	return Client.Set(
		ctx,
		key,
		data,
		ttl,
	).Err()
}

func DeleteCache(
	ctx context.Context,
	keys ...string,
) error {

	return Client.Del(
		ctx,
		keys...,
	).Err()
}
