package dto

import "catalog-service/internal/page/model"

type ProductPageResponse struct {
	PageType string `json:"page_type"`

	Category model.Category `json:"category"`

	SubCategory model.Category `json:"subcategory"`

	Product model.Product `json:"product"`

	RelatedProducts []model.Product `json:"related_products"`
}
