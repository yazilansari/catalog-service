package model

type Category struct {
	ID uint64 `gorm:"column:id"`

	Name string `gorm:"column:name"`

	Slug string `gorm:"column:slug"`

	ParentID uint64 `gorm:"column:parent_id" json:"-"`

	// Status string `gorm:"column:status"`
}

func (Category) TableName() string {
	return "categories"
}
