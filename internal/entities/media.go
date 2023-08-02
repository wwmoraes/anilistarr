package entities

type Media struct {
	AnilistID string `json:"anilist_id,omitempty" db:"anilist_id"`
	TvdbID    string `json:"thetvdb_id,omitempty" db:"thetvdb_id"`

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
