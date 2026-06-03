package dto

import (
	"catalog-service/internal/product_search/model"
)

type ProductSearchResponse struct {
	Products []model.ProductSearchDocument `json:"products"`

	Filters FilterResponse `json:"filters"`

	Pagination PaginationResponse `json:"pagination"`

	Sort string `json:"sort"`

	Total int64 `json:"total"`
}

type FilterResponse struct {
	Brands []string `json:"brands"`

	PriceRange PriceRange `json:"price_range"`
}

type PriceRange struct {
	Min float64 `json:"min"`

	Max float64 `json:"max"`
}

type PaginationResponse struct {
	NextCursor string `json:"next_cursor"`

	Limit int `json:"limit"`

	HasMore bool `json:"has_more"`
}
