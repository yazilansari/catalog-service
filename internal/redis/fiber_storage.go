package redis

import (
	"os"
	"strconv"
	"time"

	"catalog-service/internal/logger"

	redisStorage "github.com/gofiber/storage/redis/v3"

	"go.uber.org/zap"
)

var FiberStorage *redisStorage.Storage

func InitFiberStorage() {

	logger.Log.Info(
		"initializing fiber redis storage",
	)

	port, err := strconv.Atoi(
		os.Getenv("REDIS_PORT"),
	)

	if err != nil {

		logger.Log.Fatal(
			"invalid redis port",

			zap.String(
				"redis_port",
				os.Getenv("REDIS_PORT"),
			),

			zap.Error(err),
		)

		return
	}

	db, err := strconv.Atoi(
		os.Getenv("REDIS_DB"),
	)

	if err != nil {

		logger.Log.Warn(
			"invalid redis db value, defaulting to 3",

			zap.String(
				"redis_db",
				os.Getenv("REDIS_DB"),
			),

			zap.Error(err),
		)

		db = 3
	}

	start := time.Now()

	FiberStorage = redisStorage.New(
		redisStorage.Config{
			Host: os.Getenv("REDIS_HOST"),

			Port: port,

			Password: os.Getenv("REDIS_PASSWORD"),

			Database: db,

			Reset: false,
		},
	)

	duration := time.Since(start)

	logger.Log.Info(
		"fiber redis storage initialized successfully",

		zap.String(
			"host",
			os.Getenv("REDIS_HOST"),
		),

		zap.Int(
			"port",
			port,
		),

		zap.Int(
			"database",
			db,
		),

		zap.Bool(
			"reset",
			false,
		),

		zap.Duration(
			"initialization_duration",
			duration,
		),
	)
}
