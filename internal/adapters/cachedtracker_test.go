package adapters_test

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/wwmoraes/anilistarr/internal/adapters"
	"github.com/wwmoraes/anilistarr/internal/test"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

const (
	testUsername = "foo"
	testUserID   = 1
)

var (
	testUserIDStr = fmt.Sprint(testUserID)
	testMedia     = []string{
		"ID1",
		"ID2",
		"ID3",
	}
	testCacheKeyUserID    = fmt.Sprintf("anilist:user:%s:id", testUsername)
	testCacheKeyUserMedia = fmt.Sprintf("anilist:user:%s:media", testUserIDStr)
	testTTL               = adapters.CachedTrackerTTL{
		UserID:       time.Hour,
		MediaListIDs: time.Hour,
	}
)

//nolint:funlen // yeah, tests are long...
func TestCachedTracker_GetUserID(t *testing.T) {
	t.Parallel()

	type fields struct {
		Cache   adapters.Cache
		Tracker usecases.Tracker
		TTL     adapters.CachedTrackerTTL
	}

	type args struct {
		ctx  context.Context
		name string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "cache error",
			fields: fields{
				Cache: (*test.Cache)(nil),
				Tracker: &test.Tracker{
					UserIds: map[string]int{
						testCacheKeyUserID: testUserID,
					},
					MediaLists: map[int][]string{},
				},
				TTL: testTTL,
			},
			args: args{
				ctx:  context.TODO(),
				name: testUsername,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "cache hit",
			fields: fields{
				Cache: &test.Cache{
					Data: map[string]string{
						testCacheKeyUserID: testUserIDStr,
					},
				},
				Tracker: &test.Tracker{
					UserIds: map[string]int{
						testUsername: testUserID,
					},
					MediaLists: map[int][]string{},
				},
				TTL: testTTL,
			},
			args: args{
				ctx:  context.TODO(),
				name: testUsername,
			},
			want:    testUserIDStr,
			wantErr: false,
		},
		{
			name: "cache miss, tracker hit",
			fields: fields{
				Cache: &test.Cache{
					Data: map[string]string{
						testUsername + "qux": fmt.Sprint(testUserID + 1),
					},
				},
				Tracker: &test.Tracker{
					UserIds: map[string]int{
						testUsername: testUserID,
					},
					MediaLists: map[int][]string{},
				},
				TTL: testTTL,
			},
			args: args{
				ctx:  context.TODO(),
				name: testUsername,
			},
			want:    testUserIDStr,
			wantErr: false,
		},
		{
			name: "cache miss, tracker error",
			fields: fields{
				Cache: &test.Cache{
					Data: make(map[string]string),
				},
				Tracker: &test.Tracker{},
				TTL:     testTTL,
			},
			args: args{
				ctx:  context.TODO(),
				name: testUsername,
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			wrapper := &adapters.CachedTracker{
				Cache:   tt.fields.Cache,
				Tracker: tt.fields.Tracker,
				TTL:     tt.fields.TTL,
			}

			got, err := wrapper.GetUserID(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("CachedTracker.GetUserID() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if got != tt.want {
				t.Errorf("CachedTracker.GetUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}

//nolint:funlen // yeah, tests are long...
func TestCachedTracker_GetMediaListIDs(t *testing.T) {
	t.Parallel()

	type fields struct {
		Cache   adapters.Cache
		Tracker usecases.Tracker
		TTL     adapters.CachedTrackerTTL
	}

	type args struct {
		ctx    context.Context
		userId string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "cache error",
			fields: fields{
				Cache:   (*test.Cache)(nil),
				Tracker: &test.Tracker{},
				TTL:     testTTL,
			},
			args: args{
				ctx:    context.TODO(),
				userId: testUserIDStr,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "cache hit",
			fields: fields{
				Cache: &test.Cache{
					Data: map[string]string{
						testCacheKeyUserID:    testUserIDStr,
						testCacheKeyUserMedia: strings.Join(testMedia, "|"),
					},
				},
				Tracker: &test.Tracker{},
				TTL:     testTTL,
			},
			args: args{
				ctx:    context.TODO(),
				userId: testUserIDStr,
			},
			want:    testMedia,
			wantErr: false,
		},
		{
			name: "cache miss, tracker hit",
			fields: fields{
				Cache: &test.Cache{
					Data: map[string]string{
						testCacheKeyUserID: testUserIDStr,
					},
				},
				Tracker: &test.Tracker{
					MediaLists: map[int][]string{
						testUserID: testMedia,
					},
				},
				TTL: testTTL,
			},
			args: args{
				ctx:    context.TODO(),
				userId: testUserIDStr,
			},
			want:    testMedia,
			wantErr: false,
		},
		{
			name: "cache miss, tracker error",
			fields: fields{
				Cache: &test.Cache{
					Data: make(map[string]string),
				},
				Tracker: &test.Tracker{
					MediaLists: make(map[int][]string),
				},
				TTL: testTTL,
			},
			args: args{
				ctx:    context.TODO(),
				userId: "a",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			wrapper := &adapters.CachedTracker{
				Cache:   tt.fields.Cache,
				Tracker: tt.fields.Tracker,
				TTL:     tt.fields.TTL,
			}

			got, err := wrapper.GetMediaListIDs(tt.args.ctx, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("CachedTracker.GetMediaListIDs() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CachedTracker.GetMediaListIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

//nolint:funlen // yeah, tests are long...
func TestCachedTracker_Cache(t *testing.T) {
	t.Parallel()

	cache := &test.Cache{
		Data: map[string]string{},
	}

	tracker := adapters.CachedTracker{
		Cache: cache,
		Tracker: &test.Tracker{
			UserIds: map[string]int{
				testUsername: testUserID,
			},
			MediaLists: map[int][]string{
				testUserID: testMedia,
			},
		},
	}

	gotUserID, err := tracker.GetUserID(context.TODO(), testUsername)
	if err != nil {
		t.Errorf("unexpected CachedTracker.GetUserID() error: %v", err)

		return
	}

	if gotUserID != testUserIDStr {
		t.Errorf("CachedTracker.GetUserID() = %v, want %v", gotUserID, testUserIDStr)

		return
	}

	gotUserID, err = cache.GetString(context.TODO(), testCacheKeyUserID)
	if err != nil {
		t.Errorf("unexpected Cache.GetString() error: %v", err)

		return
	}

	if gotUserID != testUserIDStr {
		t.Errorf("Cache.GetString() = %v, want %v", gotUserID, testUserIDStr)

		return
	}

	gotMediaListIDs, err := tracker.GetMediaListIDs(context.TODO(), testUserIDStr)
	if err != nil {
		t.Errorf("unexpected CachedTracker.GetMediaListIDs() error: %v", err)

		return
	}

	if !reflect.DeepEqual(gotMediaListIDs, testMedia) {
		t.Errorf("CachedTracker.GetMediaListIDs() = %v, want %v", gotMediaListIDs, testMedia)

		return
	}

	gotMediaListIDsStr, err := cache.GetString(context.TODO(), testCacheKeyUserMedia)
	if err != nil {
		t.Errorf("unexpected Cache.GetString() error: %v", err)

		return
	}

	gotMediaListIDs = strings.Split(gotMediaListIDsStr, "|")

	if !reflect.DeepEqual(gotMediaListIDs, testMedia) {
		t.Errorf("CachedTracker.GetMediaListIDs() = %v, want %v", gotMediaListIDs, testMedia)

		return
	}
}
