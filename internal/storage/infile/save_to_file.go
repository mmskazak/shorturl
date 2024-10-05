package infile

import (
	"encoding/json"
	"os"
	"strconv"
)

// saveToFile асинхронно сохраняет все данные из памяти в файл.
//
// Функция запускает горутину, которая:
// 1. Логирует начало процесса сохранения данных в файл.
// 2. Получает копию данных из памяти и формирует список записей для сохранения.
// 3. Открывает файл для перезаписи, создавая его, если он не существует.
// 4. Сериализует каждую запись в формате JSON и записывает ее в файл, добавляя перенос строки после каждой записи.
// 5. Логирует успешное завершение сохранения данных.
//
// Примечания:
// - Данные сохраняются в файл асинхронно для улучшения производительности и предотвращения блокировки основной работы.
// - Ошибки при создании файла или записи данных логируются, но не останавливают выполнение.
// - Файл закрывается после завершения записи, и возможные ошибки при закрытии также логируются.
func (f *InFile) saveToFile() {
	f.zapLog.Infoln("Save urls to file")
	go func() {
		f.zapLog.Info("Starting to save storage data to file")

		// Получение копии данных из памяти
		data := f.InMe.GetCopyData()

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
		file, err := os.Create(f.filePath)
		if err != nil {
			f.zapLog.Errorf("error creating file: %v", err)
			return
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				f.zapLog.Warnf("error closing file storage: %v", err)
			}
		}(file)

		// Сериализуем данные в JSON с переносом строки в конце каждой записи
		encoder := json.NewEncoder(file)
		for _, record := range records {
			err := encoder.Encode(record)
			if err != nil {
				f.zapLog.Errorf("error encoding data: %v", err)
				return
			}
		}

		f.zapLog.Info("Successfully saved storage data to file")
	}()
}
