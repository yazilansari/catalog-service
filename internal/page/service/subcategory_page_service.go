package service

import (
	"catalog-service/internal/logger"
	"catalog-service/internal/page/dto"
	"catalog-service/internal/page/helper"
	"catalog-service/internal/page/repository"
	promotionService "catalog-service/internal/promotion/service"
	redisClient "catalog-service/internal/redis"
	"time"

	"go.uber.org/zap"
)

func GetSubCategoryPage(
	tenantCode string,
	countryCode string,
	slug string,
	page int,
	limit int,
	sort string,
	brand string,
	priceMin float64,
	priceMax float64,
) (interface{}, error) {

	// cacheKey :=
	// 	"page:" +
	// 		"subcategory:" +
	// 		tenantCode +
	// 		":" +
	// 		countryCode +
	// 		":" +
	// 		slug

	cacheKey := buildCacheKey(
		tenantCode,
		countryCode,
		slug,
		page,
		sort,
		limit,
		priceMin,
		priceMax,
	)

	// =========================
	// CACHE
	// =========================

	logger.Log.Info(
		"get subcategory page request",

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
			"subcategories",
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
	// 		"subcategory page cache hit",

	// 		zap.String(
	// 			"cache_key",
	// 			cacheKey,
	// 		),
	// 	)

	// 	var data dto.SubCategoryPageResponse

	// 	err = json.Unmarshal(
	// 		[]byte(cached),
	// 		&data,
	// 	)

	// 	if err == nil {
	// 		logger.Log.Info(
	// 			"subcategory page cache unmarshal success",

	// 			zap.String(
	// 				"cache_key",
	// 				cacheKey,
	// 			),
	// 		)

	// 		return &data, nil
	// 	}

	// 	logger.Log.Error(
	// 		"subcategory page cache unmarshal failed",

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
		redisClient.GetCache[dto.SubCategoryPageResponse](
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
		"fetching subcategory page from database",

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

	subCategory, err :=
		repository.FindSubCategoryBySlug(
			tenantCode,
			countryCode,
			slug,
		)

	subCategoryPageDuration := time.Since(start)

	if err != nil {

		logger.Log.Error(
			"failed to fetch subcategory page",

			zap.String("tenant_code", tenantCode),
			zap.String("country_code", countryCode),
			zap.String(
				"type",
				"subcategories",
			),
			zap.String("slug", slug),

			zap.Duration(
				"subcategory_page_duration",
				subCategoryPageDuration,
			),

			zap.Error(err),
		)

		return nil, err
	}

	category, err :=
		repository.GetCategoryByID(
			*subCategory.ParentID,
		)

	categoryPageDuration := time.Since(start)

	if err != nil {

		logger.Log.Error(
			"failed to fetch category page",

			zap.String("tenant_code", tenantCode),
			zap.String("country_code", countryCode),
			zap.String(
				"type",
				"subcategories",
			),
			zap.String("slug", slug),

			zap.Duration(
				"category_page_duration",
				categoryPageDuration,
			),

			zap.Error(err),
		)

		return nil, err
	}

	products, total, err :=
		repository.GetProductsBySubCategory(
			subCategory.ID,
			page,
			limit,
			sort,
			brand,
			priceMin,
			priceMax,
		)

	productPageDuration := time.Since(start)

	if err != nil {

		logger.Log.Error(
			"failed to fetch product page",

			zap.String("tenant_code", tenantCode),
			zap.String("country_code", countryCode),
			zap.String(
				"type",
				"subcategories",
			),
			zap.String("slug", slug),

			zap.Duration(
				"product_page_duration",
				productPageDuration,
			),

			zap.Int64(
				"total",
				total,
			),

			zap.Error(err),
		)

		return nil, err
	}

	promotionMap, err :=
		promotionService.GetProductsPromotions(
			tenantCode,
			countryCode,
			helper.ExtractProductIDs(products),
			helper.BuildPriceMap(products),
		)

	if err != nil {
		logger.Log.Warn(
			"failed to fetch promotions",
			zap.Error(err),
		)
	}

	duration := time.Since(start)

	for i := range products {

		if promotion, ok := promotionMap[products[i].ID]; ok {

			products[i].Promotion = promotion

		}
	}

	logger.Log.Info(
		"subcategory page fetched successfully",

		zap.String(
			"slug",
			slug,
		),

		zap.Uint64(
			"subcategory_id",
			subCategory.ID,
		),

		zap.Int(
			"products_count",
			len(products),
		),

		zap.Duration(
			"duration",
			duration,
		),
	)

	if duration > time.Second {

		logger.Log.Warn(
			"slow subcategory query detected",

			zap.String(
				"slug",
				slug,
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
	// RESPONSE
	// =========================

	data := dto.SubCategoryPageResponse{
		PageType: "subcategory",

		Category: *category,

		SubCategory: *subCategory,

		Products: products,
	}

	// =========================
	// STORE CACHE
	// =========================

	// jsonData, _ := json.Marshal(
	// 	data,
	// )

	// err = redisClient.Client.Set(
	// 	redisClient.Ctx,
	// 	cacheKey,
	// 	jsonData,
	// 	time.Hour,
	// ).Err()

	// if err != nil {

	// 	logger.Log.Error(
	// 		"failed to store subcategory page in redis",

	// 		zap.String(
	// 			"cache_key",
	// 			cacheKey,
	// 		),

	// 		zap.Error(err),
	// 	)

	// } else {

	// 	logger.Log.Info(
	// 		"subcategory page cached successfully",

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
			"failed to store subcategory page in redis",

			zap.String(
				"cache_key",
				cacheKey,
			),

			zap.Error(err),
		)

	} else {

		logger.Log.Info(
			"subcategory page cached successfully",

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

	return &data, nil
}
