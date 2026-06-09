package dto

import "catalog-service/internal/page/model"

type CategoryPageResponse struct {
	PageType string `json:"pageType"`

	Category model.Category `json:"category"`

	SubCategories []model.Category `json:"subCategories"`

	Products []model.Product `json:"products"`
}
