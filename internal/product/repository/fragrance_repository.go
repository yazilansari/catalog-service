package repository

import (
	"catalog-service/internal/database"
	"catalog-service/internal/logger"
	"catalog-service/internal/product/model"
	"time"

	"go.uber.org/zap"
)

func GetFragranceNotes(
	productID uint64,
) ([]model.FragranceNote, error) {

	logger.Log.Info(
		"fetching fragrance notes",

		zap.Uint64("product_id", productID),
	)

	var notes []model.FragranceNote

	start := time.Now()

	query := database.DB.
		Table("product_fragrance_notes fn").
		Select("fn.*").
		Joins(`
			INNER JOIN product_fragrance_map pfm
			ON pfm.fragrance_note_id = fn.id
		`).
		Where("pfm.product_id = ?", productID).
		Where("fn.status = ?", "published")

	err := query.Find(&notes).Error

	duration := time.Since(start)

	if err != nil {

		logger.Log.Error(
			"failed to fetch fragrance notes",

			zap.Error(err),
			zap.Duration("duration", duration),
		)
		return nil, err
	}

	logger.Log.Info(
		"fragrance notes fetched successfully",

		zap.Uint64("product_id", productID),
		zap.Duration("duration", duration),
	)

	if duration > time.Second {

		logger.Log.Warn(
			"slow fragrance notes query detected",

			zap.Uint64(
				"product_id",
				productID,
			),

			zap.Duration(
				"duration",
				duration,
			),
		)
	}

	return notes, nil
}
