package testdata

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/wwmoraes/anilistarr/internal/entities"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

type Tracker struct {
	UserIds    map[string]int
	MediaLists map[int][]entities.SourceID
}

func (tracker *Tracker) GetUserID(ctx context.Context, name string) (string, error) {
	if tracker.UserIds == nil {
		return "", fmt.Errorf(http.StatusText(http.StatusBadGateway))
	}

	id, ok := tracker.UserIds[name]
	if !ok {
		return "", fmt.Errorf(usecases.FailedGetUserErrorTemplate, fmt.Errorf("user id not found"))
	}

	return strconv.Itoa(id), nil
}

func (tracker *Tracker) GetMediaListIDs(ctx context.Context, userId string) ([]entities.SourceID, error) {
	if tracker.MediaLists == nil {
		return nil, fmt.Errorf(http.StatusText(http.StatusBadGateway))
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return nil, fmt.Errorf(usecases.ConvertUserIDErrorTemplate, err)
	}

	value, ok := tracker.MediaLists[userIdInt]
	if !ok {
		return []entities.SourceID{}, nil
	}

	return value, nil
}

func (tracker *Tracker) Close() error {
	*tracker = Tracker{
		UserIds:    make(map[string]int),
		MediaLists: make(map[int][]string),
	}

	return nil
}
