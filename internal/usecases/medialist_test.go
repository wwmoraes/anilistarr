package usecases_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/wwmoraes/anilistarr/internal/entities"
	"github.com/wwmoraes/anilistarr/internal/testdata"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

func TestMediaList_invalid(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.TODO()
	username := "foo"
	sourceIDs := []entities.SourceID{"1"}

	mediaLister := usecases.MediaList{
		Source:  nil,
		Store:   nil,
		Tracker: nil,
	}

	gotGenerate, err := mediaLister.Generate(ctx, username)
	require.ErrorIs(t, err, usecases.ErrStatusFailedPrecondition)
	assert.Nil(t, gotGenerate)

	gotUserID, err := mediaLister.GetUserID(ctx, username)
	require.ErrorIs(t, err, usecases.ErrStatusFailedPrecondition)
	assert.Empty(t, gotUserID)

	gotTargetIDs, err := mediaLister.MapIDs(ctx, sourceIDs)
	require.ErrorIs(t, err, usecases.ErrStatusFailedPrecondition)
	assert.Nil(t, gotTargetIDs)

	err = mediaLister.Refresh(ctx, usecases.HTTPGetter(http.DefaultClient))
	require.ErrorIs(t, err, usecases.ErrStatusFailedPrecondition)

	err = mediaLister.Close()
	require.NoError(t, err)
}

func TestMediaList_Generate(t *testing.T) {
	t.Parallel()

	username := "foo"
	userID := "1"
	sourceIDs := []entities.SourceID{"1", "2", "3", "5", "8", "13"}
	medias := []*entities.Media{
		{SourceID: "1", TargetID: "91"},
		{SourceID: "2", TargetID: "92"},
		{SourceID: "3", TargetID: "93"},
		{SourceID: "5", TargetID: "95"},
		{SourceID: "8", TargetID: "98"},
		{SourceID: "13", TargetID: "913"},
	}
	customList := entities.CustomList{
		entities.CustomEntry{TvdbID: 91},
		entities.CustomEntry{TvdbID: 92},
		entities.CustomEntry{TvdbID: 93},
		entities.CustomEntry{TvdbID: 95},
		entities.CustomEntry{TvdbID: 98},
		entities.CustomEntry{TvdbID: 913},
	}

	source := testdata.MockSource{}
	store := testdata.MockStore{}
	tracker := testdata.MockTracker{}

	tracker.On("GetUserID", mock.Anything, username).
		Return(userID, nil).Once()
	tracker.On("GetMediaListIDs", mock.Anything, userID).
		Return(sourceIDs, nil).Once()
	store.On("GetMediaBulk", mock.Anything, sourceIDs).
		Return(medias, nil).Once()

	mediaLister := usecases.MediaList{
		Source:  &source,
		Store:   &store,
		Tracker: &tracker,
	}

	got, err := mediaLister.Generate(context.TODO(), username)
	require.NoError(t, err)

	assert.Equal(t, customList, got)
	source.AssertExpectations(t)
	store.AssertExpectations(t)
	tracker.AssertExpectations(t)
}

func TestMediaList_Generate_GetUserID_error(t *testing.T) {
	t.Parallel()

	username := "foo"

	source := testdata.MockSource{}
	store := testdata.MockStore{}
	tracker := testdata.MockTracker{}

	tracker.On("GetUserID", mock.Anything, username).
		Return("", usecases.ErrStatusNotFound).Once()

	mediaLister := usecases.MediaList{
		Source:  &source,
		Store:   &store,
		Tracker: &tracker,
	}

	got, err := mediaLister.Generate(context.TODO(), username)
	require.ErrorIs(t, err, usecases.ErrStatusNotFound)

	assert.Nil(t, got)
	source.AssertExpectations(t)
	store.AssertExpectations(t)
	tracker.AssertExpectations(t)
}

func TestMediaList_Generate_GetMediaListIDs_error(t *testing.T) {
	t.Parallel()

	var sourceIDs []entities.SourceID

	username := "foo"
	userID := "1"

	source := testdata.MockSource{}
	store := testdata.MockStore{}
	tracker := testdata.MockTracker{}

	tracker.On("GetUserID", mock.Anything, username).
		Return(userID, nil).Once()
	tracker.On("GetMediaListIDs", mock.Anything, userID).
		Return(sourceIDs, usecases.ErrStatusNotFound).Once()

	mediaLister := usecases.MediaList{
		Source:  &source,
		Store:   &store,
		Tracker: &tracker,
	}

	got, err := mediaLister.Generate(context.TODO(), username)
	require.ErrorIs(t, err, usecases.ErrStatusNotFound)

	assert.Nil(t, got)
	source.AssertExpectations(t)
	store.AssertExpectations(t)
	tracker.AssertExpectations(t)
}

func TestMediaList_Generate_Store_error(t *testing.T) {
	t.Parallel()

	var medias []*entities.Media

	username := "foo"
	userID := "1"
	sourceIDs := []entities.SourceID{"1", "2", "3", "5", "8", "13"}

	source := testdata.MockSource{}
	store := testdata.MockStore{}
	tracker := testdata.MockTracker{}

	tracker.On("GetUserID", mock.Anything, username).
		Return(userID, nil).Once()
	tracker.On("GetMediaListIDs", mock.Anything, userID).
		Return(sourceIDs, nil).Once()
	store.On("GetMediaBulk", mock.Anything, sourceIDs).
		Return(medias, usecases.ErrStatusUnknown).Once()

	mediaLister := usecases.MediaList{
		Source:  &source,
		Store:   &store,
		Tracker: &tracker,
	}

	got, err := mediaLister.Generate(context.TODO(), username)
	require.ErrorIs(t, err, usecases.ErrStatusUnknown)

	assert.Nil(t, got)
	source.AssertExpectations(t)
	store.AssertExpectations(t)
	tracker.AssertExpectations(t)
}

func TestMediaList_Generate_parse_error(t *testing.T) {
	t.Parallel()

	username := "foo"
	userID := "1"
	sourceIDs := []entities.SourceID{"1", "2", "3", "5", "8", "13"}
	medias := []*entities.Media{
		{SourceID: "1", TargetID: "91"},
		{SourceID: "2", TargetID: "92"},
		{SourceID: "3", TargetID: "93"},
		{SourceID: "5", TargetID: "95"},
		{SourceID: "8", TargetID: "98"},
		{SourceID: "13", TargetID: ""},
	}

	source := testdata.MockSource{}
	store := testdata.MockStore{}
	tracker := testdata.MockTracker{}

	tracker.On("GetUserID", mock.Anything, username).
		Return(userID, nil).Once()
	tracker.On("GetMediaListIDs", mock.Anything, userID).
		Return(sourceIDs, nil).Once()
	store.On("GetMediaBulk", mock.Anything, sourceIDs).
		Return(medias, nil).Once()

	mediaLister := usecases.MediaList{
		Source:  &source,
		Store:   &store,
		Tracker: &tracker,
	}

	got, err := mediaLister.Generate(context.TODO(), username)
	require.Error(t, err)

	assert.Nil(t, got)
	source.AssertExpectations(t)
	store.AssertExpectations(t)
	tracker.AssertExpectations(t)
}

func TestMediaList_GetUserID(t *testing.T) {
	t.Parallel()

	username := "foo"
	userID := "1"

	source := testdata.MockSource{}
	store := testdata.MockStore{}
	tracker := testdata.MockTracker{}

	tracker.On("GetUserID", mock.Anything, username).
		Return(userID, nil).Once()

	mediaLister := usecases.MediaList{
		Source:  &source,
		Store:   &store,
		Tracker: &tracker,
	}

	got, err := mediaLister.GetUserID(context.TODO(), username)
	require.NoError(t, err)

	assert.Equal(t, userID, got)
	source.AssertExpectations(t)
	store.AssertExpectations(t)
	tracker.AssertExpectations(t)
}

func TestMediaList_Refresh(t *testing.T) {
	t.Parallel()

	getter := usecases.HTTPGetter(nil)
	data := []usecases.Metadata{
		testdata.Metadata{SourceID: "1", TargetID: "91"},
		testdata.Metadata{SourceID: "2", TargetID: "92"},
		testdata.Metadata{SourceID: "3", TargetID: "93"},
		testdata.Metadata{SourceID: "5", TargetID: "95"},
		testdata.Metadata{SourceID: "8", TargetID: "98"},
		testdata.Metadata{SourceID: "13", TargetID: "913"},
		testdata.Metadata{SourceID: "999", TargetID: ""},
		testdata.Metadata{SourceID: "999", TargetID: "0"},
		testdata.Metadata{SourceID: "", TargetID: "999"},
		testdata.Metadata{SourceID: "0", TargetID: "999"},
	}
	medias := []*entities.Media{
		{SourceID: "1", TargetID: "91"},
		{SourceID: "2", TargetID: "92"},
		{SourceID: "3", TargetID: "93"},
		{SourceID: "5", TargetID: "95"},
		{SourceID: "8", TargetID: "98"},
		{SourceID: "13", TargetID: "913"},
	}

	source := testdata.MockSource{}
	store := testdata.MockStore{}
	tracker := testdata.MockTracker{}

	source.On("Fetch", mock.Anything, testdata.Implements[usecases.Getter](t)).
		Return(data, nil).Once()
	store.On("PutMediaBulk", mock.Anything, medias).
		Return(nil).Once()

	mediaLister := usecases.MediaList{
		Source:  &source,
		Store:   &store,
		Tracker: &tracker,
	}

	err := mediaLister.Refresh(context.TODO(), getter)
	require.NoError(t, err)

	source.AssertExpectations(t)
	store.AssertExpectations(t)
	tracker.AssertExpectations(t)
}

func TestMediaList_Refresh_Source_error(t *testing.T) {
	t.Parallel()

	var data []usecases.Metadata

	getter := usecases.HTTPGetter(nil)

	source := testdata.MockSource{}
	store := testdata.MockStore{}
	tracker := testdata.MockTracker{}

	source.On("Fetch", mock.Anything, testdata.Implements[usecases.Getter](t)).
		Return(data, errors.New("foo")).Once()

	mediaLister := usecases.MediaList{
		Source:  &source,
		Store:   &store,
		Tracker: &tracker,
	}

	err := mediaLister.Refresh(context.TODO(), getter)
	require.Error(t, err)

	source.AssertExpectations(t)
	store.AssertExpectations(t)
	tracker.AssertExpectations(t)
}

func TestMediaList_Refresh_Store_error(t *testing.T) {
	t.Parallel()

	getter := usecases.HTTPGetter(nil)
	data := []usecases.Metadata{
		testdata.Metadata{SourceID: "1", TargetID: "91"},
		testdata.Metadata{SourceID: "2", TargetID: "92"},
		testdata.Metadata{SourceID: "3", TargetID: "93"},
		testdata.Metadata{SourceID: "5", TargetID: "95"},
		testdata.Metadata{SourceID: "8", TargetID: "98"},
		testdata.Metadata{SourceID: "13", TargetID: "913"},
		testdata.Metadata{SourceID: "999", TargetID: ""},
		testdata.Metadata{SourceID: "999", TargetID: "0"},
		testdata.Metadata{SourceID: "", TargetID: "999"},
		testdata.Metadata{SourceID: "0", TargetID: "999"},
	}
	medias := []*entities.Media{
		{SourceID: "1", TargetID: "91"},
		{SourceID: "2", TargetID: "92"},
		{SourceID: "3", TargetID: "93"},
		{SourceID: "5", TargetID: "95"},
		{SourceID: "8", TargetID: "98"},
		{SourceID: "13", TargetID: "913"},
	}

	source := testdata.MockSource{}
	store := testdata.MockStore{}
	tracker := testdata.MockTracker{}

	source.On("Fetch", mock.Anything, testdata.Implements[usecases.Getter](t)).
		Return(data, nil).Once()
	store.On("PutMediaBulk", mock.Anything, medias).
		Return(errors.New("foo")).Once()

	mediaLister := usecases.MediaList{
		Source:  &source,
		Store:   &store,
		Tracker: &tracker,
	}

	err := mediaLister.Refresh(context.TODO(), getter)
	require.Error(t, err)

	source.AssertExpectations(t)
	store.AssertExpectations(t)
	tracker.AssertExpectations(t)
}

func TestMediaList_MapIDs(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	sourceIDs := []entities.SourceID{"1", "2", "3", "5", "8", "13", "999", ""}
	targetIDs := []entities.TargetID{"91", "92", "93", "95", "98", "913"}
	medias := []*entities.Media{
		{SourceID: "1", TargetID: "91"},
		{SourceID: "2", TargetID: "92"},
		{SourceID: "3", TargetID: "93"},
		{SourceID: "5", TargetID: "95"},
		{SourceID: "8", TargetID: "98"},
		{SourceID: "13", TargetID: "913"},
	}

	source := testdata.MockSource{}
	store := testdata.MockStore{}
	tracker := testdata.MockTracker{}

	store.On("GetMediaBulk", mock.Anything, sourceIDs).
		Return(medias, nil).Once()

	mediaLister := usecases.MediaList{
		Source:  &source,
		Store:   &store,
		Tracker: &tracker,
	}

	got, err := mediaLister.MapIDs(ctx, sourceIDs)
	require.NoError(t, err)

	assert.Equal(t, targetIDs, got)
	source.AssertExpectations(t)
	store.AssertExpectations(t)
	tracker.AssertExpectations(t)
}

func TestMediaList_MapIDs_Store_error(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	sourceIDs := []entities.SourceID{"1", "2", "3", "5", "8", "13", "999", ""}

	source := testdata.MockSource{}
	tracker := testdata.MockTracker{}

	mediaLister := usecases.MediaList{
		Source:  &source,
		Store:   nil,
		Tracker: &tracker,
	}

	got, err := mediaLister.MapIDs(ctx, sourceIDs)
	require.Error(t, err)

	assert.Nil(t, got)
	source.AssertExpectations(t)
	tracker.AssertExpectations(t)
}

func TestMediaList_Close(t *testing.T) {
	t.Parallel()

	source := testdata.MockSource{}
	store := testdata.MockStore{}
	tracker := testdata.MockTracker{}

	store.On("Close").Return(nil).Once()
	tracker.On("Close").Return(nil).Once()

	mediaLister := usecases.MediaList{
		Source:  &source,
		Store:   &store,
		Tracker: &tracker,
	}

	err := mediaLister.Close()
	require.NoError(t, err)

	source.AssertExpectations(t)
	store.AssertExpectations(t)
	tracker.AssertExpectations(t)
}
