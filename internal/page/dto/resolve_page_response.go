package dto

type ResolvePageResponse struct {
	PageType string `json:"page_type"`

	Slug string `json:"slug"`

	RedirectURL string `json:"redirect_url"`
}
