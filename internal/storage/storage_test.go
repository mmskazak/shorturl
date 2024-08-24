package storage

import "testing"

func TestGetFullShortURL(t *testing.T) {
	type args struct {
		baseHost string
		shortURL string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success test",
			args: args{
				baseHost: "https://localhost",
				shortURL: "shortURL",
			},
			want:    "https://localhost/shortURL",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFullShortURL(tt.args.baseHost, tt.args.shortURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFullShortURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetFullShortURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
