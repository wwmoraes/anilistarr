package usecases

import (
	"context"
	"io"

	"github.com/wwmoraes/anilistarr/internal/entities"
)

// MediaLister converts media IDs between services and generates Sonarr custom
// list results. It handles providers, trackers and stores to fetch the data
// required for such conversions.
type MediaLister interface {
	io.Closer

	// Generate fetches the user media list from the Tracker and transform the IDs
	// found to the target service through the Mapper
	Generate(ctx context.Context, name string) (entities.CustomList, error)

	// GetUserID searches the Tracker for the user ID by their name/handle
	GetUserID(ctx context.Context, name string) (string, error)

	// MapIDs matches media IDs between two services
	MapIDs(ctx context.Context, ids []entities.SourceID) ([]entities.TargetID, error)

	// Refresh requests the Mapper to update its mapping definitions
	Refresh(ctx context.Context, client Getter) error
}
