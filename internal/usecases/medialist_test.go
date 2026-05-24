package usecases_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/wwmoraes/anilistarr/internal/entities"
	"github.com/wwmoraes/anilistarr/internal/test"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

func TestMediaList_invalid(t *testing.T) {
	t.Parallel()

	var err error

	ctx := t.Context()
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

	source := test.NewMockSource(t)
	store := test.NewMockStore(t)
	tracker := test.NewMockTracker(t)

	tracker.EXPECT().GetUserID(mock.Anything, username).
		Return(userID, nil).Once()
	tracker.EXPECT().GetMediaListIDs(mock.Anything, userID).
		Return(sourceIDs, nil).Once()
	store.EXPECT().GetMediaBulk(mock.Anything, sourceIDs).
		Return(medias, nil).Once()

	mediaLister := usecases.MediaList{
		Source:  source,
		Store:   store,
		Tracker: tracker,
	}

	got, err := mediaLister.Generate(t.Context(), username)
	require.NoError(t, err)

	assert.Equal(t, customList, got)
}

func TestMediaList_Generate_GetUserID_error(t *testing.T) {
	t.Parallel()

	username := "foo"

	source := test.NewMockSource(t)
	store := test.NewMockStore(t)
	tracker := test.NewMockTracker(t)

	tracker.EXPECT().GetUserID(mock.Anything, username).
		Return("", usecases.ErrStatusNotFound).Once()

	mediaLister := usecases.MediaList{
		Source:  source,
		Store:   store,
		Tracker: tracker,
	}

	got, err := mediaLister.Generate(t.Context(), username)
	require.ErrorIs(t, err, usecases.ErrStatusNotFound)

	assert.Nil(t, got)
}

func TestMediaList_Generate_GetMediaListIDs_error(t *testing.T) {
	t.Parallel()

	var sourceIDs []entities.SourceID

	username := "foo"
	userID := "1"

	source := test.NewMockSource(t)
	store := test.NewMockStore(t)
	tracker := test.NewMockTracker(t)

	tracker.EXPECT().GetUserID(mock.Anything, username).
		Return(userID, nil).Once()
	tracker.EXPECT().GetMediaListIDs(mock.Anything, userID).
		Return(sourceIDs, usecases.ErrStatusNotFound).Once()

	mediaLister := usecases.MediaList{
		Source:  source,
		Store:   store,
		Tracker: tracker,
	}

	got, err := mediaLister.Generate(t.Context(), username)
	require.ErrorIs(t, err, usecases.ErrStatusNotFound)

	assert.Nil(t, got)
}

func TestMediaList_Generate_Store_error(t *testing.T) {
	t.Parallel()

	var medias []*entities.Media

	username := "foo"
	userID := "1"
	sourceIDs := []entities.SourceID{"1", "2", "3", "5", "8", "13"}

	source := test.NewMockSource(t)
	store := test.NewMockStore(t)
	tracker := test.NewMockTracker(t)

	tracker.EXPECT().GetUserID(mock.Anything, username).
		Return(userID, nil).Once()
	tracker.EXPECT().GetMediaListIDs(mock.Anything, userID).
		Return(sourceIDs, nil).Once()
	store.EXPECT().GetMediaBulk(mock.Anything, sourceIDs).
		Return(medias, usecases.ErrStatusUnknown).Once()

	mediaLister := usecases.MediaList{
		Source:  source,
		Store:   store,
		Tracker: tracker,
	}

	got, err := mediaLister.Generate(t.Context(), username)
	require.ErrorIs(t, err, usecases.ErrStatusUnknown)

	assert.Nil(t, got)
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

	source := test.NewMockSource(t)
	store := test.NewMockStore(t)
	tracker := test.NewMockTracker(t)

	tracker.EXPECT().GetUserID(mock.Anything, username).
		Return(userID, nil).Once()
	tracker.EXPECT().GetMediaListIDs(mock.Anything, userID).
		Return(sourceIDs, nil).Once()
	store.EXPECT().GetMediaBulk(mock.Anything, sourceIDs).
		Return(medias, nil).Once()

	mediaLister := usecases.MediaList{
		Source:  source,
		Store:   store,
		Tracker: tracker,
	}

	got, err := mediaLister.Generate(t.Context(), username)
	require.Error(t, err)

	assert.Nil(t, got)
}

func TestMediaList_GetUserID(t *testing.T) {
	t.Parallel()

	username := "foo"
	userID := "1"

	source := test.NewMockSource(t)
	store := test.NewMockStore(t)
	tracker := test.NewMockTracker(t)

	tracker.EXPECT().GetUserID(mock.Anything, username).
		Return(userID, nil).Once()

	mediaLister := usecases.MediaList{
		Source:  source,
		Store:   store,
		Tracker: tracker,
	}

	got, err := mediaLister.GetUserID(t.Context(), username)
	require.NoError(t, err)

	assert.Equal(t, userID, got)
}

func TestMediaList_Refresh(t *testing.T) {
	t.Parallel()

	getter := usecases.HTTPGetter(nil)
	data := []usecases.Metadata{
		test.Metadata{SourceID: "1", TargetID: "91"},
		test.Metadata{SourceID: "2", TargetID: "92"},
		test.Metadata{SourceID: "3", TargetID: "93"},
		test.Metadata{SourceID: "5", TargetID: "95"},
		test.Metadata{SourceID: "8", TargetID: "98"},
		test.Metadata{SourceID: "13", TargetID: "913"},
		test.Metadata{SourceID: "999", TargetID: ""},
		test.Metadata{SourceID: "999", TargetID: "0"},
		test.Metadata{SourceID: "", TargetID: "999"},
		test.Metadata{SourceID: "0", TargetID: "999"},
	}
	medias := []*entities.Media{
		{SourceID: "1", TargetID: "91"},
		{SourceID: "2", TargetID: "92"},
		{SourceID: "3", TargetID: "93"},
		{SourceID: "5", TargetID: "95"},
		{SourceID: "8", TargetID: "98"},
		{SourceID: "13", TargetID: "913"},
	}

	source := test.NewMockSource(t)
	store := test.NewMockStore(t)
	tracker := test.NewMockTracker(t)

	source.EXPECT().Fetch(mock.Anything, implements[usecases.Getter](t)).
		Return(data, nil).Once()
	store.EXPECT().PutMediaBulk(mock.Anything, medias).
		Return(nil).Once()

	mediaLister := usecases.MediaList{
		Source:  source,
		Store:   store,
		Tracker: tracker,
	}

	err := mediaLister.Refresh(t.Context(), getter)
	require.NoError(t, err)
}

func TestMediaList_Refresh_Source_error(t *testing.T) {
	t.Parallel()

	var data []usecases.Metadata

	getter := usecases.HTTPGetter(nil)

	source := test.NewMockSource(t)
	store := test.NewMockStore(t)
	tracker := test.NewMockTracker(t)

	source.EXPECT().Fetch(mock.Anything, implements[usecases.Getter](t)).
		Return(data, errors.New("foo")).Once()

	mediaLister := usecases.MediaList{
		Source:  source,
		Store:   store,
		Tracker: tracker,
	}

	err := mediaLister.Refresh(t.Context(), getter)
	require.Error(t, err)
}

func TestMediaList_Refresh_Store_error(t *testing.T) {
	t.Parallel()

	getter := usecases.HTTPGetter(nil)
	data := []usecases.Metadata{
		test.Metadata{SourceID: "1", TargetID: "91"},
		test.Metadata{SourceID: "2", TargetID: "92"},
		test.Metadata{SourceID: "3", TargetID: "93"},
		test.Metadata{SourceID: "5", TargetID: "95"},
		test.Metadata{SourceID: "8", TargetID: "98"},
		test.Metadata{SourceID: "13", TargetID: "913"},
		test.Metadata{SourceID: "999", TargetID: ""},
		test.Metadata{SourceID: "999", TargetID: "0"},
		test.Metadata{SourceID: "", TargetID: "999"},
		test.Metadata{SourceID: "0", TargetID: "999"},
	}
	medias := []*entities.Media{
		{SourceID: "1", TargetID: "91"},
		{SourceID: "2", TargetID: "92"},
		{SourceID: "3", TargetID: "93"},
		{SourceID: "5", TargetID: "95"},
		{SourceID: "8", TargetID: "98"},
		{SourceID: "13", TargetID: "913"},
	}

	source := test.NewMockSource(t)
	store := test.NewMockStore(t)
	tracker := test.NewMockTracker(t)

	source.EXPECT().Fetch(mock.Anything, implements[usecases.Getter](t)).
		Return(data, nil).Once()
	store.EXPECT().PutMediaBulk(mock.Anything, medias).
		Return(errors.New("foo")).Once()

	mediaLister := usecases.MediaList{
		Source:  source,
		Store:   store,
		Tracker: tracker,
	}

	err := mediaLister.Refresh(t.Context(), getter)
	require.Error(t, err)
}

func TestMediaList_MapIDs(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
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

	source := test.NewMockSource(t)
	store := test.NewMockStore(t)
	tracker := test.NewMockTracker(t)

	store.EXPECT().GetMediaBulk(mock.Anything, sourceIDs).
		Return(medias, nil).Once()

	mediaLister := usecases.MediaList{
		Source:  source,
		Store:   store,
		Tracker: tracker,
	}

	got, err := mediaLister.MapIDs(ctx, sourceIDs)
	require.NoError(t, err)

	assert.Equal(t, targetIDs, got)
}

func TestMediaList_MapIDs_Store_error(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	sourceIDs := []entities.SourceID{"1", "2", "3", "5", "8", "13", "999", ""}

	source := test.NewMockSource(t)
	tracker := test.NewMockTracker(t)

	mediaLister := usecases.MediaList{
		Source:  source,
		Store:   nil,
		Tracker: tracker,
	}

	got, err := mediaLister.MapIDs(ctx, sourceIDs)
	require.Error(t, err)

	assert.Nil(t, got)
}

func TestMediaList_Close(t *testing.T) {
	t.Parallel()

	source := test.NewMockSource(t)
	store := test.NewMockStore(t)
	tracker := test.NewMockTracker(t)

	store.EXPECT().Close().Return(nil).Once()
	tracker.EXPECT().Close().Return(nil).Once()

	mediaLister := usecases.MediaList{
		Source:  source,
		Store:   store,
		Tracker: tracker,
	}

	err := mediaLister.Close()
	require.NoError(t, err)
}

func implements[T any](tb testing.TB) any {
	tb.Helper()

	return mock.MatchedBy(func(value any) bool {
		_, ok := value.(T)

		return ok
	})
}
