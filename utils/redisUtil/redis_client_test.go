package redisUtil

import (
	"github.com/alicebob/miniredis/v2"
	"testing"
)

func TestNewRedisClient(t *testing.T) {
	type args struct {
		conf Config
	}
	mr := miniredis.RunT(t)
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "invalid-addr",
			args: args{
				conf: Config{
					Addr:     "invalid-addr",
					Password: "",
					Db:       0,
				},
			},
			wantErr: true,
		},
		{
			name: "ok",
			args: args{
				conf: Config{
					Addr:     mr.Addr(),
					Password: "",
					Db:       0,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewRedisClient(tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRedisClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
