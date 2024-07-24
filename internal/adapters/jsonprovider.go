package adapters

import (
	"context"
	"fmt"

	telemetry "github.com/wwmoraes/gotell"

	"github.com/wwmoraes/anilistarr/internal/usecases"
)

// JSONLocalProvider is a HTTP-based metadata provider
type JSONProvider[F Metadata] string

func (source JSONProvider[F]) String() string {
	return string(source)
}

func (source JSONProvider[F]) Fetch(ctx context.Context, client usecases.Getter) ([]Metadata, error) {
	_, span := telemetry.Start(ctx)
	defer span.End()

	if client == nil {
		return nil, ErrNoGetter
	}

	data, err := client.Get(string(source))
	if err != nil {
		return nil, span.Assert(fmt.Errorf("failed to get JSON: %w", err))
	}

	metadata, err := unmarshalJSON[F](data)

	return metadata, span.Assert(err)
}
