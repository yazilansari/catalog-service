package model

type SEO struct {
	ID uint64 `gorm:"column:id" json:"id"`

	PageType string `gorm:"column:entity_type" json:"pageType"`

	Slug string `gorm:"column:slug" json:"slug"`

	MetaTitle string `gorm:"column:title" json:"metaTitle"`

	MetaDescription string `gorm:"column:meta_description" json:"metaDescription"`

	MetaKeywords string `gorm:"column:meta_keywords" json:"metaKeywords"`
}

func (SEO) TableName() string {
	return "seo_pages"
}
