package model

type SEO struct {
	ID uint64 `gorm:"column:id"`

	PageType string `gorm:"column:page_type"`

	Slug string `gorm:"column:slug"`

	MetaTitle string `gorm:"column:meta_title"`

	MetaDescription string `gorm:"column:meta_description"`

	MetaKeywords string `gorm:"column:meta_keywords"`
}

func (SEO) TableName() string {
	return "seo_pages"
}
