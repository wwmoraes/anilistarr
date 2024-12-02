// Package animelists provides elements to consume data from the anime lists
// and its forks.
package animelists

import (
	"strconv"
)

// Metadata represents a media with IDs for multiple services.
//
// It expects an entry in the anime-lists project format. See:
//   - https://github.com/Anime-Lists/anime-lists/
//   - https://github.com/manami-project/anime-offline-database/
//   - https://github.com/Fribb/anime-lists
//
//nolint:tagliatelle // JSON tags must match the upstream naming convention
type Metadata struct {
	// AnimePlanetID string `json:"anime,omitempty"`
	// ImdbID        string `json:"imdb_id,omitempty"`
	// NotifyMoeID   string `json:"notify,omitempty"`
	// AnidbID       uint64 `json:"anidb_id,omitempty"`

	AnilistID uint64 `json:"anilist_id,omitempty"`

	// AnisearchID   uint64 `json:"anisearch_id,omitempty"`
	// KitsuID       uint64 `json:"kitsu_id,omitempty"`
	// LiveChartID   uint64 `json:"livechart_id,omitempty"`
	// MalID         uint64 `json:"mal_id,omitempty"`
	// TheMovieDbID  uint64 `json:"themoviedb_id,omitempty"`

	TvdbID uint64 `json:"thetvdb_id,omitempty"`
}

// Anilist2TVDBMetadata represents an anime lists entry with both
// Anilist and TVDB IDs, mapping the former as source and the latter as target.
type Anilist2TVDBMetadata Metadata

// GetTargetID retrieves the target ID of this metadata entry.
func (entry Anilist2TVDBMetadata) GetTargetID() string {
	return strconv.FormatUint(entry.TvdbID, 10)
}

// GetSourceID retrieves the source ID of this metadata entry.
func (entry Anilist2TVDBMetadata) GetSourceID() string {
	return strconv.FormatUint(entry.AnilistID, 10)
}

// Valid returns true if both source and target IDs are non-zero.
func (entry Anilist2TVDBMetadata) Valid() bool {
	return entry.AnilistID > 0 && entry.TvdbID > 0
}
