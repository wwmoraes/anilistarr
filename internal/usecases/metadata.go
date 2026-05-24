package usecases

// Metadata represents any data that contains both the source Tracker media ID
// and the target service media ID
//
//mockery:generate: true
type Metadata interface {
	GetSourceID() string
	GetTargetID() string
	Valid() bool
}
