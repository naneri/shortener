package dto

// BatchLink is a link structure to which the requests are unmarshalled when processing User links in batch
type BatchLink struct {
	// CorrelationID - the users ID used to distinguish the links
	CorrelationID string `json:"correlation_id"`
	// OriginalURL - the URL to shorten
	OriginalURL string `json:"original_url"`
}

// ResponseBatchLink is a structure to which a response is shortened
type ResponseBatchLink struct {
	// CorrelationID - the users ID that the user has passed initially with the URL
	CorrelationID string `json:"correlation_id"`
	// ShortURL - the result of the URL shortening
	ShortURL string `json:"short_url"`
}
