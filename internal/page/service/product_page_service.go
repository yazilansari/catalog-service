package service

import (
	"catalog-service/internal/logger"
	"catalog-service/internal/page/repository"

	"go.uber.org/zap"
)

func GetProductPage(
	tenantCode string,
	countryCode string,
	slug string,
) (interface{}, error) {

	logger.Log.Info(
		"product page request",

		zap.String(
			"slug",
			slug,
		),
	)

	product, err :=
		repository.FindProductBySlug(
			tenantCode,
			countryCode,
			slug,
		)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"type":    "product",
		"product": product,
	}, nil
}
