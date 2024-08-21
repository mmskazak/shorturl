package inmemory

import (
	"reflect"
	"sync"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestNewInMemory(t *testing.T) {
	type args struct {
		zapLog *zap.SugaredLogger
	}

	zaplogSugar := zaptest.NewLogger(t).Sugar()

	tests := []struct {
		name    string
		args    args
		want    *InMemory
		wantErr bool
	}{
		{
			name: "test 1",
			args: args{
				zapLog: zaplogSugar,
			},
			want: &InMemory{
				mu:        &sync.Mutex{},
				data:      make(map[string]URLRecord),
				userIndex: make(map[string][]string),
				zapLog:    zaplogSugar,
			},
			wantErr: false,
		},
		{
			name: "test 2",
			args: args{
				zapLog: nil,
			},
			want: &InMemory{
				mu:        &sync.Mutex{},
				data:      make(map[string]URLRecord),
				userIndex: make(map[string][]string),
				zapLog:    nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewInMemory(tt.args.zapLog)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewInMemory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInMemory() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInMemory_GetCopyData(t *testing.T) {
	type fields struct {
		mu        *sync.Mutex
		data      map[string]URLRecord
		userIndex map[string][]string
		zapLog    *zap.SugaredLogger
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]URLRecord
	}{
		{
			name: "empty data",
			fields: fields{
				mu:        &sync.Mutex{},
				data:      make(map[string]URLRecord),
				userIndex: make(map[string][]string),
				zapLog:    zap.NewNop().Sugar(), // Используем no-op логгер для тестирования
			},
			want: make(map[string]URLRecord),
		},
		{
			name: "non-empty data",
			fields: fields{
				mu: &sync.Mutex{},
				data: map[string]URLRecord{
					"1": {ShortURL: "short1", OriginalURL: "original1", UserID: "user1", Deleted: false},
					"2": {ShortURL: "short2", OriginalURL: "original2", UserID: "user2", Deleted: true},
				},
				userIndex: make(map[string][]string),
				zapLog:    zap.NewNop().Sugar(), // Используем no-op логгер для тестирования
			},
			want: map[string]URLRecord{
				"1": {ShortURL: "short1", OriginalURL: "original1", UserID: "user1", Deleted: false},
				"2": {ShortURL: "short2", OriginalURL: "original2", UserID: "user2", Deleted: true},
			},
		},
		// Добавьте дополнительные тестовые случаи, если это необходимо
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &InMemory{
				mu:        tt.fields.mu,
				data:      tt.fields.data,
				userIndex: tt.fields.userIndex,
				zapLog:    tt.fields.zapLog,
			}
			got := m.GetCopyData()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCopyData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInMemory_NumberOfEntries(t *testing.T) {
	type fields struct {
		mu        *sync.Mutex
		data      map[string]URLRecord
		userIndex map[string][]string
		zapLog    *zap.SugaredLogger
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "empty data",
			fields: fields{
				mu:        &sync.Mutex{},
				data:      make(map[string]URLRecord),
				userIndex: make(map[string][]string),
				zapLog:    zap.NewNop().Sugar(), // Используем no-op логгер для тестирования
			},
			want: 0,
		},
		{
			name: "non-empty data",
			fields: fields{
				mu: &sync.Mutex{},
				data: map[string]URLRecord{
					"1": {ShortURL: "short1", OriginalURL: "original1", UserID: "user1", Deleted: false},
					"2": {ShortURL: "short2", OriginalURL: "original2", UserID: "user2", Deleted: true},
					"3": {ShortURL: "short3", OriginalURL: "original3", UserID: "user3", Deleted: false},
				},
				userIndex: make(map[string][]string),
				zapLog:    zap.NewNop().Sugar(), // Используем no-op логгер для тестирования
			},
			want: 3,
		},
		// Добавьте дополнительные тестовые случаи, если это необходимо
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &InMemory{
				mu:        tt.fields.mu,
				data:      tt.fields.data,
				userIndex: tt.fields.userIndex,
				zapLog:    tt.fields.zapLog,
			}
			got := m.NumberOfEntries()
			if got != tt.want {
				t.Errorf("NumberOfEntries() = %v, want %v", got, tt.want)
			}
		})
	}
}
