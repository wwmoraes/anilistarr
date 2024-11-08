package usecases_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/wwmoraes/anilistarr/internal/usecases"
)

func TestGetter(t *testing.T) {
	t.Parallel()

	type args struct {
		server *httptest.Server
		getter usecases.HTTPGetter
	}

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "200",
			args: args{
				server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("foo"))
				})),
				getter: &http.Client{},
			},
			want:    []byte("foo"),
			wantErr: false,
		},
		{
			name: "500",
			args: args{
				server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				})),
				getter: &http.Client{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid URL",
			args: args{
				server: func() *httptest.Server {
					server := httptest.NewServer(nil)

					server.URL = "bogus"

					return server
				}(),
				getter: &http.Client{},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			defer tt.args.server.Close()

			getter := usecases.HTTPGetterAsGetter(tt.args.getter)

			got, err := getter.Get(tt.args.server.URL)
			if (err != nil) != tt.wantErr {
				t.Errorf("Getter.Get() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Getter.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
