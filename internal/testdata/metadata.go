package testdata

import (
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

var (
	_ usecases.Metadata = (*Metadata)(nil)
)

//nolint:tagliatelle // JSON tags must match the upstream naming convention
type Metadata struct {
	SourceID string `json:"source_id,omitempty"`
	TargetID string `json:"target_id,omitempty"`
}

func (entry Metadata) GetSourceID() string {
	return entry.SourceID
}

func (entry Metadata) GetTargetID() string {
	return entry.TargetID
}

func (entry Metadata) Valid() bool {
	return entry.SourceID != "" &&
		entry.SourceID != "0" &&
		entry.TargetID != "" &&
		entry.TargetID != "0"
}
