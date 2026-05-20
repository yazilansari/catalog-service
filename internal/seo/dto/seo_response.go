package dto

type SEOResponse struct {
	Title string `json:"title"`

	MetaDescription string `json:"meta_description"`

	MetaKeywords string `json:"meta_keywords"`

	CanonicalURL string `json:"canonical_url"`

	OGTitle string `json:"og_title"`

	OGDescription string `json:"og_description"`

	OGImage string `json:"og_image"`

	Robots string `json:"robots"`

	SchemaJSON string `json:"schema_json"`
}
