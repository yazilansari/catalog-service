package repository

import (
	"bytes"
	"catalog-service/internal/elasticsearch"
	"catalog-service/internal/plp/dto"
	"context"
	"encoding/json"
)

func SearchProducts(
	query dto.ProductQuery,
) (
	[]map[string]interface{},
	map[string]interface{},
	int64,
	error,
) {

	filter :=
		[]map[string]interface{}{}

	// =========================
	// CATEGORY
	// =========================

	if query.Category != "" {

		filter =
			append(
				filter,

				map[string]interface{}{
					"term": map[string]interface{}{
						"category.keyword": query.Category,
					},
				},
			)
	}

	// =========================
	// SUBCATEGORY
	// =========================

	if query.SubCategory != "" {

		filter =
			append(
				filter,

				map[string]interface{}{
					"term": map[string]interface{}{
						"subcategory.keyword": query.SubCategory,
					},
				},
			)
	}

	// =========================
	// BRAND
	// =========================

	if query.Brand != "" {

		filter =
			append(
				filter,

				map[string]interface{}{
					"term": map[string]interface{}{
						"brand.keyword": query.Brand,
					},
				},
			)
	}

	// =========================
	// PRICE RANGE
	// =========================

	if query.MinPrice > 0 ||
		query.MaxPrice > 0 {

		filter =
			append(
				filter,

				map[string]interface{}{
					"range": map[string]interface{}{
						"price": map[string]interface{}{
							"gte": query.MinPrice,

							"lte": query.MaxPrice,
						},
					},
				},
			)
	}

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
						"field": "brand.keyword",

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

		searchQuery["search_after"] =
			[]interface{}{
				query.Cursor,
			}
	}

	body, _ :=
		json.Marshal(
			searchQuery,
		)

	res, err :=
		elasticsearch.Client.Search(
			elasticsearch.Client.Search.WithContext(
				context.Background(),
			),

			elasticsearch.Client.Search.WithIndex(
				"products",
			),

			elasticsearch.Client.Search.WithBody(
				bytes.NewReader(body),
			),
		)

	if err != nil {
		return nil, nil, 0, err
	}

	defer res.Body.Close()

	var result map[string]interface{}

	_ = json.NewDecoder(
		res.Body,
	).Decode(
		&result,
	)

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

		source :=
			row.(map[string]interface{})["_source"]

		products =
			append(
				products,
				source.(map[string]interface{}),
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
