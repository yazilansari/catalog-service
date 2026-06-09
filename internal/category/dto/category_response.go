package dto

type CategoryResponse struct {
	ID                   uint64              `json:"id"`
	Name                 string              `json:"name"`
	Slug                 string              `json:"slug"`
	Image                *string             `json:"image,omitempty"`
	IconImage            *string             `json:"iconImage,omitempty"`
	MenuImage            *string             `json:"menuImage,omitempty"`
	MenuImage2           *string             `json:"menuImage2,omitempty"`
	MobileImage          *string             `json:"mobileImage,omitempty"`
	Video                *string             `json:"video"`
	ProductSubCategories []*CategoryResponse `json:"productSubCategories,omitempty"`
}
