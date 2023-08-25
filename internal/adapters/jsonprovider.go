package adapters

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/wwmoraes/anilistarr/internal/telemetry"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

// JSONLocalProvider is a HTTP-based metadata provider
type JSONProvider[F Metadata] string

func (source JSONProvider[F]) String() string {
	return string(source)
}

func (source JSONProvider[F]) Fetch(ctx context.Context, client usecases.Getter) ([]Metadata, error) {
	_, span := telemetry.StartFunction(ctx)
	defer span.End()

	if client == nil {
		client = http.DefaultClient
	}

	res, err := client.Get(string(source))
	if err != nil {
		return nil, span.Assert(fmt.Errorf("failed to fetch remote JSON: %w", err))
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("provider data not found")
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, span.Assert(fmt.Errorf("failed to read fetched JSON response body: %w", err))
	}

	metadata, err := unmarshalJSON[F](data)

	return metadata, span.Assert(err)
}
