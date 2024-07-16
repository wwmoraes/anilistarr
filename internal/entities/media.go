package entities

// TODO change to uint64
type (
	SourceID = string
	TargetID = string
)

//nolint:tagliatelle // JSON tags must match the upstream naming convention
type Media struct {
	SourceID SourceID `db:"source_id" json:"source_id,omitempty"`
	TargetID TargetID `db:"target_id" json:"target_id,omitempty"`
}
