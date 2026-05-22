package service

import (
	"catalog-service/internal/logger"
	"catalog-service/internal/page/repository"

	"go.uber.org/zap"
)

func GetSubCategoryPage(
	tenantCode string,
	countryCode string,
	slug string,
) (interface{}, error) {

	logger.Log.Info(
		"subcategory page request",

		zap.String(
			"slug",
			slug,
		),
	)

	subcategory, err :=
		repository.FindSubCategoryBySlug(
			tenantCode,
			countryCode,
			slug,
		)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"type":        "subcategory",
		"subcategory": subcategory,
	}, nil
}
