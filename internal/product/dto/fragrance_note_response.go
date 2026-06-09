package dto

type FragranceNoteResponse struct {
	TopNote string `json:"topNote"`

	HeartNote string `json:"heartNote"`

	BaseNote string `json:"baseNote"`

	TopNoteImage string `json:"topNoteImage"`

	HeartNoteImage string `json:"heartNoteImage"`

	BaseNoteImage string `json:"baseNoteImage"`

	TopNoteDescription string `json:"topNoteDescription"`

	HeartNoteDescription string `json:"heartNoteDescription"`

	BaseNoteDescription string `json:"baseNoteDescription"`
}
