package entities

type (
	// SourceID represents a media ID from a source service
	SourceID = string

	// TargetID represents a media ID from a target service
	TargetID = string
)

// Media represents a relationship between two services. It contains the ID of
// the same media on them both.
//
//nolint:tagliatelle // JSON tags must match upstream naming convention
type Media struct {
	SourceID SourceID `db:"source_id" json:"source_id,omitempty"`
	TargetID TargetID `db:"target_id" json:"target_id,omitempty"`
}

// Valid returns true if this is a valid media i.e. it contains both non-empty,
// non-zero source and target IDs.
func (media *Media) Valid() bool {
	return media.SourceID != "" &&
		media.TargetID != "" &&
		media.SourceID != "0" &&
		media.TargetID != "0"
}
