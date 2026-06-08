package service

import (
	"catalog-service/internal/logger"
	"catalog-service/internal/product/dto"
	"catalog-service/internal/product/repository"
	promotionService "catalog-service/internal/promotion/service"
	redisClient "catalog-service/internal/redis"
	"time"

	"go.uber.org/zap"
)

func GetProductPage(
	tenantCode string,
	countryCode string,
	slug string,
) (*dto.ProductResponse, error) {

	cacheKey :=
		"product:" +
			tenantCode +
			":" +
			countryCode +
			":" +
			slug

	// =========================
	// LOG REQUEST
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
			"slug",
			slug,
		),

		zap.String(
			"cache_key",
			cacheKey,
		),
	)

	// =========================
	// REDIS CACHE
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
	// 		"product cache hit",

	// 		zap.String(
	// 			"cache_key",
	// 			cacheKey,
	// 		),
	// 	)

	// 	var data dto.ProductResponse

	// 	err = json.Unmarshal(
	// 		[]byte(cached),
	// 		&data,
	// 	)

	// 	if err == nil {

	// 		logger.Log.Info(
	// 			"product cache unmarshal success",

	// 			zap.String(
	// 				"cache_key",
	// 				cacheKey,
	// 			),
	// 		)

	// 		return &data, nil
	// 	}

	// 	logger.Log.Error(
	// 		"product cache unmarshal failed",

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
	// 		"product redis cache miss",

	// 		zap.String(
	// 			"cache_key",
	// 			cacheKey,
	// 		),

	// 		zap.Error(err),
	// 	)
	// }

	redisStart := time.Now()

	cached, err :=
		redisClient.GetCache[dto.ProductResponse](
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
	// PRODUCT
	// =========================

	start := time.Now()

	product, err :=
		repository.GetProductBySlug(
			tenantCode,
			countryCode,
			slug,
		)

	if err != nil {

		logger.Log.Error(
			"failed to fetch product",

			zap.Error(err),
		)

		return nil, err
	}

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

	logger.Log.Info(
		"product fetched",

		zap.Uint64(
			"product_id",
			product.ID,
		),
	)

	// =========================
	// CATEGORY + SUBCATEGORY
	// =========================

	category,
		subCategory,
		err :=
		repository.GetCategoryAndSubCategory(
			product.ID,
		)

	if err != nil {

		logger.Log.Error(
			"failed to fetch category mapping",

			zap.Uint64(
				"product_id",
				product.ID,
			),

			zap.Error(err),
		)

		return nil, err
	}

	// =========================
	// IMAGES
	// =========================

	images, err :=
		repository.GetProductImages(
			product.ID,
		)

	if err != nil {
		return nil, err
	}

	// =========================
	// VARIANTS
	// =========================

	variantsDB, err :=
		repository.GetVariantsByProductID(
			product.ID,
		)

	if err != nil {
		return nil, err
	}

	var variants []dto.ProductVariantResponse

	for _, v := range variantsDB {

		variants =
			append(
				variants,

				dto.ProductVariantResponse{
					ID: v.ID,

					Name: v.Name,

					SKU: v.SKU,

					Price: v.Price,

					DiscountPrice: v.DiscountPrice,

					Stock: v.Stock,
				},
			)
	}

	// =========================
	// FRAGRANCE NOTES
	// =========================

	notesDB, err :=
		repository.GetFragranceNotes(
			product.ID,
		)

	if err != nil {
		return nil, err
	}

	var fragranceNotes []dto.FragranceNoteResponse

	for _, n := range notesDB {

		fragranceNotes =
			append(
				fragranceNotes,

				dto.FragranceNoteResponse{
					TopNote: n.TopNote,

					HeartNote: n.HeartNote,

					BaseNote: n.BaseNote,

					TopNoteImage: n.TopNoteImage,

					HeartNoteImage: n.HeartNoteImage,

					BaseNoteImage: n.BaseNoteImage,

					TopNoteDescription: n.TopNoteDescription,

					HeartNoteDescription: n.HeartNoteDescription,

					BaseNoteDescription: n.BaseNoteDescription,
				},
			)
	}

	// =========================
	// RELATED PRODUCTS
	// =========================

	relatedDB, err :=
		repository.GetRelatedProducts(
			category.ID,
			product.ID,
		)

	if err != nil {

		logger.Log.Error(
			"failed to fetch related products",

			zap.Error(err),
		)

		return nil, err
	}

	productIDs := make([]uint64, 0, len(relatedDB))
	priceMap := make(map[uint64]float64)

	for _, p := range relatedDB {
		productIDs = append(productIDs, p.ID)
		priceMap[p.ID] = p.Price
	}

	relatedPromotionMap, err :=
		promotionService.GetProductsPromotions(
			tenantCode,
			countryCode,
			productIDs,
			priceMap,
		)

	if err != nil {

		logger.Log.Error(
			"failed to fetch related products promotions",

			zap.Error(err),
		)

		return nil, err
	}

	for i := range relatedDB {

		if promotion,
			ok := relatedPromotionMap[relatedDB[i].ID]; ok {

			relatedDB[i].Promotion = promotion
		}
	}

	var relatedProducts []dto.RelatedProductResponse

	for _, p := range relatedDB {

		relatedProducts =
			append(
				relatedProducts,

				dto.RelatedProductResponse{
					ID: p.ID,

					Name: p.Name,

					Slug: p.Slug,

					Price: p.Price,

					// DiscountPrice: p.DiscountPrice,

					Promotion: p.Promotion,
				},
			)
	}

	// =========================
	// SEO
	// =========================

	seo, _ :=
		repository.GetProductSEO(
			slug,
		)

	// =========================
	// RESPONSE
	// =========================

	response := dto.ProductResponse{
		PageType: "product",

		Category: *category,

		SubCategory: *subCategory,

		Product: *product,

		Images: images,

		Variants: variants,

		FragranceNotes: fragranceNotes,

		RelatedProducts: relatedProducts,
	}

	if seo != nil {
		response.SEO = *seo
	}

	duration := time.Since(start)

	logger.Log.Info(
		"product built",

		zap.Duration(
			"duration",
			duration,
		),
	)

	// =========================
	// STORE CACHE
	// =========================

	// jsonData, _ :=
	// 	json.Marshal(response)

	// err =
	// 	redisClient.Client.Set(
	// 		redisClient.Ctx,
	// 		cacheKey,
	// 		jsonData,
	// 		time.Hour,
	// 	).Err()

	// if err != nil {

	// 	logger.Log.Error(
	// 		"failed to cache product",

	// 		zap.String(
	// 			"cache_key",
	// 			cacheKey,
	// 		),

	// 		zap.Error(err),
	// 	)
	// } else {

	// 	logger.Log.Info(
	// 		"product cached successfully",

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
		redisClient.ProductTTL,
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
			"failed to cache product",

			zap.String(
				"cache_key",
				cacheKey,
			),

			zap.Error(err),
		)
	} else {

		logger.Log.Info(
			"product cached successfully",

			zap.String(
				"cache_key",
				cacheKey,
			),

			zap.Duration(
				"ttl",
				redisClient.ProductTTL,
			),
		)
	}

	if duration > time.Second {

		logger.Log.Warn(
			"slow product details query detected",

			zap.Any(
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

	return &response, nil
}
