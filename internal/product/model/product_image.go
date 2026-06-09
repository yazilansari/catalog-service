package model

type ProductImage struct {
	ID uint64 `gorm:"column:id" json:"id"`

	// ProductID uint64 `gorm:"column:product_id"`

	Image string `gorm:"column:image_url" json:"imageUrl"`

	SortOrder int `gorm:"column:sort_order" json:"sortOrder"`
}

func (ProductImage) TableName() string {
	return "product_images"
}
