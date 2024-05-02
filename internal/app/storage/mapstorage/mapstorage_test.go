package mapstorage //nolint: golint

import (
	"reflect"
	"sync"
	"testing"
)

func TestMapStorage_GetShortURL(t *testing.T) {
	type fields struct {
		mu   sync.Mutex
		data map[string]string
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Existing URL",
			fields: fields{
				data: map[string]string{
					"existingID": "https://example.com",
				},
			},
			args: args{
				id: "existingID",
			},
			want:    "https://example.com",
			wantErr: false,
		},
		{
			name: "Non-existing URL",
			fields: fields{
				data: map[string]string{
					"existingID": "https://example.com",
				},
			},
			args: args{
				id: "nonExistingID",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MapStorage{
				mu:   tt.fields.mu,
				data: tt.fields.data,
			}
			got, err := m.GetShortURL(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetShortURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetShortURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapStorage_SetShortURL(t *testing.T) {
	type fields struct {
		mu   sync.Mutex
		data map[string]string
	}
	type args struct {
		id        string
		targetURL string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "New URL",
			fields: fields{
				data: map[string]string{},
			},
			args: args{
				id:        "newID",
				targetURL: "https://example.com/new",
			},
			wantErr: false,
		},
		{
			name: "Existing URL",
			fields: fields{
				data: map[string]string{
					"existingID": "https://example.com/existing",
				},
			},
			args: args{
				id:        "existingID",
				targetURL: "https://example.com/updated",
			},
			wantErr: true,
		},
		{
			name: "Empty ID",
			fields: fields{
				data: map[string]string{},
			},
			args: args{
				id:        "",
				targetURL: "https://example.com/emptyid",
			},
			wantErr: true,
		},
		{
			name: "Empty URL",
			fields: fields{
				data: map[string]string{},
			},
			args: args{
				id:        "emptyURL",
				targetURL: "",
			},
			wantErr: true,
		},
	}
	for i := range tests {
		tt := &tests[i]
		t.Run(tt.name, func(t *testing.T) {
			m := &MapStorage{
				mu:   tt.fields.mu,
				data: tt.fields.data,
			}
			if err := m.SetShortURL(tt.args.id, tt.args.targetURL); (err != nil) != tt.wantErr {
				t.Errorf("SetShortURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewMapStorage(t *testing.T) {
	tests := []struct {
		name string
		want *MapStorage
	}{
		{
			name: "New instance",
			want: &MapStorage{data: make(map[string]string), mu: sync.Mutex{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMapStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMapStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}
