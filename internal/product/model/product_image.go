package model

type ProductImage struct {
	ID uint64 `gorm:"column:id"`

	// ProductID uint64 `gorm:"column:product_id"`

	Image string `gorm:"column:image_url"`

	SortOrder int `gorm:"column:sort_order"`
}

func (ProductImage) TableName() string {
	return "product_images"
}
