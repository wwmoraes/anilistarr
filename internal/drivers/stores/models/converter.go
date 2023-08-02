package models

import (
	"database/sql"

	"github.com/wwmoraes/anilistarr/internal/entities"
)

func (mapping *Mapping) ToMedia() *entities.Media {
	return &entities.Media{
		AnilistID: mapping.AnilistID.String,
		TvdbID:    mapping.TvdbID,
	}
}

func MappingFromMedia(media *entities.Media) *Mapping {
	if media == nil {
		return nil
	}

	return &Mapping{
		AnilistID: sql.NullString{
			String: media.AnilistID,
			Valid:  len(media.AnilistID) > 0,
		},
		TvdbID: media.TvdbID,
	}
}
