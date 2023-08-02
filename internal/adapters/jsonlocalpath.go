package adapters

import (
	"context"
	"fmt"
	"os"

	"github.com/wwmoraes/anilistarr/internal/telemetry"
)

type JSONLocalPath[F Metadata] string

func (source JSONLocalPath[F]) Fetch(ctx context.Context) ([]Metadata, error) {
	_, span := telemetry.StartFunction(ctx)
	defer span.End()

	data, err := os.ReadFile(string(source))
	if err != nil {
		return nil, span.Assert(fmt.Errorf("failed to read local JSON: %w", err))
	}

	metadata, err := unmarshalJSON[F](data)

	return metadata, span.Assert(err)
}
