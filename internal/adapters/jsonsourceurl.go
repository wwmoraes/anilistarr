package adapters

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/wwmoraes/anilistarr/internal/telemetry"
)

type JSONSourceURL[F Metadata] string

func (source JSONSourceURL[F]) Fetch(ctx context.Context, client Getter) ([]Metadata, error) {
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

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, span.Assert(fmt.Errorf("failed to read fetched JSON response body: %w", err))
	}

	metadata, err := unmarshalJSON[F](data)

	return metadata, span.Assert(err)
}
