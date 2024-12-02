package usecases

// Metadata represents any data that contains both the source Tracker media ID
// and the target service media ID
type Metadata interface {
	GetSourceID() string
	GetTargetID() string
	Valid() bool
}
