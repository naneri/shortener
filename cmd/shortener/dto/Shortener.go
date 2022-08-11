package dto

// ShortenerDto is used to parse the users passed URL
type ShortenerDto struct {
	// URL - the URL passed by the user to shorten
	URL string `json:"url"`
}
