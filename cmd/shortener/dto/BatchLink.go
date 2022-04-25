package dto

type BatchLink struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ResponseBatchLink struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
