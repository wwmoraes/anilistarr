package usecases

import (
	"context"
	"fmt"
)

// Source is a data origin that provides metadata for mapping media IDs
//
//mockery:generate: true
type Source interface {
	fmt.Stringer

	Fetch(ctx context.Context, client Getter) ([]Metadata, error)
}
