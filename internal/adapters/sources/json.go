// Package sources contains [usecases.Source] implementations to retrieve
// metadata from specific encode formats such as JSON. They are agnostic of
// its source location, as in, its their client responsibility to take action
// locally (e.g. read a local file) or remotely (e.g. fetch a file over HTTP).
package sources

import (
	"context"
	"errors"
	"fmt"

	"github.com/goccy/go-json"
	telemetry "github.com/wwmoraes/gotell"

	"github.com/wwmoraes/anilistarr/internal/usecases"
)

// JSON is a URI to a metadata source marshaled in JSON.
type JSON[F usecases.Metadata] string

// String returns the provider URI.
func (source JSON[F]) String() string {
	return string(source)
}

// Fetch retrieves and parses metadata from the provider URI.
func (source JSON[F]) Fetch(
	ctx context.Context,
	client usecases.Getter,
) ([]usecases.Metadata, error) {
	_, span := telemetry.Start(ctx)
	defer span.End()

	if client == nil {
		return nil, fmt.Errorf("%w: %s", usecases.ErrStatusInternal, "no getter set")
	}

	data, err := client.Get(ctx, source.String())
	if err != nil {
		return nil, span.Assert(
			errors.Join(usecases.ErrStatusUnavailable, fmt.Errorf("failed to get JSON: %w", err)),
		)
	}

	metadata, err := unmarshalJSON[F](data)

	return metadata, span.Assert(err)
}

// unmarshalJSON parses raw bytes into an list of metadata elements.
//
// TODO refactor to use JSON interfaces instead
func unmarshalJSON[F usecases.Metadata](data []byte) ([]usecases.Metadata, error) {
	var dataEntries []F

	err := json.Unmarshal(data, &dataEntries)
	if err != nil {
		return nil, errors.Join(
			usecases.ErrStatusFailedPrecondition,
			fmt.Errorf("failed to unmarshal JSON: %w", err),
		)
	}

	entries := make([]usecases.Metadata, 0, len(dataEntries))
	for _, entry := range dataEntries {
		entries = append(entries, entry)
	}

	return entries, nil
}
