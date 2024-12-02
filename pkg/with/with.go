// Package with implements an idiomatic builder pattern.
//
// The difference from object-oriented builder patterns is that here we use an
// auxiliary struct with all relevant properties to its consumer. This struct
// then changes using functors that apply some mutation to it in-place.
package with

// Option values can apply modifications to a target options type.
type Option[O any] interface {
	Apply(options *O)
}

// Functor is a function that modifies a target options type.
type Functor[O any] func(options *O)

// Apply changes a target options type.
func (fn Functor[O]) Apply(options *O) {
	fn(options)
}

// Apply folds all option values in order on top of the defaults, returning a
// copy of the resulting target options type.
func Apply[O any](defaults O, opts ...Option[O]) O {
	for _, opt := range opts {
		opt.Apply(&defaults)
	}

	return defaults
}
