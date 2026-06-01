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
	"go.uber.org/zap"
)

func main() {

	// =========================
	// INIT
	// =========================

	logger.InitLogger()

	database.ConnectPostgres()

	logger.Log.Info("postgres connected")

	elasticsearch.InitElasticSearch()

	logger.Log.Info("elasticsearch connected")

	logger.Log.Info(
		"starting elasticsearch seed...",
	)

	// =========================
	// FETCH PRODUCTS
	// =========================

	products,
		err :=
		repository.GetProductsForElastic()

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

				Index: "products",

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

				"discount_price": product.SalePrice,

				"status": product.Status,

				"created_at": product.CreatedAt,

				"sales_count": 0,
			}

		body, _ :=
			json.Marshal(
				document,
			)

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

					OnFailure: func(
						ctx context.Context,
						item esutil.BulkIndexerItem,
						resp esutil.BulkIndexerResponseItem,
						err error,
					) {

						fmt.Println(
							"failed:",
							item.DocumentID,
						)

						if err != nil {
							fmt.Println(err)
						}
					},
				},
			)

		if err != nil {

			logger.Log.Error(
				"failed to add product to bulk indexer",
			)
			panic(err)
		}
	}

	// =========================
	// CLOSE
	// =========================

	err =
		indexer.Close(
			context.Background(),
		)

	if err != nil {

		logger.Log.Error(
			"failed to close bulk indexer",
		)
		panic(err)
	}

	stats :=
		indexer.Stats()

	logger.Log.Info(
		"indexed:",
		zap.Uint64("count", stats.NumIndexed),
	)
}
