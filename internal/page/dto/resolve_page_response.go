package dto

type ResolvePageResponse struct {
	PageType string `json:"pageType"`

	Slug string `json:"slug"`

	RedirectURL string `json:"redirectUrl"`
}
