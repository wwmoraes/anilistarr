package adapters

import (
	"encoding/json"
	"fmt"
)

type Metadata interface {
	GetAnilistID() string
	GetTvdbID() string
}

func unmarshalJSON[F Metadata](data []byte) ([]Metadata, error) {
	var dataEntries []F
	err := json.Unmarshal(data, &dataEntries)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	entries := make([]Metadata, len(dataEntries))
	for index, entry := range dataEntries {
		entries[index] = entry
	}

	return entries, nil
}
