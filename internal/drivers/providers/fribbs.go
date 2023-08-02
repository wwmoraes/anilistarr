package providers

import (
	"strconv"

	"github.com/wwmoraes/anilistarr/internal/adapters"
)

const FribbsSource adapters.JSONSourceURL[FribbsEntry] = "https://github.com/Fribb/anime-lists/raw/master/anime-list-full.json"

type FribbsEntry struct {
	AnilistID uint64 `json:"anilist_id,omitempty"`
	TvdbID    uint64 `json:"thetvdb_id,omitempty"`

	//// useless
	// Type      string `json:"type,omitempty"`

	//// commented out as we don't need these
	// AnidbID     uint   `json:"anidb_id,omitempty"`
	// AnisearchID uint   `json:"anisearch_id,omitempty"`
	// ImdbID      string `json:"imdb_id,omitempty"`
	// KitsuID     uint   `json:"kitsu_id,omitempty"`
	// LivechartID uint   `json:"livechart_id,omitempty"`
	// MalID       uint   `json:"mal_id,omitempty"`
	// NotifyMoeID string `json:"notify.moe_id,omitempty"`

	//// those are even worse as they mix strings and numbers
	// AnimePlanetID string `json:"anime-planet_id,omitempty"`
	// TmdbID      uint   `json:"themoviedb_id,omitempty"`
}

func (entry FribbsEntry) GetTvdbID() string {
	return strconv.FormatUint(entry.TvdbID, 10)
}

func (entry FribbsEntry) GetAnilistID() string {
	return strconv.FormatUint(entry.AnilistID, 10)
}
