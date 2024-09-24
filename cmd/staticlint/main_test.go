package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig_ValidFile(t *testing.T) {
	// Создаем временный файл конфигурации
	file, err := os.CreateTemp("", "_staticlint.json")
	require.NoError(t, err)
	defer func(name string) {
		err := os.Remove(name)
		require.NoError(t, err)
	}(file.Name()) // Удаляем файл после теста

	content := `{"checks":{"S1001": true,"S1002": true}}`
	_, err = file.WriteString(content)
	require.NoError(t, err)

	// Закрываем файл после записи
	require.NoError(t, file.Close())

	// Считываем конфигурацию из файла
	actual, err := loadConfig(file.Name())
	require.NoError(t, err)

	expected := map[string]bool{"S1001": true, "S1002": true}
	require.Equal(t, expected, actual)
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	// Создаем временный файл конфигурации
	file, err := os.CreateTemp("", "_staticlint.json")
	require.NoError(t, err)
	defer func(name string) {
		err := os.Remove(name)
		require.NoError(t, err)
	}(file.Name()) // Удаляем файл после теста

	content := `{"checks":{"S1001": true,"S1002": true` // Некорректный JSON
	_, err = file.WriteString(content)
	require.NoError(t, err)

	// Закрываем файл после записи
	require.NoError(t, file.Close())

	// Считываем конфигурацию из файла
	actual, err := loadConfig(file.Name())
	require.Error(t, err)
	require.Nil(t, actual)
}

func TestLoadConfig_EmptyFile(t *testing.T) {
	// Создаем временный файл конфигурации
	file, err := os.CreateTemp("", "config*.json")
	require.NoError(t, err)
	defer func(name string) {
		err := os.Remove(name)
		require.NoError(t, err)
	}(file.Name()) // Удаляем файл после теста

	// Создаем пустой файл
	_, err = file.WriteString("")
	require.NoError(t, err)

	// Закрываем файл после записи
	require.NoError(t, file.Close())

	// Считываем конфигурацию из файла
	actual, err := loadConfig(file.Name())
	require.Error(t, err)
	require.Nil(t, actual)
}

func TestLoadConfig_FileDoesNotExist(t *testing.T) {
	// Путь к несуществующему файлу
	filePath := "nonexistent_file.json"

	// Считываем конфигурацию из файла
	actual, err := loadConfig(filePath)
	require.Error(t, err)
	require.Nil(t, actual)
}

func Test_getAnalyzers(t *testing.T) {
	got := getAnalyzers()
	assert.NotNil(t, got)
}
