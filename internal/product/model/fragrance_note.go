package model

type FragranceNote struct {
	ID uint64 `gorm:"column:id"`

	ItemFamily string `gorm:"column:itemFamily"`

	TopNote   string `gorm:"column:top_note"`
	HeartNote string `gorm:"column:heart_note"`
	BaseNote  string `gorm:"column:base_note"`

	TopNoteImage   string `gorm:"column:top_note_image"`
	HeartNoteImage string `gorm:"column:heart_note_image"`
	BaseNoteImage  string `gorm:"column:base_note_image"`

	TopNoteDescription   string `gorm:"column:top_note_description"`
	HeartNoteDescription string `gorm:"column:heart_note_description"`
	BaseNoteDescription  string `gorm:"column:base_note_description"`

	Status string `gorm:"column:status"`
}

func (FragranceNote) TableName() string {
	return "product_fragrance_notes"
}
