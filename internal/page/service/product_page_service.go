package service

import (
	"catalog-service/internal/logger"
	"catalog-service/internal/page/dto"
	"catalog-service/internal/page/helper"
	"catalog-service/internal/page/repository"
	productRepository "catalog-service/internal/product/repository"
	promotionService "catalog-service/internal/promotion/service"
	redisClient "catalog-service/internal/redis"
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

	// cacheKey :=
	// 	"page:" +
	// 		"product:" +
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

	redisStart := time.Now()

	cached, err :=
		redisClient.GetCache[dto.ProductPageResponse](
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

	productPromotion, err :=
		promotionService.GetProductPromotions(
			tenantCode,
			countryCode,
			product.ID,
			product.Price,
		)

	if err != nil {
		logger.Log.Warn(
			"failed to fetch promotions",
			zap.Error(err),
		)
	}

	product.Promotion = productPromotion

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

	// =========================
	// IMAGES
	// =========================

	images, err :=
		productRepository.GetProductImages(
			product.ID,
		)

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

			zap.Error(err),
		)

		return nil, err
	}

	imageDTOs := make([]dto.ProductImage, 0, len(images))

	for _, img := range images {
		imageDTOs = append(imageDTOs, dto.ProductImage{
			ID:        img.ID,
			Image:     img.Image,
			SortOrder: img.SortOrder,
		})
	}

	// =========================
	// RELATED PRODUCTS
	// =========================

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

	relatedPromotionMap, err :=
		promotionService.GetProductsPromotions(
			tenantCode,
			countryCode,
			helper.ExtractProductIDs(relatedProducts),
			helper.BuildPriceMap(relatedProducts),
		)

	for i := range relatedProducts {

		if promotion,
			ok := relatedPromotionMap[relatedProducts[i].ID]; ok {

			relatedProducts[i].Promotion = promotion
		}
	}

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

	data := dto.ProductPageResponse{
		PageType: "product",

		Category: *category,

		SubCategory: *subCategory,

		Images: imageDTOs,

		Product: *product,

		RelatedProducts: relatedProducts,
	}

	// =========================
	// CACHE STORE
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
	// 		"failed to cache product page",

	// 		zap.String(
	// 			"cache_key",
	// 			cacheKey,
	// 		),

	// 		zap.Error(err),
	// 	)

	// } else {

	// 	logger.Log.Info(
	// 		"product page cached successfully",

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
				redisClient.PageTTL,
			),
		)
	}

	return &data, nil
}
