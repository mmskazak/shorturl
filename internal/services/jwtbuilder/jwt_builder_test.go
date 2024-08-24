package jwtbuilder

import "testing"

func TestJWTBuilder_Create(t *testing.T) {
	type args struct {
		header  HeaderJWT
		payload PayloadJWT
		secret  string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				header: HeaderJWT{
					Alg: "HS256",
					Typ: "JWT",
				},
				payload: PayloadJWT{
					UserID: "exampleUserID",
				},
				secret: "secret",
			},
			want: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZXhhbXBsZVVzZXJJRCJ9." +
				"kfLHSHUEyv8izn-Rc38PicbbtO9EOFNkXUeEtW05jEA",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JWTBuilder{}
			got, err := j.Create(tt.args.header, tt.args.payload, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}
