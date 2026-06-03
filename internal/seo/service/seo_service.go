package service

import (
	"time"

	"catalog-service/internal/logger"
	"catalog-service/internal/seo/dto"
	"catalog-service/internal/seo/repository"

	redisClient "catalog-service/internal/redis"

	"encoding/json"

	"go.uber.org/zap"
)

func GetSEOPage(
	tenantCode string,
	countryCode string,
	entityType string,
	slug string,
) (*dto.SEOResponse, error) {

	cacheKey :=
		"seo:" +
			entityType +
			":" +
			tenantCode +
			":" +
			countryCode +
			":" +
			slug

	// =========================
	// CACHE
	// =========================

	logger.Log.Info(
		"get seo request",

		zap.String(
			"tenant_code",
			tenantCode,
		),

		zap.String(
			"country_code",
			countryCode,
		),

		zap.String(
			"entity_type",
			entityType,
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
	// 		"seo cache hit",

	// 		zap.String(
	// 			"cache_key",
	// 			cacheKey,
	// 		),
	// 	)

	// 	var data dto.SEOResponse

	// 	err = json.Unmarshal(
	// 		[]byte(cached),
	// 		&data,
	// 	)

	// 	if err == nil {
	// 		logger.Log.Info(
	// 			"seo cache unmarshal success",

	// 			zap.String(
	// 				"cache_key",
	// 				cacheKey,
	// 			),
	// 		)

	// 		return &data, nil
	// 	}

	// 	logger.Log.Error(
	// 		"seo cache unmarshal failed",

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

	cached, err :=
		redisClient.GetCache[dto.SEOResponse](
			redisClient.Ctx,
			cacheKey,
		)

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
		"fetching seo from database",

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

	seo, err := repository.GetSEO(
		tenantCode,
		countryCode,
		entityType,
		slug,
	)

	duration := time.Since(start)

	if err != nil {
		logger.Log.Error(
			"failed to fetch seo from database",

			zap.String(
				"tenant_code",
				tenantCode,
			),

			zap.String(
				"country_code",
				countryCode,
			),

			zap.String(
				"entity_type",
				entityType,
			),

			zap.String(
				"slug",
				slug,
			),

			zap.Error(err),
		)
		return nil, err
	}

	logger.Log.Info(
		"seo fetched from database",

		zap.Duration(
			"query_duration",
			duration,
		),

		zap.String(
			"entity_type",
			entityType,
		),

		zap.String(
			"slug",
			slug,
		),
	)

	// =========================
	// SLOW QUERY DETECTION
	// =========================

	if duration > time.Second {

		logger.Log.Warn(
			"slow database query detected",

			zap.Duration(
				"duration",
				duration,
			),

			zap.String(
				"tenant_code",
				tenantCode,
			),

			zap.String(
				"country_code",
				countryCode,
			),

			zap.String(
				"entity_type",
				entityType,
			),

			zap.String(
				"slug",
				slug,
			),
		)
	}

	// =========================
	// RESPONSE
	// =========================

	data := dto.SEOResponse{
		Title: seo.Title,

		MetaDescription: seo.MetaDescription,

		MetaKeywords: seo.MetaKeywords,

		CanonicalURL: seo.CanonicalURL,

		OGTitle: seo.OGTitle,

		OGDescription: seo.OGDescription,

		OGImage: seo.OGImage,

		Robots: seo.Robots,

		SchemaJSON: seo.SchemaJSON,
	}

	// =========================
	// STORE CACHE
	// =========================

	jsonData, _ := json.Marshal(data)

	err = redisClient.Client.Set(
		redisClient.Ctx,
		cacheKey,
		jsonData,
		time.Hour,
	).Err()

	if err != nil {

		logger.Log.Error(
			"failed to store seo in redis",

			zap.String(
				"cache_key",
				cacheKey,
			),

			zap.Error(err),
		)

	} else {

		logger.Log.Info(
			"seo cached successfully",

			zap.String(
				"cache_key",
				cacheKey,
			),

			zap.Duration(
				"ttl",
				time.Hour,
			),
		)
	}

	return &data, nil
}
