package models

// Incoming представляет данные, полученные при создании или обновлении короткого URL.
// CorrelationID - идентификатор запроса, связанный с этим URL.
// OriginalURL - оригинальный (длинный) URL, который нужно сократить.
type Incoming struct {
	CorrelationID string `json:"correlation_id"` // строковый идентификатор
	OriginalURL   string `json:"original_url"`   // оригинальный URL
}

// Output представляет результат операции создания или обновления короткого URL.
// CorrelationID - идентификатор запроса, связанный с этим URL.
// ShortURL - сгенерированный короткий URL.
type Output struct {
	CorrelationID string `json:"correlation_id"` // строковый идентификатор
	ShortURL      string `json:"short_url"`      // короткий URL
}

// URL представляет собой короткий URL и его оригинальный (длинный) URL.
// ShortURL - сгенерированный короткий URL.
// OriginalURL - оригинальный URL.
type URL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// JSONRequest представляет структуру запроса для создания короткого URL.
type JSONRequest struct {
	URL string `json:"url"`
}

// JSONResponse представляет структуру ответа с созданным коротким URL.
type JSONResponse struct {
	ShortURL string `json:"result"`
}

// URLRecord - структура для хранения URL с дополнительной информацией.
type URLRecord struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"` // Оригинальный URL
	UserID      string `json:"user_id"`      // Идентификатор пользователя
	Deleted     bool   `json:"deleted"`      // Флаг, указывающий на удаление URL
}

// Stats - структура содержит статистику]
// "urls": <int> количество сокращённых URL в сервисе
// "users": <int>  количество пользователей в сервисе
type Stats struct {
	Urls  string `json:"urls"`
	Users string `json:"users"`
}
