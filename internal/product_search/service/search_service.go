package service

import (
	"catalog-service/internal/logger"
	"catalog-service/internal/product_search/dto"
	"catalog-service/internal/product_search/model"
	"catalog-service/internal/product_search/repository"
	promotionService "catalog-service/internal/promotion/service"
	redisClient "catalog-service/internal/redis"
	"encoding/json"
	"strconv"
	"time"

	"go.uber.org/zap"
)

func SearchProducts(
	query dto.ProductSearchQuery,
	tenantCode string,
	countryCode string,
) (
	*dto.ProductSearchResponse,
	error,
) {

	logger.Log.Info(
		"search products service called",
		zap.String("tenant_code", tenantCode),
		zap.String("country_code", countryCode),
		zap.Any("query", query),
	)

	// =========================
	// DEFAULTS
	// =========================

	if query.Limit <= 0 {
		query.Limit = 20
	}

	if query.Limit > 100 {
		query.Limit = 100
	}

	if query.Sort == "" {
		query.Sort = "relevance"
	}

	// =========================
	// CACHE KEY
	// =========================

	cacheKey :=
		buildCacheKey(
			query,
			tenantCode,
			countryCode,
		)

	logger.Log.Info(
		"product search request",

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
	// 		"product search cache hit",

	// 		zap.String(
	// 			"cache_key",
	// 			cacheKey,
	// 		),
	// 	)

	// 	var response dto.ProductSearchResponse

	// 	err =
	// 		json.Unmarshal(
	// 			[]byte(cached),
	// 			&response,
	// 		)

	// 	if err == nil {

	// 		logger.Log.Info(
	// 			"product search unmarshal success",

	// 			zap.String(
	// 				"cache_key",
	// 				cacheKey,
	// 			),
	// 		)

	// 		return &response, nil
	// 	}

	// 	logger.Log.Error(
	// 		"product search cache unmarshal failed",

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
	// 		"product search redis cache miss",

	// 		zap.String(
	// 			"cache_key",
	// 			cacheKey,
	// 		),

	// 		zap.Error(err),
	// 	)
	// }

	redisStart := time.Now()

	cached, err :=
		redisClient.GetCache[dto.ProductSearchResponse](
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
			"slow product details query detected",

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

	var products []model.ProductSearchDocument

	for _, row := range rows {

		data, _ :=
			json.Marshal(row)

		var product model.ProductSearchDocument

		_ =
			json.Unmarshal(
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
	// PROMOTIONS
	// =========================

	productIDs :=
		make(
			[]uint64,
			0,
			len(products),
		)

	priceMap :=
		make(
			map[uint64]float64,
		)

	for _, product := range products {

		productIDs =
			append(
				productIDs,
				product.ID,
			)

		priceMap[product.ID] =
			product.Price
	}

	promotionMap,
		err :=
		promotionService.GetProductsPromotions(
			tenantCode,
			countryCode,
			productIDs,
			priceMap,
		)

	if err != nil {

		logger.Log.Error(
			"failed to fetch product promotions",

			zap.Error(err),
		)
	}

	for i := range products {

		if promotion,
			ok := promotionMap[products[i].ID]; ok {

			products[i].Promotion =
				promotion
		}
	}

	// =========================
	// PAGINATION
	// =========================

	hasMore := int64(len(products)) < total

	nextCursor := ""

	if hasMore &&
		len(products) > 0 {

		last :=
			products[len(products)-1]

		switch v :=
			last.Cursor.(type) {

		case float64:

			nextCursor =
				strconv.FormatFloat(
					v,
					'f',
					-1,
					64,
				)

		case string:

			nextCursor = v
		}
	}

	// =========================
	// BRAND FILTERS
	// =========================

	var brands []string

	if brandsAgg, ok :=
		aggs["brands"].(map[string]interface{}); ok {

		if buckets, ok :=
			brandsAgg["buckets"].([]interface{}); ok {

			for _, bucket := range buckets {

				brand :=
					bucket.(map[string]interface{})["key"]

				if value, ok :=
					brand.(string); ok {

					brands =
						append(
							brands,
							value,
						)
				}
			}
		}
	}

	// =========================
	// PRICE RANGE
	// =========================

	var minPrice float64
	var maxPrice float64

	if value, ok :=
		aggs["min_price"].(map[string]interface{})["value"]; ok &&
		value != nil {

		minPrice =
			value.(float64)
	}

	if value, ok :=
		aggs["max_price"].(map[string]interface{})["value"]; ok &&
		value != nil {

		maxPrice =
			value.(float64)
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
		dto.ProductSearchResponse{
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
	// CACHE RESPONSE
	// =========================

	// jsonData, _ :=
	// 	json.Marshal(
	// 		response,
	// 	)

	// err =
	// 	redisClient.Client.Set(
	// 		redisClient.Ctx,
	// 		cacheKey,
	// 		jsonData,
	// 		time.Hour,
	// 	).Err()

	// if err != nil {

	// 	logger.Log.Error(
	// 		"failed to cache product search",

	// 		zap.String(
	// 			"cache_key",
	// 			cacheKey,
	// 		),

	// 		zap.Error(err),
	// 	)
	// } else {

	// 	logger.Log.Info(
	// 		"product search cached successfully",

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
		redisClient.SearchTTL,
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
			"failed to cache product search",

			zap.String(
				"cache_key",
				cacheKey,
			),

			zap.Error(err),
		)
	} else {

		logger.Log.Info(
			"product search cached successfully",

			zap.String(
				"cache_key",
				cacheKey,
			),

			zap.Duration(
				"ttl",
				redisClient.SearchTTL,
			),
		)
	}

	return &response, nil
}
