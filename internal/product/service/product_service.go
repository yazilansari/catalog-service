package service

import (
	"catalog-service/internal/logger"
	"catalog-service/internal/product/dto"
	"catalog-service/internal/product/repository"
	redisClient "catalog-service/internal/redis"
	"encoding/json"
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

	cached, err :=
		redisClient.GetCache[dto.ProductResponse](
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
		return nil, err
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

					DiscountPrice: p.DiscountPrice,
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

	jsonData, _ :=
		json.Marshal(response)

	err =
		redisClient.Client.Set(
			redisClient.Ctx,
			cacheKey,
			jsonData,
			time.Hour,
		).Err()

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
				time.Hour,
			),
		)
	}

	return &response, nil
}
