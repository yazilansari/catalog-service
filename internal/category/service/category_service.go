package service

import (
	"time"

	"catalog-service/internal/category/dto"
	"catalog-service/internal/category/model"
	"catalog-service/internal/category/repository"
	"catalog-service/internal/logger"

	redisClient "catalog-service/internal/redis"

	"encoding/json"

	"go.uber.org/zap"
)

func GetCategoryTree(
	tenantCode string,
	countryCode string,
) ([]*dto.CategoryResponse, error) {

	cacheKey :=
		"categories:" +
			tenantCode +
			":" +
			countryCode

	logger.Log.Info(
		"get category tree request",

		zap.String(
			"tenant_code",
			tenantCode,
		),

		zap.String(
			"country_code",
			countryCode,
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
	// 		"redis cache hit",

	// 		zap.String(
	// 			"cache_key",
	// 			cacheKey,
	// 		),
	// 	)

	// 	var data []*dto.CategoryResponse

	// 	err = json.Unmarshal(
	// 		[]byte(cached),
	// 		&data,
	// 	)

	// 	if err == nil {
	// 		logger.Log.Info(
	// 			"redis unmarshal success",

	// 			zap.String(
	// 				"cache_key",
	// 				cacheKey,
	// 			),

	// 			zap.Int(
	// 				"root_categories",
	// 				len(data),
	// 			),
	// 		)

	// 		return data, nil
	// 	}

	// 	logger.Log.Error(
	// 		"redis unmarshal failed",

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
		redisClient.GetCache[[]*dto.CategoryResponse](
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

		return *cached, nil
	}

	// =========================
	// DATABASE FALLBACK
	// =========================

	logger.Log.Info(
		"fetching categories from database",

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

	categories, err := repository.GetCategories(
		tenantCode,
		countryCode,
	)

	duration := time.Since(start)

	if err != nil {
		logger.Log.Error(
			"failed to fetch categories from database",

			zap.String(
				"tenant_code",
				tenantCode,
			),

			zap.String(
				"country_code",
				countryCode,
			),

			zap.Error(err),
		)
		return nil, err
	}

	logger.Log.Info(
		"categories fetched from database",

		zap.Int(
			"count",
			len(categories),
		),

		zap.Duration(
			"query_duration",
			duration,
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
		)
	}

	tree := buildCategoryTree(categories)

	logger.Log.Info(
		"category tree built successfully",

		zap.Int(
			"root_categories",
			len(tree),
		),
	)

	// =========================
	// STORE CACHE
	// =========================

	jsonData, _ := json.Marshal(tree)

	err = redisClient.Client.Set(
		redisClient.Ctx,
		cacheKey,
		jsonData,
		time.Hour,
	).Err()

	if err != nil {

		logger.Log.Error(
			"failed to store categories in redis",

			zap.String(
				"cache_key",
				cacheKey,
			),

			zap.Error(err),
		)

	} else {

		logger.Log.Info(
			"categories cached successfully",

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

	return tree, nil
}

func buildCategoryTree(
	categories []model.Category,
) []*dto.CategoryResponse {

	logger.Log.Info(
		"building category tree",

		zap.Int(
			"categories_count",
			len(categories),
		),
	)

	categoryMap := make(map[uint64]*dto.CategoryResponse)

	var roots []*dto.CategoryResponse

	// Step 1: create all nodes
	for _, cat := range categories {

		categoryMap[cat.ID] = &dto.CategoryResponse{
			ID:                   cat.ID,
			Name:                 cat.Name,
			Slug:                 cat.Slug,
			Image:                cat.Image,
			IconImage:            cat.IconImage,
			MenuImage:            cat.MenuImage,
			MenuImage2:           cat.MenuImage2,
			MobileImage:          cat.MobileImage,
			Video:                cat.Video,
			ProductSubCategories: []*dto.CategoryResponse{},
		}
	}

	logger.Log.Info(
		"category map created",

		zap.Int(
			"map_size",
			len(categoryMap),
		),
	)

	// Step 2: build hierarchy
	for _, cat := range categories {

		node := categoryMap[cat.ID]

		// root node
		if cat.ParentID == nil || *cat.ParentID == 0 {
			roots = append(roots, node)
			continue
		}

		parent, ok := categoryMap[*cat.ParentID]
		if !ok {
			logger.Log.Warn(
				"missing parent category",

				zap.Uint64(
					"category_id",
					cat.ID,
				),

				zap.Uint64(
					"parent_id",
					*cat.ParentID,
				),
			)
			continue
		}

		parent.ProductSubCategories = append(parent.ProductSubCategories, node)

		logger.Log.Debug(
			"child category attached",

			zap.Uint64(
				"child_id",
				node.ID,
			),

			zap.Uint64(
				"parent_id",
				parent.ID,
			),
		)
	}

	logger.Log.Info(
		"category tree completed",

		zap.Int(
			"root_categories",
			len(roots),
		),
	)

	return roots
}
