package testdata

import "github.com/wwmoraes/anilistarr/internal/drivers/memory"

const (
	Username = "test"
	UserID   = 1234
)

var (
	SourceIDs    = []string{"1", "2", "3", "5", "8", "13"}
	TargetIDs    = []string{"91", "92", "93", "95", "98", "913"}
	SampleClient = HTTPClient{
		Data: map[string]string{
			Provider.String(): `[
				{"anilist_id": 1, "thetvdb_id": 91},
				{"anilist_id": 2, "thetvdb_id": 92},
				{"anilist_id": 3, "thetvdb_id": 93},
				{"anilist_id": 5, "thetvdb_id": 95},
				{"anilist_id": 8, "thetvdb_id": 98},
				{"anilist_id": 13, "thetvdb_id": 913}
			]`,
		},
	}
	SampleTracker = &Tracker{
		UserIds: map[string]int{
			Username: UserID,
		},
		MediaLists: map[int][]string{
			UserID: SourceIDs,
		},
	}
	SampleStore = &memory.Memory{
		"1":  "91",
		"2":  "92",
		"3":  "93",
		"5":  "95",
		"8":  "98",
		"13": "913",
	}
)
