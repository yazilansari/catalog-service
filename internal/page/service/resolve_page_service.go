package service

import (
	"catalog-service/internal/logger"
	"catalog-service/internal/page/dto"
	"catalog-service/internal/page/repository"
	redisClient "catalog-service/internal/redis"
	"errors"
	"time"

	"go.uber.org/zap"
)

func ResolvePage(
	tenantCode string,
	countryCode string,
	slug string,
) (*dto.ResolvePageResponse, error) {

	cacheKey :=
		"page:" +
			"resolve:" +
			tenantCode +
			":" +
			countryCode +
			":" +
			slug

	// =========================
	// CACHE
	// =========================

	logger.Log.Info(
		"get resolve page request",

		zap.String(
			"tenant_code",
			tenantCode,
		),

		zap.String(
			"country_code",
			countryCode,
		),

		zap.String(
			"type",
			"products",
		),

		zap.String(
			"slug",
			slug,
		),

		zap.String(
			"cache_key",
			cacheKey,
		),
	)

	// =========================
	// TRY REDIS
	// =========================

	logger.Log.Info(
		"checking redis cache",

		zap.String(
			"cache_key",
			cacheKey,
		),
	)

	// cached, err :=
	// 	redisClient.Client.Get(
	// 		redisClient.Ctx,
	// 		cacheKey,
	// 	).Result()

	// // CACHE HIT

	// if err == nil {

	// 	logger.Log.Info(
	// 		"resolve page cache hit",

	// 		zap.String(
	// 			"cache_key",
	// 			cacheKey,
	// 		),
	// 	)

	// 	var data dto.ResolvePageResponse

	// 	err = json.Unmarshal(
	// 		[]byte(cached),
	// 		&data,
	// 	)

	// 	if err == nil {
	// 		logger.Log.Info(
	// 			"resolve page cache unmarshal success",

	// 			zap.String(
	// 				"cache_key",
	// 				cacheKey,
	// 			),
	// 		)

	// 		return &data, nil
	// 	}

	// 	logger.Log.Error(
	// 		"resolve page cache unmarshal failed",

	// 		zap.String(
	// 			"cache_key",
	// 			cacheKey,
	// 		),

	// 		zap.Error(err),
	// 	)
	// }

	// // CACHE MISS

	// if err != nil {

	// 	logger.Log.Warn(
	// 		"redis cache miss",

	// 		zap.String(
	// 			"cache_key",
	// 			cacheKey,
	// 		),

	// 		zap.Error(err),
	// 	)
	// }

	redisStart := time.Now()

	cached, err :=
		redisClient.GetCache[dto.ResolvePageResponse](
			redisClient.Ctx,
			cacheKey,
		)

	redisDuration := time.Since(redisStart)

	// =========================
	// SLOW REDIS QUERY DETECTION
	// =========================

	if redisDuration > time.Second {

		logger.Log.Warn(
			"slow redis operation detected",

			zap.Duration(
				"duration",
				redisDuration,
			),

			zap.String(
				"operation",
				"Redis.Get",
			),

			zap.String(
				"redis_key",
				cacheKey,
			),
		)
	}

	if err == nil &&
		cached != nil {

		logger.Log.Info(
			"redis cache hit",

			zap.String(
				"cache_key",
				cacheKey,
			),
		)

		return cached, nil
	}

	// =========================
	// DATABASE FALLBACK
	// =========================

	logger.Log.Info(
		"fetching resolve page from database",

		zap.String(
			"tenant_code",
			tenantCode,
		),

		zap.String(
			"country_code",
			countryCode,
		),
	)

	start := time.Now()

	// =========================
	// CATEGORY
	// =========================

	data, err :=
		repository.ResolvePage(
			tenantCode,
			countryCode,
			slug,
		)

	duration := time.Since(start)

	if err != nil {

		logger.Log.Error(
			"failed to resolve page",

			zap.String("tenant_code", tenantCode),
			zap.String("country_code", countryCode),
			zap.String("slug", slug),
			zap.Duration("query_duration", duration),
			zap.Error(err),
		)

		return nil, err
	}

	if data == nil {

		logger.Log.Warn(
			"page not found",

			zap.String("slug", slug),
		)

		return nil, errors.New("page not found")
	}

	logger.Log.Info(
		"resolve page fetched successfully",

		zap.String(
			"page_type",
			data.PageType,
		),

		zap.Duration(
			"query_duration",
			duration,
		),
	)

	// =========================
	// SLOW QUERY
	// =========================

	if duration > time.Second {

		logger.Log.Warn(
			"slow resolve query detected",

			zap.String("slug", slug),
			zap.Duration("duration", duration),
			zap.String("country_code", countryCode),
			zap.String("tenant_code", tenantCode),
		)
	}

	// =========================
	// STORE CACHE
	// =========================

	// jsonData, _ := json.Marshal(data)

	// err = redisClient.Client.Set(
	// 	redisClient.Ctx,
	// 	cacheKey,
	// 	jsonData,
	// 	time.Hour,
	// ).Err()

	// if err != nil {

	// 	logger.Log.Error(
	// 		"failed to cache resolve page",

	// 		zap.String(
	// 			"cache_key",
	// 			cacheKey,
	// 		),

	// 		zap.Error(err),
	// 	)

	// } else {

	// 	logger.Log.Info(
	// 		"resolve page cached successfully",

	// 		zap.String(
	// 			"cache_key",
	// 			cacheKey,
	// 		),

	// 		zap.Duration(
	// 			"ttl",
	// 			time.Hour,
	// 		),
	// 	)
	// }

	redisSetStart := time.Now()

	err = redisClient.SetCache(
		redisClient.Ctx,
		cacheKey,
		data,
		redisClient.PageTTL,
	)

	redisSetDuration := time.Since(redisSetStart)

	// =========================
	// SLOW REDIS QUERY DETECTION
	// =========================

	if redisSetDuration > time.Second {

		logger.Log.Warn(
			"slow redis operation detected",

			zap.Duration(
				"duration",
				redisSetDuration,
			),

			zap.String(
				"operation",
				"Redis.Get",
			),

			zap.String(
				"redis_key",
				cacheKey,
			),
		)
	}

	if err != nil {

		logger.Log.Error(
			"failed to cache resolve page",

			zap.String(
				"cache_key",
				cacheKey,
			),

			zap.Error(err),
		)

	} else {

		logger.Log.Info(
			"resolve page cached successfully",

			zap.String(
				"cache_key",
				cacheKey,
			),

			zap.Duration(
				"ttl",
				redisClient.PageTTL,
			),
		)
	}

	return nil, err
}
