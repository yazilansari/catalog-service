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

	logger.Log.Info(
		"fetching cache",

		zap.String(
			"cache_key",
			key,
		),
	)

	start := time.Now()

	result, err :=
		Client.Get(
			ctx,
			key,
		).Result()

	duration := time.Since(start)

	// =========================
	// SLOW REDIS QUERY DETECTION
	// =========================

	if duration > time.Second {

		logger.Log.Warn(
			"slow redis operation detected",

			zap.Duration(
				"duration",
				duration,
			),

			zap.String(
				"operation",
				"Redis.Get",
			),

			zap.String(
				"redis_key",
				key,
			),
		)
	}

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

	logger.Log.Info(
		"cache hit",

		zap.String(
			"cache_key",
			key,
		),

		zap.Duration(
			"duration",
			duration,
		),
	)

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

	logger.Log.Info(
		"setting cache",

		zap.String(
			"cache_key",
			key,
		),
	)

	start := time.Now()

	data, err :=
		json.Marshal(
			value,
		)

	if err != nil {

		logger.Log.Error(
			"failed to marshal tenant cache data",

			zap.String(
				"cache_key",
				key,
			),

			zap.Error(err),
		)

		return err
	}

	err = Client.Set(
		ctx,
		key,
		data,
		ttl,
	).Err()

	duration := time.Since(start)

	// =========================
	// SLOW REDIS QUERY DETECTION
	// =========================

	if duration > time.Second {

		logger.Log.Warn(
			"slow redis operation detected",

			zap.Duration(
				"duration",
				duration,
			),

			zap.String(
				"operation",
				"Redis.Get",
			),

			zap.String(
				"redis_key",
				key,
			),
		)
	}

	if err != nil {

		logger.Log.Error(
			"failed to set cache",

			zap.String(
				"cache_key",
				key,
			),

			zap.Duration(
				"duration",
				duration,
			),

			zap.Error(err),
		)

		return err
	}

	logger.Log.Info(
		"cache set successfully",

		zap.String(
			"cache_key",
			key,
		),

		zap.Duration(
			"ttl",
			ttl,
		),

		zap.Duration(
			"duration",
			duration,
		),
	)

	return nil
}

func DeleteCache(
	ctx context.Context,
	keys ...string,
) error {

	logger.Log.Info("deleting cache")

	start := time.Now()

	err := Client.Del(
		ctx,
		keys...,
	).Err()

	duration := time.Since(start)

	// =========================
	// SLOW REDIS QUERY DETECTION
	// =========================

	if duration > time.Second {

		logger.Log.Warn(
			"slow redis operation detected",

			zap.Duration(
				"duration",
				duration,
			),

			zap.String(
				"operation",
				"Redis.Get",
			),

			zap.String(
				"redis_key",
				keys[0],
			),
		)
	}

	if err != nil {

		logger.Log.Error(
			"failed to delete cache",

			zap.Duration(
				"duration",
				duration,
			),

			zap.Error(err),
		)

		return err
	}

	logger.Log.Info(
		"cache delete success",

		zap.String(
			"cache_key",
			keys[0],
		),
	)

	return nil

}
