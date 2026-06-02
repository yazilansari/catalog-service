package repository

import (
	"bytes"
	"catalog-service/internal/elasticsearch"
	"catalog-service/internal/logger"
	"catalog-service/internal/plp/dto"
	"context"
	"encoding/json"

	"go.uber.org/zap"
)

func SearchProducts(
	query dto.ProductQuery,
	tenantCode string,
	countryCode string,
) (
	[]map[string]interface{},
	map[string]interface{},
	int64,
	error,
) {

	indexName :=
		elasticsearch.GetProductIndex(
			tenantCode,
			countryCode,
		)

	logger.Log.Debug(
		"index name",
		zap.String("indexName", indexName),
	)

	filter :=
		[]map[string]interface{}{}

	// =========================
	// CATEGORY
	// =========================

	if query.Category != "" {

		logger.Log.Debug(
			"category query",
		)

		filter =
			append(
				filter,

				map[string]interface{}{
					"term": map[string]interface{}{
						"category": query.Category,
					},
				},
			)
	}

	// =========================
	// SUBCATEGORY
	// =========================

	if query.SubCategory != "" {

		logger.Log.Debug(
			"subcategory query",
		)

		filter =
			append(
				filter,

				map[string]interface{}{
					"term": map[string]interface{}{
						"subcategory": query.SubCategory,
					},
				},
			)
	}

	// =========================
	// BRAND
	// =========================

	if query.Brand != "" {

		logger.Log.Debug(
			"brand query",
		)

		filter =
			append(
				filter,

				map[string]interface{}{
					"term": map[string]interface{}{
						"brand": query.Brand,
					},
				},
			)
	}

	// =========================
	// PRICE RANGE
	// =========================

	if query.MinPrice > 0 ||
		query.MaxPrice > 0 {

		logger.Log.Debug(
			"price range query",
		)

		priceRange :=
			map[string]interface{}{}

		if query.MinPrice > 0 {

			priceRange["gte"] =
				query.MinPrice
		}

		if query.MaxPrice > 0 {

			priceRange["lte"] =
				query.MaxPrice
		}

		filter =
			append(
				filter,

				map[string]interface{}{
					"range": map[string]interface{}{
						"price": priceRange,
					},
				},
			)
	}

	logger.Log.Debug(
		"filter",
		zap.Any("filter", filter),
	)

	// =========================
	// SORT
	// =========================

	sort :=
		buildSort(
			query.Sort,
		)

	searchQuery :=
		map[string]interface{}{
			"size": query.Limit,

			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"filter": filter,
				},
			},

			"sort": sort,

			"aggs": map[string]interface{}{
				"brands": map[string]interface{}{
					"terms": map[string]interface{}{
						"field": "brand",

						"size": 100,
					},
				},

				"min_price": map[string]interface{}{
					"min": map[string]interface{}{
						"field": "price",
					},
				},

				"max_price": map[string]interface{}{
					"max": map[string]interface{}{
						"field": "price",
					},
				},
			},
		}

	// =========================
	// CURSOR PAGINATION
	// =========================

	if query.Cursor != "" {

		logger.Log.Debug(
			"cursor pagination",
		)

		searchQuery["search_after"] =
			[]interface{}{
				query.Cursor,
			}
	}

	logger.Log.Debug(
		"search query",
		zap.Any("searchQuery", searchQuery),
	)

	body, err :=
		json.Marshal(
			searchQuery,
		)

	if err != nil {

		logger.Log.Error(
			"failed to marshal search query",

			zap.Error(err),
		)

		return nil, nil, 0, err
	}

	res, err :=
		elasticsearch.Client.Search(
			elasticsearch.Client.Search.WithContext(
				context.Background(),
			),

			elasticsearch.Client.Search.WithIndex(
				indexName,
			),

			elasticsearch.Client.Search.WithBody(
				bytes.NewReader(body),
			),
		)

	if res.IsError() {

		logger.Log.Error(
			"elasticsearch search failed",

			zap.String(
				"status",
				res.Status(),
			),
		)

		return nil, nil, 0, err
	}

	if err != nil {

		logger.Log.Error(
			"failed to search products",

			zap.Error(err),
		)

		return nil, nil, 0, err
	}

	defer res.Body.Close()

	var result map[string]interface{}

	err = json.NewDecoder(
		res.Body,
	).Decode(
		&result,
	)

	if err != nil {

		logger.Log.Error(
			"failed to decode search response",

			zap.Error(err),
		)

		return nil, nil, 0, err
	}

	// =========================
	// HITS
	// =========================

	hitsMap, ok :=
		result["hits"].(map[string]interface{})

	if !ok {

		logger.Log.Error(
			"failed to get hits from search response",
		)

		return []map[string]interface{}{},
			nil,
			0,
			nil
	}

	total := int64(0)

	if totalMap, ok :=
		hitsMap["total"].(map[string]interface{}); ok {

		if value, ok :=
			totalMap["value"].(float64); ok {

			total =
				int64(value)
		}
	}

	rows, ok :=
		hitsMap["hits"].([]interface{})

	if !ok {

		logger.Log.Error(
			"failed to get hits from search response",
		)

		return []map[string]interface{}{},
			nil,
			total,
			nil
	}

	var products []map[string]interface{}

	for _, row := range rows {

		rowMap, ok :=
			row.(map[string]interface{})

		if !ok {
			continue
		}

		source, ok :=
			rowMap["_source"].(map[string]interface{})

		if !ok {
			continue
		}

		// attach cursor

		if sortValues, ok :=
			rowMap["sort"].([]interface{}); ok &&
			len(sortValues) > 0 {

			source["cursor"] =
				sortValues[0]
		}

		products =
			append(
				products,
				source,
			)
	}

	// =========================
	// AGGREGATIONS
	// =========================

	aggs :=
		map[string]interface{}{}

	if aggregations, ok :=
		result["aggregations"].(map[string]interface{}); ok {

		aggs =
			aggregations
	}

	return products,
		aggs,
		total,
		nil
}

func buildSort(
	sort string,
) []map[string]interface{} {

	switch sort {

	case "price_asc":

		return []map[string]interface{}{
			{
				"price": map[string]interface{}{
					"order": "asc",
				},
			},
		}

	case "price_desc":

		return []map[string]interface{}{
			{
				"price": map[string]interface{}{
					"order": "desc",
				},
			},
		}

	case "best_seller":

		return []map[string]interface{}{
			{
				"sales_count": map[string]interface{}{
					"order": "desc",
				},
			},
			{
				"id": map[string]interface{}{
					"order": "desc",
				},
			},
		}

	default:

		return []map[string]interface{}{
			{
				"id": map[string]interface{}{
					"order": "desc",
				},
			},
		}
	}
}
