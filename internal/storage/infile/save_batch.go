package infile

import (
	"context"

	"mmskazak/shorturl/internal/storage"
)

// SaveBatch сохраняет пакет входящих данных в хранилище и обновляет файл хранения.
//
// Аргументы:
//   - ctx: Контекст выполнения, используемый для управления временем жизни запроса и отмены.
//   - items: Срез объектов, представляющих входящие данные для сохранения.
//   - baseHost: Базовый URL-адрес для формирования полного короткого URL.
//   - userID: Идентификатор пользователя, связанный с URL.
//   - generator: Интерфейс для генерации уникальных идентификаторов коротких URL.
//
// Возвращает:
//   - []storage.Output: Срез объектов, представляющих результат сохранения данных.
//   - error: Ошибка, если она произошла при сохранении данных или обновлении файла.
//
// Примечание:
// Функция передает сохранение данных в хранилище в памяти (InMemory), затем сохраняет обновленные данные в файл.
// Если при сохранении данных в памяти возникает ошибка, она возвращается без оборачивания.
func (m *InFile) SaveBatch(
	ctx context.Context,
	items []storage.Incoming,
	baseHost string,
	userID string,
	generator storage.IGenIDForURL,
) ([]storage.Output, error) {
	outputs, err := m.InMe.SaveBatch(ctx, items, baseHost, userID, generator)
	if err != nil {
		return nil, err //nolint:wrapcheck // прокидываем оригинальную ошибку
	}
	m.saveToFile()
	return outputs, nil
}
