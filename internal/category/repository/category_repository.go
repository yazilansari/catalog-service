package repository

import (
	"log"

	"catalog-service/internal/category/model"
	"catalog-service/internal/database"
)

func GetCategories(
	tenantCode string,
	countryCode string,
) ([]model.Category, error) {

	log.Printf("[GetCategories] called tenantCode=%s countryCode=%s", tenantCode, countryCode)

	var categories []model.Category

	query := database.DB.
		Where("tenant_code = ?", tenantCode).
		Where("country_code = ?", countryCode).
		Where("status = ?", "published").
		Order("sort_order asc")

	log.Printf("[GetCategories] executing DB query")

	err := query.Find(&categories).Error

	if err != nil {
		log.Printf("[GetCategories] DB error: %v", err)
		return nil, err
	}

	log.Printf("[GetCategories] success count=%d", len(categories))

	return categories, nil
}
