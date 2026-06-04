package service

import (
	"catalog-service/internal/logger"
	"catalog-service/internal/plp/dto"
	"catalog-service/internal/plp/model"
	"catalog-service/internal/plp/repository"
	redisClient "catalog-service/internal/redis"
	"encoding/json"
	"strconv"
	"time"

	"go.uber.org/zap"
)

func GetProducts(
	query dto.ProductQuery,
	tenantCode string,
	countryCode string,
) (*dto.ProductListResponse, error) {

	// =========================
	// DEFAULTS
	// =========================

	if query.Limit <= 0 {
		query.Limit = 20
	}

	if query.Sort == "" {
		query.Sort = "latest"
	}

	// =========================
	// CACHE KEY
	// =========================

	cacheKey := buildCacheKey(
		query,
		tenantCode,
		countryCode,
	)

	logger.Log.Info(
		"get products request",

		zap.String(
			"cache_key",
			cacheKey,
		),
	)

	// =========================
	// CACHE
	// =========================

	// cached, err :=
	// 	redisClient.Client.Get(
	// 		redisClient.Ctx,
	// 		cacheKey,
	// 	).Result()

	// // CACHE HIT

	// if err == nil {

	// 	logger.Log.Info(
	// 		"plp cache hit",

	// 		zap.String(
	// 			"cache_key",
	// 			cacheKey,
	// 		),
	// 	)

	// 	var response dto.ProductListResponse

	// 	err =
	// 		json.Unmarshal(
	// 			[]byte(cached),
	// 			&response,
	// 		)

	// 	if err == nil {

	// 		logger.Log.Info(
	// 			"plp cache unmarshal success",

	// 			zap.String(
	// 				"cache_key",
	// 				cacheKey,
	// 			),
	// 		)

	// 		return &response, nil
	// 	}

	// 	logger.Log.Error(
	// 		"plp page cache unmarshal failed",

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
	// 		"plp redis cache miss",

	// 		zap.String(
	// 			"cache_key",
	// 			cacheKey,
	// 		),

	// 		zap.Error(err),
	// 	)
	// }

	redisStart := time.Now()

	cached, err :=
		redisClient.GetCache[dto.ProductListResponse](
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
	// ELASTICSEARCH
	// =========================

	start :=
		time.Now()

	rows,
		aggs,
		total,
		err :=
		repository.SearchProducts(
			query,
			tenantCode,
			countryCode,
		)

	logger.Log.Info(
		"search result",

		zap.Int(
			"rows",
			len(rows),
		),

		zap.Bool(
			"aggs_nil",
			aggs == nil,
		),

		zap.Int64(
			"total",
			total,
		),
	)

	duration :=
		time.Since(start)

	if err != nil {

		logger.Log.Error(
			"failed to search products",

			zap.Error(err),
		)

		return nil, err
	}

	logger.Log.Info(
		"products fetched from elasticsearch",

		zap.Int64(
			"total",
			total,
		),

		zap.Duration(
			"duration",
			duration,
		),
	)

	if duration > time.Second {

		logger.Log.Warn(
			"slow plp query detected",

			zap.Any(
				"query",
				query,
			),

			zap.Duration(
				"duration",
				duration,
			),

			zap.String(
				"country_code",
				countryCode,
			),

			zap.String(
				"tenant_code",
				tenantCode,
			),
		)
	}

	// =========================
	// MAP PRODUCTS
	// =========================

	var products []model.ProductDocument

	for _, row := range rows {

		data, _ :=
			json.Marshal(
				row,
			)

		var product model.ProductDocument

		_ = json.Unmarshal(
			data,
			&product,
		)

		products =
			append(
				products,
				product,
			)
	}

	// =========================
	// PAGINATION
	// =========================

	hasMore := false

	nextCursor := ""

	if len(products) > 0 {

		hasMore =
			len(products) >= query.Limit

		if hasMore {

			last :=
				products[len(products)-1]

			nextCursor =
				strconv.FormatUint(
					last.ID,
					10,
				)
		}
	}

	// =========================
	// FILTERS
	// =========================

	var brands []string

	var minPrice float64
	var maxPrice float64

	if aggs != nil {

		// =========================
		// BRANDS
		// =========================

		if brandsAgg, ok :=
			aggs["brands"].(map[string]interface{}); ok {

			if buckets, ok :=
				brandsAgg["buckets"].([]interface{}); ok {

				for _, bucket := range buckets {

					bucketMap, ok :=
						bucket.(map[string]interface{})

					if !ok {
						continue
					}

					key, ok :=
						bucketMap["key"].(string)

					if !ok {
						continue
					}

					brands =
						append(
							brands,
							key,
						)
				}
			}
		}

		// =========================
		// MIN PRICE
		// =========================

		if minAgg, ok :=
			aggs["min_price"].(map[string]interface{}); ok {

			if value, ok :=
				minAgg["value"].(float64); ok {

				minPrice =
					value
			}
		}

		// =========================
		// MAX PRICE
		// =========================

		if maxAgg, ok :=
			aggs["max_price"].(map[string]interface{}); ok {

			if value, ok :=
				maxAgg["value"].(float64); ok {

				maxPrice =
					value
			}
		}
	}

	filters :=
		dto.FilterResponse{
			Brands: brands,

			PriceRange: dto.PriceRange{
				Min: minPrice,
				Max: maxPrice,
			},
		}

	// =========================
	// RESPONSE
	// =========================

	response :=
		dto.ProductListResponse{
			Products: products,

			Filters: filters,

			Pagination: dto.PaginationResponse{
				NextCursor: nextCursor,

				Limit: query.Limit,

				HasMore: hasMore,
			},

			Sort: query.Sort,

			Total: total,
		}

	// =========================
	// STORE CACHE
	// =========================

	// jsonData, _ :=
	// 	json.Marshal(
	// 		response,
	// 	)

	// err = redisClient.Client.Set(
	// 	redisClient.Ctx,
	// 	cacheKey,
	// 	jsonData,
	// 	15*time.Minute,
	// ).Err()

	// if err != nil {

	// 	logger.Log.Error(
	// 		"failed to store plp in redis",

	// 		zap.String(
	// 			"cache_key",
	// 			cacheKey,
	// 		),

	// 		zap.Error(err),
	// 	)

	// } else {

	// 	logger.Log.Info(
	// 		"plp cached successfully",

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
		response,
		redisClient.PLPTTL,
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
			"failed to store plp in redis",

			zap.String(
				"cache_key",
				cacheKey,
			),

			zap.Error(err),
		)

	} else {

		logger.Log.Info(
			"plp cached successfully",

			zap.String(
				"cache_key",
				cacheKey,
			),

			zap.Duration(
				"ttl",
				redisClient.PLPTTL,
			),
		)
	}

	return &response, nil
}
