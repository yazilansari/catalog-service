package model

type ProductVariant struct {
	ID uint64 `gorm:"column:id"`

	ProductID uint64 `gorm:"column:product_id"`

	Name string `gorm:"column:name"`

	SKU string `gorm:"column:sku"`

	Price float64 `gorm:"column:price"`

	DiscountPrice float64 `gorm:"column:discount_price"`

	Stock int `gorm:"column:stock"`

	Status string `gorm:"column:status"`
}

func (ProductVariant) TableName() string {
	return "product_variants"
}
