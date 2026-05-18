package database

import (
	"fmt"
	"os"
	"time"

	"catalog-service/internal/logger"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectPostgres() {

	logger.Log.Info(
		"connecting to postgres database",
	)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	start := time.Now()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	duration := time.Since(start)

	if err != nil {
		logger.Log.Fatal(
			"failed to connect postgres database",

			zap.Error(err),

			zap.Duration(
				"connection_duration",
				duration,
			),
		)

		return
	}

	sqlDB, err := db.DB()

	if err != nil {

		logger.Log.Fatal(
			"failed to get sql db instance",

			zap.Error(err),
		)

		return
	}

	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db

	logger.Log.Info(
		"postgres database connected successfully",

		zap.Duration(
			"connection_duration",
			duration,
		),

		zap.Int(
			"max_open_connections",
			100,
		),

		zap.Int(
			"max_idle_connections",
			20,
		),
	)
}
