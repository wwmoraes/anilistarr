package adapters_test

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/wwmoraes/anilistarr/internal/adapters"
)

func TestNewCacheParams(t *testing.T) {
	t.Parallel()

	type args struct {
		options []adapters.CacheOption
	}

	tests := []struct {
		name    string
		args    args
		want    *adapters.CacheParams
		wantErr bool
	}{
		{
			name: "broken option",
			args: args{
				options: []adapters.CacheOption{
					adapters.CacheOptionFn(func(params *adapters.CacheParams) error {
						return errors.New("yay error")
					}),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "with TTL",
			args: args{
				options: []adapters.CacheOption{
					adapters.WithTTL(time.Hour),
				},
			},
			want: &adapters.CacheParams{
				TTL: time.Hour,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := adapters.NewCacheParams(tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCacheParams() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCacheParams() = %v, want %v", got, tt.want)
			}
		})
	}
}
