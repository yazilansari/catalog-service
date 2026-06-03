package service

import (
	"catalog-service/internal/logger"
	"catalog-service/internal/page/dto"
	"catalog-service/internal/page/repository"
	redisClient "catalog-service/internal/redis"
	"encoding/json"
	"time"

	"go.uber.org/zap"
)

func GetProductPage(
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

	cacheKey :=
		"page:" +
			"product:" +
			tenantCode +
			":" +
			countryCode +
			":" +
			slug

	// =========================
	// CACHE
	// =========================

	logger.Log.Info(
		"get product page request",

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
	// 		"product page cache hit",

	// 		zap.String(
	// 			"cache_key",
	// 			cacheKey,
	// 		),
	// 	)

	// 	var data dto.ProductPageResponse

	// 	err = json.Unmarshal(
	// 		[]byte(cached),
	// 		&data,
	// 	)

	// 	if err == nil {
	// 		logger.Log.Info(
	// 			"product page cache unmarshal success",

	// 			zap.String(
	// 				"cache_key",
	// 				cacheKey,
	// 			),
	// 		)

	// 		return &data, nil
	// 	}

	// 	logger.Log.Error(
	// 		"product page cache unmarshal failed",

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
		redisClient.GetCache[dto.ProductPageResponse](
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
		"fetching product page from database",

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

	product, err :=
		repository.FindProductBySlug(
			tenantCode,
			countryCode,
			slug,
		)

	ProductPageDuration := time.Since(start)

	if err != nil {

		logger.Log.Error(
			"failed to fetch product page",

			zap.String("tenant_code", tenantCode),
			zap.String("country_code", countryCode),
			zap.String(
				"type",
				"products",
			),
			zap.String("slug", slug),

			zap.Duration(
				"product_page_duration",
				ProductPageDuration,
			),

			zap.Error(err),
		)

		return nil, err
	}

	category, err :=
		repository.GetCategoryByProductID(
			product.ID,
		)

	categoryPageDuration := time.Since(start)

	if err != nil {

		logger.Log.Error(
			"failed to fetch product page",

			zap.String("tenant_code", tenantCode),
			zap.String("country_code", countryCode),
			zap.String(
				"type",
				"products",
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

	subCategory, err :=
		repository.GetSubCategoryByID(
			product.ID,
		)

	subCategoryPageDuration := time.Since(start)

	if err != nil {

		logger.Log.Error(
			"failed to fetch product page",

			zap.String("tenant_code", tenantCode),
			zap.String("country_code", countryCode),
			zap.String(
				"type",
				"products",
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

	relatedProducts, total, err :=
		repository.GetRelatedProducts(
			subCategory.ID,
			product.ID,
			page,
			limit,
			sort,
			brand,
			priceMin,
			priceMax,
		)

	relatedProductsDuration := time.Since(start)

	duration := time.Since(start)

	if err != nil {

		logger.Log.Error(
			"failed to fetch product page",

			zap.String("tenant_code", tenantCode),
			zap.String("country_code", countryCode),
			zap.String(
				"type",
				"products",
			),
			zap.String("slug", slug),

			zap.Int64(
				"total",
				total,
			),

			zap.Duration(
				"related_products_duration",
				relatedProductsDuration,
			),

			zap.Error(err),
		)

		return nil, err
	}

	logger.Log.Info(
		"product page fetched successfully",

		zap.String(
			"slug",
			slug,
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
			"slow database query detected",

			zap.Duration(
				"duration",
				duration,
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

	data := dto.ProductPageResponse{
		PageType: "product",

		Category: *category,

		SubCategory: *subCategory,

		Product: *product,

		RelatedProducts: relatedProducts,
	}

	// =========================
	// CACHE STORE
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
			"failed to cache product page",

			zap.String(
				"cache_key",
				cacheKey,
			),

			zap.Error(err),
		)

	} else {

		logger.Log.Info(
			"product page cached successfully",

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
