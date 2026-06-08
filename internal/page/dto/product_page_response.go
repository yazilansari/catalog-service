package dto

import "catalog-service/internal/page/model"

type ProductImage struct {
	ID        uint64 `json:"ID"`
	Image     string `json:"Image"`
	SortOrder int    `json:"SortOrder"`
}

type ProductPageResponse struct {
	PageType string `json:"page_type"`

	Category model.Category `json:"category"`

	SubCategory model.Category `json:"subcategory"`

	Images []ProductImage `json:"images"`

	Product model.Product `json:"product"`

	RelatedProducts []model.Product `json:"related_products"`
}
