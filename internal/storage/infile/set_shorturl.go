package infile

import (
	"context"
)

// SetShortURL сохраняет короткий URL и ассоциированные с ним данные в хранилище.
//
// Функция выполняет следующие шаги:
// 1. Делегирует сохранение короткого URL в внутреннее хранилище InMemory.
// 2. Если сохранение успешно, инициирует асинхронное сохранение данных в файл.
// 3. Логирует успешное добавление короткой ссылки, включая её идентификатор, оригинальный URL и идентификатор пользователя.
//
// Возможные ошибки:
//   - Если внутреннее хранилище возвращает ошибку, она пробрасывается дальше.
//     Возможные ошибки включают:
//   - ErrKeyAlreadyExists — если ключ уже существует в хранилище.
//   - ConflictError (с ErrOriginalURLAlreadyExists) — если оригинальный URL уже существует.
//
// Примечания:
// - Ошибки обрабатываются и пробрасываются дальше без обёртывания, чтобы сохранять оригинальное сообщение об ошибке.
// - После успешного сохранения вызывается функция `saveToFile` для записи данных в файл асинхронно.
func (m *InFile) SetShortURL(
	ctx context.Context,
	idShortPath string,
	originalURL string,
	userID string,
	deleted bool,
) error {
	err := m.InMe.SetShortURL(ctx, idShortPath, originalURL, userID, deleted)
	if err != nil {
		return err //nolint:wrapcheck // пробрасываем дальше оригинальную ошибку
	}
	m.saveToFile()
	m.zapLog.Infof("Added short link id: %s, URL: %s, for UserID %s", idShortPath, originalURL, userID)
	return nil
}
