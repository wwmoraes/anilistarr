package test

import (
	"strconv"

	"github.com/wwmoraes/anilistarr/internal/adapters"
)

// type Provider[T adapters.Metadata] []T

const Provider adapters.JSONProvider[Metadata] = `memory:///test`

//nolint:tagliatelle // JSON tags must match the upstream naming convention
type Metadata struct {
	SourceID uint64 `json:"anilist_id,omitempty"`
	TargetID uint64 `json:"thetvdb_id,omitempty"`
}

func (metadata Metadata) GetSourceID() string {
	return strconv.FormatUint(metadata.SourceID, 10)
}

func (metadata Metadata) GetTargetID() string {
	return strconv.FormatUint(metadata.TargetID, 10)
}
