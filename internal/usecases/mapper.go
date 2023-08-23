package usecases

import (
	"context"
	"io"

	"github.com/wwmoraes/anilistarr/internal/entities"
)

type Mapper interface {
	io.Closer

	MapIDs(context.Context, []entities.SourceID) ([]entities.TargetID, error)
	MapID(context.Context, entities.SourceID) (entities.TargetID, error)
	Refresh(context.Context, Getter) error
}
