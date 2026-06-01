package dto

import "catalog-service/internal/plp/model"

type ProductListResponse struct {
	Products []model.ProductDocument `json:"products"`

	Filters FilterResponse `json:"filters"`

	Pagination PaginationResponse `json:"pagination"`

	Sort string `json:"sort"`

	Total int64 `json:"total"`
}
