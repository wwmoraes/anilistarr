package adapters

import (
	"context"
	"fmt"

	"github.com/wwmoraes/anilistarr/internal/usecases"
)

// Provider is a source of mapping media data
type Provider[T Metadata] interface {
	fmt.Stringer

	Fetch(ctx context.Context, client usecases.Getter) ([]T, error)
}
