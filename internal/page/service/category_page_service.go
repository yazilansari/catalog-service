package service

import (
	"catalog-service/internal/logger"
	"catalog-service/internal/page/dto"
	"catalog-service/internal/page/repository"
	"encoding/json"
	"time"

	redisClient "catalog-service/internal/redis"

	"go.uber.org/zap"
)

func GetCategoryPage(
	tenantCode string,
	countryCode string,
	slug string,
) (interface{}, error) {

	cacheKey :=
		"page:" +
			"category:" +
			tenantCode +
			":" +
			countryCode +
			":" +
			slug

	// =========================
	// CACHE
	// =========================

	logger.Log.Info(
		"get category page request",

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
			"categories",
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

	cached, err :=
		redisClient.Client.Get(
			redisClient.Ctx,
			cacheKey,
		).Result()

	// CACHE HIT

	if err == nil {

		logger.Log.Info(
			"category page cache hit",

			zap.String(
				"cache_key",
				cacheKey,
			),
		)

		var data dto.CategoryPageResponse

		err = json.Unmarshal(
			[]byte(cached),
			&data,
		)

		if err == nil {
			logger.Log.Info(
				"category page cache unmarshal success",

				zap.String(
					"cache_key",
					cacheKey,
				),
			)

			return &data, nil
		}

		logger.Log.Error(
			"category page cache unmarshal failed",

			zap.String(
				"cache_key",
				cacheKey,
			),

			zap.Error(err),
		)
	}

	// CACHE MISS

	if err != nil {

		logger.Log.Warn(
			"redis cache miss",

			zap.String(
				"cache_key",
				cacheKey,
			),

			zap.Error(err),
		)
	}

	// =========================
	// DATABASE FALLBACK
	// =========================

	logger.Log.Info(
		"fetching category page from database",

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

	category, err := repository.FindCategoryBySlug(
		tenantCode,
		countryCode,
		slug,
	)

	categoryPageDuration := time.Since(start)

	if err != nil {

		logger.Log.Error(
			"failed to fetch category page",

			zap.String("tenant_code", tenantCode),
			zap.String("country_code", countryCode),
			zap.String(
				"type",
				"categories",
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

	subCategories, err := repository.GetSubCategories(
		category.ID,
	)

	subCategoryPageDuration := time.Since(start)

	if err != nil {
		logger.Log.Error(
			"failed to fetch sub category page",

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

			zap.Duration(
				"sub_category_page_duration",
				subCategoryPageDuration,
			),

			zap.Error(err),
		)
		return nil, err
	}

	products, err := repository.GetProductsByCategory(
		category.ID,
	)

	productPageDuration := time.Since(start)

	if err != nil {
		logger.Log.Error(
			"failed to fetch product page",

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

			zap.Duration(
				"product_page_duration",
				productPageDuration,
			),

			zap.Error(err),
		)
		return nil, err
	}

	duration := time.Since(start)

	logger.Log.Info(
		"page fetched from database",

		zap.Duration(
			"query_duration",
			duration,
		),

		zap.Int(
			"subcategories_count",
			len(subCategories),
		),

		zap.Int(
			"products_count",
			len(products),
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
				"slug",
				slug,
			),
		)
	}

	// =========================
	// RESPONSE
	// =========================

	data := dto.CategoryPageResponse{
		PageType: "category",

		Category: *category,

		SubCategories: subCategories,

		Products: products,
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
			"failed to store category page in redis",

			zap.String(
				"cache_key",
				cacheKey,
			),

			zap.Error(err),
		)

	} else {

		logger.Log.Info(
			"category page cached successfully",

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
