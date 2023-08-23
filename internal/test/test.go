package test

import (
	"strconv"
	"testing"

	"github.com/wwmoraes/anilistarr/internal/entities"
)

func SonarrCustomListFromIDs(tb testing.TB, ids ...string) entities.SonarrCustomList {
	tb.Helper()

	customList := make(entities.SonarrCustomList, len(ids))

	for index, id := range ids {
		tvdbId, err := strconv.ParseUint(id, 10, 0)
		if err != nil {
			tb.Fatal(err)
		}

		customList[index] = entities.SonarrCustomEntry{
			TvdbID: tvdbId,
		}
	}

	return customList
}
