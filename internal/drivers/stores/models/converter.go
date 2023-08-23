package models

import (
	"github.com/wwmoraes/anilistarr/internal/entities"
)

func (mapping *Mapping) ToMedia() *entities.Media {
	return &entities.Media{
		SourceID: mapping.SourceID,
		TargetID: mapping.TargetID,
	}
}

func MappingFromMedia(media *entities.Media) *Mapping {
	if media == nil {
		return nil
	}

	return &Mapping{
		SourceID: media.SourceID,
		TargetID: media.TargetID,
	}
}
