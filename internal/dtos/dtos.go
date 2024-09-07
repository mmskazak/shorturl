package dtos

// DTOShortURL представляет данные, необходимые для создания короткого URL.
type DTOShortURL struct {
	OriginalURL string // Оригинальный URL
	UserID      string // Идентификатор пользователя
	BaseHost    string // Базовый хост для формирования короткого URL
	Deleted     bool   // Флаг, указывающий, удален ли URL
}
