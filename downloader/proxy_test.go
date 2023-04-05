package downloader

import (
	"reflect"
	"testing"
)

func TestNewProxy(t *testing.T) {
	type args struct {
		id   int
		host string
		port int
		typ  string
	}
	tests := []struct {
		name    string
		args    args
		want    IProxy
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				id:   1,
				host: "127.0.0.1",
				port: 8080,
				typ:  "http",
			},
			want:    &Proxy{Id: 1, Host: "127.0.0.1", Port: 8080, Type: "http"},
			wantErr: false,
		},
		{
			name: "invalid type",
			args: args{
				id:   1,
				host: "127.0.0.1",
				port: 8080,
				typ:  "htp",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid port",
			args: args{
				id:   1,
				host: "127.0.0.1",
				port: -2,
				typ:  "http",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProxy(tt.args.id, tt.args.host, tt.args.port, tt.args.typ)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProxy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProxy() got = %v, want %v", got, tt.want)
			}
		})
	}
}
