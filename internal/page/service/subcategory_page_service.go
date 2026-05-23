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

func GetSubCategoryPage(
	tenantCode string,
	countryCode string,
	slug string,
) (interface{}, error) {

	cacheKey :=
		"page:" +
			"subcategory:" +
			tenantCode +
			":" +
			countryCode +
			":" +
			slug

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

	cached, err :=
		redisClient.Client.Get(
			redisClient.Ctx,
			cacheKey,
		).Result()

	// CACHE HIT

	if err == nil {

		logger.Log.Info(
			"subcategory page cache hit",

			zap.String(
				"cache_key",
				cacheKey,
			),
		)

		var data dto.SubCategoryPageResponse

		err = json.Unmarshal(
			[]byte(cached),
			&data,
		)

		if err == nil {
			logger.Log.Info(
				"subcategory page cache unmarshal success",

				zap.String(
					"cache_key",
					cacheKey,
				),
			)

			return &data, nil
		}

		logger.Log.Error(
			"subcategory page cache unmarshal failed",

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

	products, err :=
		repository.GetProductsBySubCategory(
			subCategory.ID,
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

			zap.Error(err),
		)

		return nil, err
	}

	duration := time.Since(start)

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

	jsonData, _ := json.Marshal(
		data,
	)

	err = redisClient.Client.Set(
		redisClient.Ctx,
		cacheKey,
		jsonData,
		time.Hour,
	).Err()

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
				time.Hour,
			),
		)
	}

	return &data, nil
}
