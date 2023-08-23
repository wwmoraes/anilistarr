package entities

// TODO change to uint64
type SourceID = string
type TargetID = string

type Media struct {
	SourceID SourceID `json:"source_id,omitempty" db:"source_id"`
	TargetID TargetID `json:"target_id,omitempty" db:"target_id"`
}
