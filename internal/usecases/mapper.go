package usecases

import (
	"context"
	"io"

	"github.com/wwmoraes/anilistarr/internal/entities"
)

// Mapper transforms IDs between two services based on the data from a provider
type Mapper interface {
	io.Closer

	MapIDs(context.Context, []entities.SourceID) ([]entities.TargetID, error)
	MapID(context.Context, entities.SourceID) (entities.TargetID, error)
	Refresh(context.Context, Getter) error
}
