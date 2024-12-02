package usecases_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/wwmoraes/anilistarr/internal/usecases"
)

func TestCacheOptions(t *testing.T) {
	t.Parallel()

	type args struct {
		options []usecases.CacheOption
	}

	tests := []struct {
		want *usecases.CacheOptions
		name string
		args args
	}{
		{
			name: "default",
			args: args{
				options: []usecases.CacheOption{},
			},
			want: &usecases.CacheOptions{
				TTL: 0,
			},
		},
		{
			name: "WithTTL",
			args: args{
				options: []usecases.CacheOption{
					usecases.WithTTL(time.Nanosecond),
				},
			},
			want: &usecases.CacheOptions{
				TTL: time.Nanosecond,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, usecases.NewCacheOptions(tt.args.options...))
		})
	}
}
