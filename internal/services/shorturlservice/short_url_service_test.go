package shorturlservice

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"mmskazak/shorturl/internal/services/shorturlservice/mocks"

	"github.com/golang/mock/gomock"
)

func TestShortURLService_GenerateShortURL_Success(t *testing.T) {
	// Создание нового контроллера
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Контекст и создание моков
	ctx := context.Background()
	generator := mocks.NewMockIGenIDForURL(ctrl)
	data := mocks.NewMockISetShortURL(ctrl)

	// Настройка полей и аргументов
	s := &ShortURLService{
		maxIteration: 10,
	}
	dto := DTOShortURL{
		OriginalURL: "ya.ru",
		UserID:      "1",
		BaseHost:    "http://localhost",
		Deleted:     false,
	}

	// Настройка ожиданий
	generator.EXPECT().Generate().Return("expectedID", nil)
	data.EXPECT().SetShortURL(ctx,
		"expectedID",
		dto.OriginalURL,
		dto.UserID,
		false).
		Return(nil)

	// Вызов тестируемого метода
	got, err := s.GenerateShortURL(ctx, dto, generator, data)

	// Проверка на ошибки
	if err != nil {
		t.Errorf("GenerateShortURL() error = %v, wantErr %v", err, false)
		return
	}

	// Проверка результата
	want := "http://localhost/expectedID"
	if got != want {
		t.Errorf("GenerateShortURL() got = %v, want %v", got, want)
	}
}

func TestShortURLService_GenerateShortURL_ErrGenerate(t *testing.T) {
	// Создание нового контроллера
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Контекст и создание моков
	ctx := context.Background()
	generator := mocks.NewMockIGenIDForURL(ctrl)
	data := mocks.NewMockISetShortURL(ctrl)

	// Настройка полей и аргументов
	s := &ShortURLService{
		maxIteration: 1,
	}
	dto := DTOShortURL{
		OriginalURL: "ya.ru",
		UserID:      "1",
		BaseHost:    "http://localhost",
		Deleted:     false,
	}

	// Настройка ожиданий
	generator.EXPECT().Generate().Return("", errors.New("test error generator"))
	// Вызов тестируемого метода
	_, err := s.GenerateShortURL(ctx, dto, generator, data)
	assert.Error(t, err)
}

func TestShortURLService_GenerateShortURL_ErrSetShortURL(t *testing.T) {
	// Создание нового контроллера
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Контекст и создание моков
	ctx := context.Background()
	generator := mocks.NewMockIGenIDForURL(ctrl)
	data := mocks.NewMockISetShortURL(ctrl)

	// Настройка полей и аргументов
	s := &ShortURLService{
		maxIteration: 2,
	}
	dto := DTOShortURL{
		OriginalURL: "ya.ru",
		UserID:      "1",
		BaseHost:    "http://localhost",
		Deleted:     false,
	}

	// Настройка ожиданий
	generator.EXPECT().Generate().Return("expectedID", nil).Times(2)
	data.EXPECT().SetShortURL(ctx,
		"expectedID",
		dto.OriginalURL,
		dto.UserID,
		false).
		Return(errors.New("test error set_short_url")).Times(2)

	// Вызов тестируемого метода
	_, err := s.GenerateShortURL(ctx, dto, generator, data)
	assert.Error(t, err)
}

func TestShortURLService_GenerateShortURL_ErrBaseHost(t *testing.T) {
	// Создание нового контроллера
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Контекст и создание моков
	ctx := context.Background()
	generator := mocks.NewMockIGenIDForURL(ctrl)
	data := mocks.NewMockISetShortURL(ctrl)

	// Настройка полей и аргументов
	s := &ShortURLService{
		maxIteration: 10,
	}
	dto := DTOShortURL{
		OriginalURL: "ya.ru",
		UserID:      "1",
		BaseHost:    "рtp://example.com",
		Deleted:     false,
	}

	// Вызов тестируемого метода
	_, err := s.GenerateShortURL(ctx, dto, generator, data)
	assert.Error(t, err)
}

// ConflictError должен реализовывать интерфейс error.
type ConflictError struct {
	Err      error
	ShortURL string
}

func (e ConflictError) Error() string {
	return e.Err.Error()
}

func TestShortURLService_GenerateShortURL_ErrConflictError(t *testing.T) {
	// Создание нового контроллера
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Контекст и создание моков
	ctx := context.Background()
	generator := mocks.NewMockIGenIDForURL(ctrl)
	data := mocks.NewMockISetShortURL(ctrl)

	// Настройка полей и аргументов
	s := &ShortURLService{
		maxIteration: 2,
	}
	dto := DTOShortURL{
		OriginalURL: "ya.ru",
		UserID:      "1",
		BaseHost:    "http://localhost",
		Deleted:     false,
	}

	conflictError := ConflictError{
		Err:      errors.New("test error set_short_url"),
		ShortURL: "expectedID",
	}

	// Настройка ожиданий
	generator.EXPECT().Generate().Return("expectedID", nil).Times(2)
	data.EXPECT().SetShortURL(ctx,
		"expectedID",
		dto.OriginalURL,
		dto.UserID,
		false).
		Return(conflictError).Times(2)

	// Вызов тестируемого метода
	_, err := s.GenerateShortURL(ctx, dto, generator, data)
	assert.Error(t, err)
}
