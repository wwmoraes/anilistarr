package testdata

import (
	"strconv"
	"testing"

	"github.com/wwmoraes/anilistarr/internal/entities"
)

func CustomListFromIDs(tb testing.TB, ids ...string) entities.CustomList {
	tb.Helper()

	customList := make(entities.CustomList, len(ids))

	for index, id := range ids {
		tvdbId, err := strconv.ParseUint(id, 10, 0)
		if err != nil {
			tb.Fatal(err)
		}

		customList[index] = entities.CustomEntry{
			TvdbID: tvdbId,
		}
	}

	return customList
}
