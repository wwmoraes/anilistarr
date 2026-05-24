package with_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/wwmoraes/anilistarr/pkg/with"
)

type testConfig struct {
	foo string
}

func TestApply(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		defaults testConfig
		want     testConfig
		options  []with.Option[testConfig]
	}{
		{
			name: "no-op",
			defaults: testConfig{
				foo: "bar",
			},
			options: nil,
			want: testConfig{
				foo: "bar",
			},
		},
		{
			name: "change foo",
			defaults: testConfig{
				foo: "bar",
			},
			options: []with.Option[testConfig]{
				with.Functor[testConfig](func(options *testConfig) {
					options.foo = "qux"
				}),
			},
			want: testConfig{
				foo: "qux",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			config := with.Apply(tt.defaults, tt.options...)
			require.Equal(t, tt.want, config)
		})
	}
}
