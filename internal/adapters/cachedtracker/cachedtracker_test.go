package cachedtracker_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/wwmoraes/anilistarr/internal/adapters/cachedtracker"
	"github.com/wwmoraes/anilistarr/internal/test"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

func TestCachedTracker_GetUserID_uncached(t *testing.T) {
	t.Parallel()

	username := "foo"
	userID := "1"
	cacheKeyUserID := "anilist:user:foo:id"

	cache := test.NewMockCache(t)
	tracker := test.NewMockTracker(t)

	// cache miss for user ID
	cache.EXPECT().GetString(
		mock.Anything,
		cacheKeyUserID,
	).Return("", usecases.ErrStatusNotFound).Once()
	// store user ID in the cache
	cache.EXPECT().SetString(
		mock.Anything,
		cacheKeyUserID,
		userID,
		mock.Anything,
	).Return(nil).Once()

	tracker.EXPECT().GetUserID(
		mock.Anything,
		username,
	).Return(userID, nil).Once()

	cachedTracker := cachedtracker.CachedTracker{
		Cache:   cache,
		Tracker: tracker,
	}

	gotUserID, err := cachedTracker.GetUserID(t.Context(), username)
	require.NoError(t, err)

	assert.Equal(t, userID, gotUserID)
}

func TestCachedTracker_GetUserID_cached(t *testing.T) {
	t.Parallel()

	username := "foo"
	userID := "1"
	cacheKeyUserID := "anilist:user:foo:id"

	cache := test.NewMockCache(t)
	tracker := test.NewMockTracker(t)

	// cache miss for user ID
	cache.EXPECT().GetString(
		mock.Anything,
		cacheKeyUserID,
	).Return(userID, nil).Once()

	cachedTracker := cachedtracker.CachedTracker{
		Cache:   cache,
		Tracker: tracker,
	}

	gotUserID, err := cachedTracker.GetUserID(t.Context(), username)
	require.NoError(t, err)

	assert.Equal(t, userID, gotUserID)
}

func TestCachedTracker_GetMediaListIDs_uncached(t *testing.T) {
	t.Parallel()

	userID := "1"
	medias := []string{
		"ID1",
		"ID2",
		"ID3",
	}
	mediasStr := strings.Join(medias, "|")
	cacheKeyUserMedia := "anilist:user:1:media"

	cache := test.NewMockCache(t)
	tracker := test.NewMockTracker(t)

	// cache miss for user medias
	cache.EXPECT().GetString(
		mock.Anything,
		cacheKeyUserMedia,
	).Return("", usecases.ErrStatusNotFound).Once()
	// store user medias in the cache
	cache.EXPECT().SetString(
		mock.Anything,
		cacheKeyUserMedia,
		mediasStr,
		mock.Anything,
	).Return(nil).Once()
	// retrieve cached medias
	tracker.EXPECT().GetMediaListIDs(
		mock.Anything,
		userID,
	).Return(medias, nil).Once()

	cachedTracker := cachedtracker.CachedTracker{
		Cache:   cache,
		Tracker: tracker,
	}

	gotMediaListIDs, err := cachedTracker.GetMediaListIDs(t.Context(), userID)
	require.NoError(t, err)

	assert.Equal(t, medias, gotMediaListIDs)
}

func TestCachedTracker_GetMediaListIDs_cached(t *testing.T) {
	t.Parallel()

	userID := "1"
	medias := []string{
		"ID1",
		"ID2",
		"ID3",
	}
	mediasStr := strings.Join(medias, "|")
	cacheKeyUserMedia := "anilist:user:1:media"

	cache := test.NewMockCache(t)
	tracker := test.NewMockTracker(t)

	// retrieve cached medias
	cache.EXPECT().GetString(
		mock.Anything,
		cacheKeyUserMedia,
	).Return(mediasStr, nil).Once()

	cachedTracker := cachedtracker.CachedTracker{
		Cache:   cache,
		Tracker: tracker,
	}

	gotMediaListIDs, err := cachedTracker.GetMediaListIDs(t.Context(), userID)
	require.NoError(t, err)

	assert.Equal(t, medias, gotMediaListIDs)
}

func TestCachedTracker_Close(t *testing.T) {
	t.Parallel()

	cache := test.NewMockCache(t)
	tracker := test.NewMockTracker(t)

	cache.EXPECT().Close().Return(nil).Once()
	tracker.EXPECT().Close().Return(nil).Once()

	cachedTracker := cachedtracker.CachedTracker{
		Cache:   cache,
		Tracker: tracker,
	}

	err := cachedTracker.Close()
	require.NoError(t, err)
}

func TestCachedTracker_Cache_error(t *testing.T) {
	t.Parallel()

	cacheError := errors.New("qux")

	cache := test.NewMockCache(t)
	tracker := test.NewMockTracker(t)

	cache.EXPECT().GetString(
		mock.Anything,
		"anilist:user:foo:id",
	).Return("", cacheError).Once()
	cache.EXPECT().GetString(
		mock.Anything,
		"anilist:user:1:media",
	).Return("", cacheError).Once()

	cachedTracker := cachedtracker.CachedTracker{
		Cache:   cache,
		Tracker: tracker,
	}

	gotUserID, err := cachedTracker.GetUserID(t.Context(), "foo")
	require.ErrorIs(t, err, cacheError)

	assert.Empty(t, gotUserID)

	gotMedias, err := cachedTracker.GetMediaListIDs(t.Context(), "1")
	require.ErrorIs(t, err, cacheError)

	assert.Nil(t, gotMedias)
}

func TestCachedTracker_Tracker_error(t *testing.T) {
	t.Parallel()

	var medias []string

	trackerError := errors.New("qux")

	cache := test.NewMockCache(t)
	tracker := test.NewMockTracker(t)

	cache.EXPECT().GetString(
		mock.Anything,
		"anilist:user:foo:id",
	).Return("", usecases.ErrStatusNotFound).Once()
	cache.EXPECT().GetString(
		mock.Anything,
		"anilist:user:1:media",
	).Return("", usecases.ErrStatusNotFound).Once()
	tracker.EXPECT().GetUserID(
		mock.Anything,
		"foo",
	).Return("", trackerError).Once()
	tracker.EXPECT().GetMediaListIDs(
		mock.Anything,
		"1",
	).Return(medias, trackerError).Once()

	cachedTracker := cachedtracker.CachedTracker{
		Cache:   cache,
		Tracker: tracker,
	}

	gotUserID, err := cachedTracker.GetUserID(t.Context(), "foo")
	require.ErrorIs(t, err, trackerError)

	assert.Empty(t, gotUserID)

	gotMedias, err := cachedTracker.GetMediaListIDs(t.Context(), "1")
	require.ErrorIs(t, err, trackerError)

	assert.Nil(t, gotMedias)
}
