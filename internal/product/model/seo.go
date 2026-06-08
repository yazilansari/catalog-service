package model

type SEO struct {
	ID uint64 `gorm:"column:id"`

	PageType string `gorm:"column:entity_type"`

	Slug string `gorm:"column:slug"`

	MetaTitle string `gorm:"column:title"`

	MetaDescription string `gorm:"column:meta_description"`

	MetaKeywords string `gorm:"column:meta_keywords"`
}

func (SEO) TableName() string {
	return "seo_pages"
}
