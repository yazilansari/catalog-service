package model

import "time"

type SEOPage struct {
	ID uint64 `gorm:"primaryKey"`

	TenantCode  string `gorm:"index"`
	CountryCode string `gorm:"index"`

	EntityType string `gorm:"index"`
	// product/category/blog

	EntityID uint64 `gorm:"index"`

	Slug string `gorm:"uniqueIndex"`

	Title string

	MetaDescription string
	MetaKeywords    string

	CanonicalURL string

	OGTitle       string
	OGDescription string
	OGImage       string

	Robots string

	SchemaJSON string `gorm:"type:jsonb"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (SEOPage) TableName() string {
	return "seo_pages"
}
