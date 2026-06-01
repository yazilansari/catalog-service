package dto

type FragranceNoteResponse struct {
	TopNote string `json:"top_note"`

	HeartNote string `json:"heart_note"`

	BaseNote string `json:"base_note"`

	TopNoteImage string `json:"top_note_image"`

	HeartNoteImage string `json:"heart_note_image"`

	BaseNoteImage string `json:"base_note_image"`

	TopNoteDescription string `json:"top_note_description"`

	HeartNoteDescription string `json:"heart_note_description"`

	BaseNoteDescription string `json:"base_note_description"`
}
