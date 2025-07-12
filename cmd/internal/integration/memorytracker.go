package main

import (
	"context"
	"errors"
	"strconv"

	"github.com/wwmoraes/anilistarr/internal/entities"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

var _ usecases.Tracker = (*memoryTracker)(nil)

type memoryTracker struct {
	UserIDs    map[string]int
	MediaLists map[int][]entities.SourceID
}

// GetUserID returns the internal ID of a registered user
func (tracker *memoryTracker) GetUserID(_ context.Context, name string) (string, error) {
	if tracker.UserIDs == nil {
		return "", usecases.ErrStatusFailedPrecondition
	}

	id, ok := tracker.UserIDs[name]
	if !ok {
		return "", usecases.ErrStatusNotFound
	}

	return strconv.Itoa(id), nil
}

// GetMediaListIDs retrieves the media IDs of a registered user
func (tracker *memoryTracker) GetMediaListIDs(
	_ context.Context,
	userID string,
) ([]entities.SourceID, error) {
	if tracker.MediaLists == nil {
		return nil, usecases.ErrStatusFailedPrecondition
	}

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return nil, errors.Join(usecases.ErrStatusInvalidArgument, err)
	}

	value, ok := tracker.MediaLists[userIDInt]
	if !ok {
		return []entities.SourceID{}, nil
	}

	return value, nil
}

// Close cleans up data
func (tracker *memoryTracker) Close() error {
	*tracker = memoryTracker{
		UserIDs:    make(map[string]int),
		MediaLists: make(map[int][]string),
	}

	return nil
}
