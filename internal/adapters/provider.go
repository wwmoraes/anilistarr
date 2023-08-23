package adapters

import (
	"context"
	"fmt"

	"github.com/wwmoraes/anilistarr/internal/usecases"
)

type Provider[T Metadata] interface {
	fmt.Stringer

	Fetch(ctx context.Context, client usecases.Getter) ([]T, error)
}
