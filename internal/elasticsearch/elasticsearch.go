package elasticsearch

import (
	"catalog-service/internal/logger"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"
)

var Client *elasticsearch.Client

func InitElasticSearch() {

	logger.Log.Info(
		"connecting to elasticsearch database",
	)

	start := time.Now()

	cfg :=
		elasticsearch.Config{
			Addresses: []string{
				os.Getenv(
					"ELASTICSEARCH_URL",
				),
			},
		}

	client, err :=
		elasticsearch.NewClient(
			cfg,
		)

	duration := time.Since(start)

	if err != nil {

		logger.Log.Fatal(
			"failed to connect elasticsearch database",

			zap.Error(err),

			zap.Duration(
				"connection_duration",
				duration,
			),
		)
		panic(err)
	}

	Client = client

	logger.Log.Info(
		"elasticsearch connected successfully",

		zap.Duration(
			"connection_duration",
			duration,
		),

		zap.String(
			"elasticsearch_url",
			os.Getenv("ELASTICSEARCH_URL"),
		),
	)
}
