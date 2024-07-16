package usecases

import (
	"context"
	"io"

	"github.com/wwmoraes/anilistarr/internal/entities"
)

// Mapper transforms IDs between two services based on the data from a provider
type Mapper interface {
	io.Closer

	MapIDs(ctx context.Context, ids []entities.SourceID) ([]entities.TargetID, error)
	Refresh(ctx context.Context, client Getter) error
}
