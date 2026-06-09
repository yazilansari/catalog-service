package dto

import "catalog-service/internal/page/model"

type ProductImage struct {
	ID        uint64 `json:"id"`
	Image     string `json:"image"`
	SortOrder int    `json:"SortOrder"`
}

type ProductPageResponse struct {
	PageType string `json:"pageType"`

	Category model.Category `json:"category"`

	SubCategory model.Category `json:"subCategory"`

	Images []ProductImage `json:"images"`

	Product model.Product `json:"product"`

	RelatedProducts []model.Product `json:"relatedProducts"`
}
