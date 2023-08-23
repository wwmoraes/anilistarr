package providers

import (
	"strconv"

	"github.com/wwmoraes/anilistarr/internal/adapters"
)

const AnilistFribbsProvider adapters.JSONProvider[AnilistFribbsMetadata] = "https://github.com/Fribb/anime-lists/raw/master/anime-list-full.json"

type AnilistFribbsMetadata struct {
	AnilistID uint64 `json:"anilist_id,omitempty"`
	TvdbID    uint64 `json:"thetvdb_id,omitempty"`
}

func (entry AnilistFribbsMetadata) GetTargetID() string {
	return strconv.FormatUint(entry.TvdbID, 10)
}

func (entry AnilistFribbsMetadata) GetSourceID() string {
	return strconv.FormatUint(entry.AnilistID, 10)
}
