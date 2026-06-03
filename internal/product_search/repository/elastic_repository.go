package repository

import (
	"bytes"
	"catalog-service/internal/elasticsearch"
	"catalog-service/internal/logger"
	"catalog-service/internal/product_search/dto"
	"context"
	"encoding/json"

	"go.uber.org/zap"
)

func SearchProducts(
	query dto.ProductSearchQuery,
	tenantCode string,
	countryCode string,
) (
	[]map[string]interface{},
	map[string]interface{},
	int64,
	error,
) {

	logger.Log.Info(
		"search products request received",
		zap.String("tenant_code", tenantCode),
		zap.String("country_code", countryCode),
	)

	indexName :=
		elasticsearch.GetProductIndex(
			tenantCode,
			countryCode,
		)

	logger.Log.Info(
		"index name",
		zap.String("index_name", indexName),
	)

	filter :=
		[]map[string]interface{}{}

	// =========================
	// CATEGORY
	// =========================

	if query.Category != "" {

		logger.Log.Info("category", zap.String("category", query.Category))

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

		logger.Log.Info("subcategory", zap.String("subcategory", query.SubCategory))

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

		logger.Log.Info("brand", zap.String("brand", query.Brand))

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

		logger.Log.Info(
			"price range",
			zap.Float64("min_price", query.MinPrice),
			zap.Float64("max_price", query.MaxPrice),
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

	logger.Log.Info("filter", zap.Any("filter", filter))

	// =========================
	// SORT
	// =========================

	sort :=
		buildSort(
			query.Sort,
		)

	// =========================
	// SEARCH QUERY
	// =========================

	searchQuery :=
		map[string]interface{}{
			"size": query.Limit,

			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": []map[string]interface{}{
						{
							"multi_match": map[string]interface{}{
								"query": query.Query,

								"fields": []string{
									"name^5",
									"brand^3",
									"category^2",
									"subcategory^2",
								},

								"fuzziness": "AUTO",
							},
						},
					},

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

		logger.Log.Info("cursor", zap.String("cursor", query.Cursor))

		searchQuery["search_after"] =
			[]interface{}{
				query.Cursor,
			}
	}

	logger.Log.Info("search query", zap.Any("search_query", searchQuery))

	body, err :=
		json.Marshal(searchQuery)

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

	if err != nil {

		logger.Log.Error(
			"failed elastic search",

			zap.Error(err),
		)

		return nil, nil, 0, err
	}

	defer res.Body.Close()

	var result map[string]interface{}

	err =
		json.NewDecoder(
			res.Body,
		).Decode(
			&result,
		)

	if err != nil {

		logger.Log.Error(
			"failed to unmarshal elastic search response",
			zap.Error(err),
		)

		return nil, nil, 0, err
	}

	// =========================
	// HITS
	// =========================

	hits :=
		result["hits"].(map[string]interface{})

	total :=
		int64(
			hits["total"].(map[string]interface{})["value"].(float64),
		)

	rows :=
		hits["hits"].([]interface{})

	var products []map[string]interface{}

	for _, row := range rows {

		rowMap :=
			row.(map[string]interface{})

		source :=
			rowMap["_source"].(map[string]interface{})

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

	aggs :=
		result["aggregations"].(map[string]interface{})

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
			{
				"id": map[string]interface{}{
					"order": "desc",
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
			{
				"id": map[string]interface{}{
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
				"_score": map[string]interface{}{
					"order": "desc",
				},
			},
			{
				"id": map[string]interface{}{
					"order": "desc",
				},
			},
		}
	}
}
