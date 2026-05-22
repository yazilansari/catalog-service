package dto

import "catalog-service/internal/page/model"

type CategoryPageResponse struct {
	PageType string `json:"page_type"`

	Category model.Category `json:"category"`

	SubCategories []model.Category `json:"subcategories"`

	Products []model.Product `json:"products"`
}
