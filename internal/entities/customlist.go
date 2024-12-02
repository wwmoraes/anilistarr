package entities

// CustomList contains custom entries in the Sonarr format.
type CustomList []CustomEntry

// CustomEntry represents one media entry in the Sonarr format.
type CustomEntry struct {
	TvdbID uint64
}
