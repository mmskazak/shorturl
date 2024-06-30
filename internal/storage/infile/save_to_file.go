package infile

import (
	"encoding/json"
	"os"
	"strconv"
)

// saveToFile - асинхронная запись всех данных из памяти в файл.
func (m *InFile) saveToFile() {
	go func() {
		m.zapLog.Info("Starting to save storage data to file")

		// Получение копии данных из памяти
		data := m.InMe.GetCopyData()

		// Формирование списка записей для сериализации
		records := make([]shortURLStruct, 0, len(data))
		numberItem := 1
		for k, v := range data {
			records = append(records, shortURLStruct{
				ID:          strconv.Itoa(numberItem),
				ShortURL:    k,
				OriginalURL: v.OriginalURL,
				UserID:      v.UserID,
				Deleted:     v.Deleted,
			})
			numberItem++
		}

		// Открываем файл для перезаписи
		file, err := os.Create(m.filePath)
		if err != nil {
			m.zapLog.Errorf("error creating file: %v", err)
			return
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				m.zapLog.Warnf("error closing file storage: %v", err)
			}
		}(file)

		// Сериализуем данные в JSON с переносом строки в конце каждой записи
		encoder := json.NewEncoder(file)
		for _, record := range records {
			err := encoder.Encode(record)
			if err != nil {
				m.zapLog.Errorf("error encoding data: %v", err)
				return
			}
		}

		m.zapLog.Info("Successfully saved storage data to file")
	}()
}
