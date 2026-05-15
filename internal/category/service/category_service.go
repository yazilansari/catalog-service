package service

import (
	"log"
	"time"

	"catalog-service/internal/category/dto"
	"catalog-service/internal/category/model"
	"catalog-service/internal/category/repository"

	redisClient "catalog-service/internal/redis"

	"encoding/json"
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

	log.Printf(
		"[GetCategoryTree] request tenant=%s country=%s cacheKey=%s",
		tenantCode,
		countryCode,
		cacheKey,
	)

	// =========================
	// TRY REDIS
	// =========================

	log.Printf(
		"[GetCategoryTree] checking redis cache key=%s",
		cacheKey,
	)

	cached, err :=
		redisClient.Client.Get(
			redisClient.Ctx,
			cacheKey,
		).Result()

	// CACHE HIT

	if err == nil {

		log.Printf(
			"[GetCategoryTree] redis cache HIT key=%s",
			cacheKey,
		)

		var data []*dto.CategoryResponse

		err = json.Unmarshal(
			[]byte(cached),
			&data,
		)

		if err == nil {
			log.Printf(
				"[GetCategoryTree] redis unmarshal success key=%s roots=%d",
				cacheKey,
				len(data),
			)

			return data, nil
		}

		log.Printf(
			"[GetCategoryTree] redis unmarshal failed key=%s error=%v",
			cacheKey,
			err,
		)
	}

	// CACHE MISS

	if err != nil {

		log.Printf(
			"[GetCategoryTree] redis cache MISS key=%s error=%v",
			cacheKey,
			err,
		)
	}

	// =========================
	// DATABASE FALLBACK
	// =========================

	log.Printf(
		"[GetCategoryTree] fetching categories from database",
	)

	categories, err := repository.GetCategories(
		tenantCode,
		countryCode,
	)

	if err != nil {
		log.Printf("[GetCategoryTree] repository error: %v", err)
		return nil, err
	}

	log.Printf("[GetCategoryTree] categories fetched count=%d", len(categories))

	tree := buildCategoryTree(categories)

	log.Printf("[GetCategoryTree] tree built success roots=%d", len(tree))

	// =========================
	// STORE CACHE
	// =========================

	jsonData, _ := json.Marshal(tree)

	_ = redisClient.Client.Set(
		redisClient.Ctx,
		cacheKey,
		jsonData,
		time.Hour,
	).Err()

	return tree, nil
}

func buildCategoryTree(
	categories []model.Category,
) []*dto.CategoryResponse {

	log.Printf("[buildCategoryTree] building tree from %d categories", len(categories))

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

	log.Printf("[buildCategoryTree] map built size=%d", len(categoryMap))

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
			log.Printf("[buildCategoryTree] missing parent id=%v for category id=%d", *cat.ParentID, cat.ID)
			continue
		}

		parent.ProductSubCategories = append(parent.ProductSubCategories, node)

		log.Printf("[buildCategoryTree] attached child id=%d -> parent id=%d", node.ID, parent.ID)
	}

	log.Printf("[buildCategoryTree] roots=%d", len(roots))

	return roots
}
