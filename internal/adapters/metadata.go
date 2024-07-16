package adapters

import (
	"fmt"

	"github.com/goccy/go-json"
)

// Metadata represents any data that contains both the source Tracker media ID
// and the target service media ID
type Metadata interface {
	GetSourceID() string
	GetTargetID() string
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
