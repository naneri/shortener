package dto

// ListLink the DTO is used to display the list of user shortened URLs
type ListLink struct {
	// ShortURL - the shortened URL
	ShortURL string `json:"short_url"`
	// OriginalURL - the original URL
	OriginalURL string `json:"original_url"`
}
