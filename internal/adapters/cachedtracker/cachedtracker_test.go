package cachedtracker_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/wwmoraes/anilistarr/internal/adapters/cachedtracker"
	"github.com/wwmoraes/anilistarr/internal/testdata"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

func TestCachedTracker_GetUserID_uncached(t *testing.T) {
	t.Parallel()

	username := "foo"
	userID := "1"
	cacheKeyUserID := "anilist:user:foo:id"

	cache := testdata.MockCache{}
	tracker := testdata.MockTracker{}

	// cache miss for user ID
	cache.On(
		"GetString",
		mock.Anything,
		cacheKeyUserID,
	).Return("", usecases.ErrStatusNotFound).Once()
	// store user ID in the cache
	cache.On(
		"SetString",
		mock.Anything,
		cacheKeyUserID,
		userID,
		mock.Anything,
	).Return(nil).Once()

	tracker.On("GetUserID", mock.Anything, username).
		Return(userID, nil).Once()

	cachedTracker := cachedtracker.CachedTracker{
		Cache:   &cache,
		Tracker: &tracker,
	}

	gotUserID, err := cachedTracker.GetUserID(context.TODO(), username)
	require.NoError(t, err)

	assert.Equal(t, userID, gotUserID)

	cache.AssertExpectations(t)
	tracker.AssertExpectations(t)
}

func TestCachedTracker_GetUserID_cached(t *testing.T) {
	t.Parallel()

	username := "foo"
	userID := "1"
	cacheKeyUserID := "anilist:user:foo:id"

	cache := testdata.MockCache{}
	tracker := testdata.MockTracker{}

	// cache miss for user ID
	cache.On(
		"GetString",
		mock.Anything,
		cacheKeyUserID,
	).Return(userID, nil).Once()

	cachedTracker := cachedtracker.CachedTracker{
		Cache:   &cache,
		Tracker: &tracker,
	}

	gotUserID, err := cachedTracker.GetUserID(context.TODO(), username)
	require.NoError(t, err)

	assert.Equal(t, userID, gotUserID)

	cache.AssertExpectations(t)
	tracker.AssertExpectations(t)
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

	cache := testdata.MockCache{}
	tracker := testdata.MockTracker{}

	// cache miss for user medias
	cache.On(
		"GetString",
		mock.Anything,
		cacheKeyUserMedia,
	).Return("", usecases.ErrStatusNotFound).Once()
	// store user medias in the cache
	cache.On(
		"SetString",
		mock.Anything,
		cacheKeyUserMedia,
		mediasStr,
		mock.Anything,
	).Return(nil).Once()
	// retrieve cached medias
	tracker.On("GetMediaListIDs", mock.Anything, userID).
		Return(medias, nil).Once()

	cachedTracker := cachedtracker.CachedTracker{
		Cache:   &cache,
		Tracker: &tracker,
	}

	gotMediaListIDs, err := cachedTracker.GetMediaListIDs(context.TODO(), userID)
	require.NoError(t, err)

	assert.Equal(t, medias, gotMediaListIDs)

	cache.AssertExpectations(t)
	tracker.AssertExpectations(t)
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

	cache := testdata.MockCache{}
	tracker := testdata.MockTracker{}

	// retrieve cached medias
	cache.On(
		"GetString",
		mock.Anything,
		cacheKeyUserMedia,
	).Return(mediasStr, nil).Once()

	cachedTracker := cachedtracker.CachedTracker{
		Cache:   &cache,
		Tracker: &tracker,
	}

	gotMediaListIDs, err := cachedTracker.GetMediaListIDs(context.TODO(), userID)
	require.NoError(t, err)

	assert.Equal(t, medias, gotMediaListIDs)

	cache.AssertExpectations(t)
	tracker.AssertExpectations(t)
}

func TestCachedTracker_Close(t *testing.T) {
	t.Parallel()

	cache := testdata.MockCache{}
	tracker := testdata.MockTracker{}

	cache.On("Close").Return(nil).Once()
	tracker.On("Close").Return(nil).Once()

	cachedTracker := cachedtracker.CachedTracker{
		Cache:   &cache,
		Tracker: &tracker,
	}

	err := cachedTracker.Close()
	require.NoError(t, err)

	cache.AssertExpectations(t)
	tracker.AssertExpectations(t)
}

func TestCachedTracker_Cache_error(t *testing.T) {
	t.Parallel()

	cacheError := errors.New("qux")

	cache := testdata.MockCache{}
	tracker := testdata.MockTracker{}

	cache.On(
		"GetString",
		mock.Anything,
		"anilist:user:foo:id",
	).Return("", cacheError).Once()
	cache.On(
		"GetString",
		mock.Anything,
		"anilist:user:1:media",
	).Return("", cacheError).Once()

	cachedTracker := cachedtracker.CachedTracker{
		Cache:   &cache,
		Tracker: &tracker,
	}

	gotUserID, err := cachedTracker.GetUserID(context.TODO(), "foo")
	require.ErrorIs(t, err, cacheError)

	assert.Empty(t, gotUserID)

	gotMedias, err := cachedTracker.GetMediaListIDs(context.TODO(), "1")
	require.ErrorIs(t, err, cacheError)

	assert.Nil(t, gotMedias)

	cache.AssertExpectations(t)
	tracker.AssertExpectations(t)
}

func TestCachedTracker_Tracker_error(t *testing.T) {
	t.Parallel()

	var medias []string

	trackerError := errors.New("qux")

	cache := testdata.MockCache{}
	tracker := testdata.MockTracker{}

	cache.On(
		"GetString",
		mock.Anything,
		"anilist:user:foo:id",
	).Return("", usecases.ErrStatusNotFound).Once()
	cache.On(
		"GetString",
		mock.Anything,
		"anilist:user:1:media",
	).Return("", usecases.ErrStatusNotFound).Once()
	tracker.On("GetUserID", mock.Anything, "foo").
		Return("", trackerError).Once()
	tracker.On("GetMediaListIDs", mock.Anything, "1").
		Return(medias, trackerError).Once()

	cachedTracker := cachedtracker.CachedTracker{
		Cache:   &cache,
		Tracker: &tracker,
	}

	gotUserID, err := cachedTracker.GetUserID(context.TODO(), "foo")
	require.ErrorIs(t, err, trackerError)

	assert.Empty(t, gotUserID)

	gotMedias, err := cachedTracker.GetMediaListIDs(context.TODO(), "1")
	require.ErrorIs(t, err, trackerError)

	assert.Nil(t, gotMedias)

	cache.AssertExpectations(t)
	tracker.AssertExpectations(t)
}
