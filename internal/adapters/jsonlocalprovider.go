package adapters

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/wwmoraes/anilistarr/internal/telemetry"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

// JSONLocalProvider is a filesystem-based metadata provider
type JSONLocalProvider[F Metadata] struct {
	Fs   fs.FS
	Name string
}

func (source JSONLocalProvider[F]) String() string {
	return string(source.Name)
}

func (source JSONLocalProvider[F]) Fetch(ctx context.Context, client usecases.Getter) ([]Metadata, error) {
	_, span := telemetry.StartFunction(ctx)
	defer span.End()

	data, err := fs.ReadFile(source.Fs, source.String())
	if err != nil {
		return nil, span.Assert(fmt.Errorf("failed to read local JSON: %w", err))
	}

	metadata, err := unmarshalJSON[F](data)

	return metadata, span.Assert(err)
}
