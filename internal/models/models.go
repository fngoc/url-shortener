package models

type (
	Request struct {
		URL string `json:"url"`
	}

	Response struct {
		Result string `json:"result"`
	}

	RequestBatch struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	ResponseBatch struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}

	URLData struct {
		UUID        int    `json:"uuid"`
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}
)
