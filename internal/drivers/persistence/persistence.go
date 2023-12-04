package persistence

import "github.com/wwmoraes/anilistarr/internal/adapters"

// Persistence represents objects that provide data storage and persistence,
// eventual or immediate, for use as both Cache and/or Store
type Persistence interface {
	adapters.Cache
	adapters.Store
}
