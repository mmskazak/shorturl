package inmemory

import (
	"context"
	"sync"
	"testing"
)

func TestInMemory_SetShortURL(t *testing.T) {
	type fields struct {
		mu           *sync.Mutex
		data         map[string]string
		indexForData map[string]string
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
			name: "first success",
			fields: fields{
				mu:           &sync.Mutex{},
				data:         map[string]string{},
				indexForData: map[string]string{},
			},
			args: args{
				id:        "test0001",
				targetURL: "https://www.google.com",
			},
			wantErr: false,
		},
		{
			name: "second success",
			fields: fields{
				mu: &sync.Mutex{},
				data: func() map[string]string {
					return map[string]string{
						"test0001": "https://www.google.com",
					}
				}(),
				indexForData: func() map[string]string {
					return map[string]string{
						"https://www.google.com": "test0001",
					}
				}(),
			},
			args: args{
				id:        "test0001",
				targetURL: "https://www.google.com",
			},
			wantErr: true,
		},
		{
			name: "tree has error",
			fields: fields{
				mu:           &sync.Mutex{},
				data:         map[string]string{},
				indexForData: map[string]string{},
			},
			args: args{
				id:        "test0001",
				targetURL: "",
			},
			wantErr: true,
		},
		{
			name: "tree has error",
			fields: fields{
				mu:           &sync.Mutex{},
				data:         map[string]string{},
				indexForData: map[string]string{},
			},
			args: args{
				id:        "",
				targetURL: "https://www.google.com",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &InMemory{
				mu:           tt.fields.mu,
				data:         tt.fields.data,
				indexForData: tt.fields.indexForData,
			}
			ctx := context.TODO()
			if err := m.SetShortURL(ctx, tt.args.id, tt.args.targetURL); (err != nil) != tt.wantErr {
				t.Errorf("SetShortURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
