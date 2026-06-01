package service

import (
	"catalog-service/internal/logger"
	"catalog-service/internal/plp/dto"
	"catalog-service/internal/plp/model"
	"catalog-service/internal/plp/repository"
	redisClient "catalog-service/internal/redis"
	"encoding/json"
	"strconv"
	"time"

	"go.uber.org/zap"
)

func GetProducts(
	query dto.ProductQuery,
) (*dto.ProductListResponse, error) {

	// =========================
	// DEFAULTS
	// =========================

	if query.Limit <= 0 {
		query.Limit = 20
	}

	if query.Sort == "" {
		query.Sort = "latest"
	}

	// =========================
	// CACHE KEY
	// =========================

	cacheKey :=
		"plp:" +
			query.Category + ":" +
			query.SubCategory + ":" +
			query.Brand + ":" +
			query.Sort + ":" +
			strconv.Itoa(
				query.Limit,
			) + ":" +
			query.Cursor

	logger.Log.Info(
		"get products request",

		zap.String(
			"cache_key",
			cacheKey,
		),
	)

	// =========================
	// CACHE
	// =========================

	cached, err :=
		redisClient.Client.Get(
			redisClient.Ctx,
			cacheKey,
		).Result()

	if err == nil {

		var response dto.ProductListResponse

		err =
			json.Unmarshal(
				[]byte(cached),
				&response,
			)

		if err == nil {

			logger.Log.Info(
				"plp cache hit",

				zap.String(
					"cache_key",
					cacheKey,
				),
			)

			return &response, nil
		}
	}

	logger.Log.Info(
		"plp cache miss",
	)

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

	// =========================
	// MAP PRODUCTS
	// =========================

	var products []model.ProductDocument

	for _, row := range rows {

		data, _ :=
			json.Marshal(
				row,
			)

		var product model.ProductDocument

		_ = json.Unmarshal(
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
	// PAGINATION
	// =========================

	hasMore :=
		len(products) >= query.Limit

	nextCursor := ""

	if hasMore {

		last :=
			products[len(products)-1]

		nextCursor =
			strconv.FormatUint(
				last.ID,
				10,
			)
	}

	// =========================
	// FILTERS
	// =========================

	var brands []string

	brandsAgg, ok :=
		aggs["brands"].(map[string]interface{})

	if ok {

		buckets, ok :=
			brandsAgg["buckets"].([]interface{})

		if ok {

			for _, bucket := range buckets {

				brand :=
					bucket.(map[string]interface{})["key"]

				brands =
					append(
						brands,
						brand.(string),
					)
			}
		}
	}

	var minPrice float64
	var maxPrice float64

	// =========================
	// MIN PRICE
	// =========================

	if value, ok :=
		aggs["min_price"].(map[string]interface{})["value"]; ok && value != nil {

		minPrice =
			value.(float64)
	}

	// =========================
	// MAX PRICE
	// =========================

	if value, ok :=
		aggs["max_price"].(map[string]interface{})["value"]; ok && value != nil {

		maxPrice =
			value.(float64)
	}

	// =========================
	// FILTERS
	// =========================

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
		dto.ProductListResponse{
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
	// STORE CACHE
	// =========================

	jsonData, _ :=
		json.Marshal(
			response,
		)

	_ = redisClient.Client.Set(
		redisClient.Ctx,
		cacheKey,
		jsonData,
		15*time.Minute,
	).Err()

	return &response, nil
}
