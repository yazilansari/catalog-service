package main

import (
	"catalog-service/internal/database"
	"catalog-service/internal/elasticsearch"
	"catalog-service/internal/logger"
	"catalog-service/internal/plp/repository"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	esutil "github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {

	// Load ENV

	err := godotenv.Load()

	if err != nil {
		logger.Log.Fatal(".env file not loaded")
	}

	// =========================
	// INIT
	// =========================

	logger.InitLogger()

	defer logger.Log.Sync()

	database.ConnectPostgres()

	logger.Log.Info("postgres connected")

	elasticsearch.InitElasticSearch()

	logger.Log.Info("elasticsearch connected")

	// =========================
	// TENANT / COUNTRY
	// =========================

	tenantCode := "AE"

	countryCode := "AE"

	indexName :=
		elasticsearch.GetProductIndex(
			tenantCode,
			countryCode,
		)

	logger.Log.Info(
		"starting elasticsearch seed",

		zap.String(
			"tenant_code",
			tenantCode,
		),

		zap.String(
			"country_code",
			countryCode,
		),

		zap.String(
			"index_name",
			indexName,
		),
	)

	// =========================
	// CREATE INDEX
	// =========================

	err =
		elasticsearch.CreateProductIndex(
			tenantCode,
			countryCode,
		)

	if err != nil {

		logger.Log.Fatal(
			"failed to create index",

			zap.Error(err),
		)
	}

	logger.Log.Info(
		"index ready",

		zap.String(
			"index_name",
			indexName,
		),
	)

	logger.Log.Info(
		"starting elasticsearch seed...",
	)

	// =========================
	// FETCH PRODUCTS
	// =========================

	products,
		err :=
		repository.GetProductsForElastic(
			tenantCode,
			countryCode,
		)

	if err != nil {

		logger.Log.Error(
			"failed to fetch products",
		)
		panic(err)
	}

	logger.Log.Info(
		"products found:",
		zap.Int("count", len(products)),
	)

	// =========================
	// BULK INDEXER
	// =========================

	indexer, err :=
		esutil.NewBulkIndexer(
			esutil.BulkIndexerConfig{
				Client: elasticsearch.Client,

				Index: indexName,

				NumWorkers: 4,

				FlushBytes: 5e+6,
			},
		)

	if err != nil {

		logger.Log.Error(
			"failed to create bulk indexer",
		)
		panic(err)
	}

	// =========================
	// LOOP PRODUCTS
	// =========================

	for _, product := range products {

		document :=
			map[string]interface{}{
				"id": product.ID,

				"name": product.Name,

				"slug": product.Slug,

				"price": product.Price,

				"category": product.Category,

				"subcategory": product.SubCategory,

				"brand": product.Brand,

				"discount_price": product.DiscountPrice,

				"status": product.Status,

				"created_at": product.CreatedAt,

				"sales_count": 0,
			}

		body, err :=
			json.Marshal(
				document,
			)

		if err != nil {

			logger.Log.Error(
				"failed to marshal product",

				zap.Uint64(
					"product_id",
					product.ID,
				),

				zap.Error(err),
			)

			continue
		}

		err =
			indexer.Add(
				context.Background(),

				esutil.BulkIndexerItem{
					Action: "index",

					DocumentID: fmt.Sprintf(
						"%d",
						product.ID,
					),

					Body: strings.NewReader(
						string(body),
					),

					OnSuccess: func(
						ctx context.Context,
						item esutil.BulkIndexerItem,
						resp esutil.BulkIndexerResponseItem,
					) {

						logger.Log.Debug(
							"product indexed",

							zap.String(
								"document_id",
								item.DocumentID,
							),
						)
					},

					OnFailure: func(
						ctx context.Context,
						item esutil.BulkIndexerItem,
						resp esutil.BulkIndexerResponseItem,
						err error,
					) {

						logger.Log.Error(
							"failed to index product",

							zap.String(
								"document_id",
								item.DocumentID,
							),

							zap.Error(err),
						)
					},
				},
			)

		if err != nil {

			logger.Log.Error(
				"failed to add product to bulk indexer",

				zap.Uint64(
					"product_id",
					product.ID,
				),

				zap.Error(err),
			)
		}
	}

	// =========================
	// CLOSE INDEXER
	// =========================

	err =
		indexer.Close(
			context.Background(),
		)

	if err != nil {

		logger.Log.Fatal(
			"failed to close bulk indexer",

			zap.Error(err),
		)
	}

	stats :=
		indexer.Stats()

	logger.Log.Info(
		"elasticsearch seed completed",

		zap.Uint64(
			"indexed_count",
			stats.NumIndexed,
		),

		zap.Uint64(
			"failed_count",
			stats.NumFailed,
		),
	)
}
