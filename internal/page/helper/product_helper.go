package helper

import (
	"catalog-service/internal/page/model"
)

func ExtractProductIDs(
	products []model.Product,
) []uint64 {

	ids := make(
		[]uint64,
		0,
		len(products),
	)

	for _, product := range products {

		ids = append(
			ids,
			product.ID,
		)
	}

	return ids
}

func BuildPriceMap(
	products []model.Product,
) map[uint64]float64 {

	result := make(
		map[uint64]float64,
	)

	for _, product := range products {

		result[product.ID] = product.Price
	}

	return result
}
