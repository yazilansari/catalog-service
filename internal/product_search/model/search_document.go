package model

import promotionDto "catalog-service/internal/promotion/dto"

type ProductSearchDocument struct {
	ID uint64 `json:"id"`

	Name string `json:"name"`

	Slug string `json:"slug"`

	Category string `json:"category"`

	SubCategory string `json:"subCategory"`

	Brand string `json:"brand"`

	Price float64 `json:"price"`

	// DiscountPrice float64 `json:"discount_price"`

	// Status string `json:"status"`

	// CreatedAt string `json:"created_at"`

	Cursor interface{} `json:"cursor,omitempty"`

	Promotion *promotionDto.PromotionResponse `json:"promotion,omitempty"`
}
