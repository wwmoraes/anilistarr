// Package anilist provides an [usecases.Tracker] that communicates with
// Anilist's GraphQL API.
package anilist_test

import (
	"context"
	"net/http"
	"strconv"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/wwmoraes/anilistarr/internal/drivers/trackers/anilist"
	"github.com/wwmoraes/anilistarr/internal/testdata"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

func TestTracker_GetMediaListIDs(t *testing.T) {
	t.Parallel()

	type fields struct {
		Client   graphql.Client
		PageSize int
	}

	type args struct {
		ctx    context.Context
		userID string
	}

	tests := []struct {
		assertion assert.ErrorAssertionFunc
		args      args
		fields    fields
		name      string
		want      []string
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tracker := &anilist.Tracker{
				Client:   tt.fields.Client,
				PageSize: tt.fields.PageSize,
			}

			got, err := tracker.GetMediaListIDs(tt.args.ctx, tt.args.userID)
			tt.assertion(t, err)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTracker_Close(t *testing.T) {
	t.Parallel()

	type fields struct {
		Client   graphql.Client
		PageSize int
	}

	tests := []struct {
		assertion assert.ErrorAssertionFunc
		fields    fields
		name      string
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tracker := &anilist.Tracker{
				Client:   tt.fields.Client,
				PageSize: tt.fields.PageSize,
			}

			tt.assertion(t, tracker.Close())
		})
	}
}

func TestTracker(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	username := "foo"
	userID := 1
	mediaList := []string{"11"}

	transport := testdata.MockHTTPRoundTripper{}

	transport.On(
		"RoundTrip",
		testdata.HTTPRequestWithJSONBody(t, graphql.Request{
			OpName: "GetUserByName",
			Query:  anilist.GetUserByName_Operation,
			Variables: &struct {
				Name string `json:"name"`
			}{
				Name: username,
			},
		}),
	).Return(
		//nolint:bodyclose // client transport closes it
		testdata.HTTPResponseWithJSONBody(t, graphql.Response{
			Data: &anilist.GetUserByNameResponse{
				User: anilist.GetUserByNameUser{
					Id: userID,
				},
			},
		}),
		nil,
	).Once()

	transport.On(
		"RoundTrip",
		testdata.HTTPRequestWithJSONBody(t, graphql.Request{
			OpName: "GetWatching",
			Query:  anilist.GetWatching_Operation,
			Variables: &struct {
				//nolint:tagliatelle // upstream format
				UserID  int `json:"userId"`
				Page    int `json:"page"`
				PerPage int `json:"perPage"`
			}{
				UserID:  userID,
				Page:    1,
				PerPage: 10,
			},
		}),
	).Return(
		//nolint:bodyclose // client transport closes it
		testdata.HTTPResponseWithJSONBody(t, graphql.Response{
			Data: &anilist.GetWatchingResponse{
				Page: anilist.GetWatchingPage{
					MediaList: []anilist.GetWatchingPageMediaList{
						{
							Media: anilist.GetWatchingPageMediaListMedia{
								Id: 11,
							},
						},
					},
				},
			},
		}),
		nil,
	).Once()

	client := anilist.New(
		"http://example.com",
		anilist.WithClient(&http.Client{
			Transport: &transport,
		}),
		anilist.WithPageSize(10),
	)
	defer client.Close()

	gotUserID, err := client.GetUserID(ctx, username)
	require.NoError(t, err)

	assert.Equal(t, strconv.Itoa(userID), gotUserID)

	gotMediaList, err := client.GetMediaListIDs(ctx, strconv.Itoa(userID))
	require.NoError(t, err)

	assert.Equal(t, mediaList, gotMediaList)
	transport.AssertExpectations(t)
}

func TestTracker_GetUserID_not_found(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	username := "foo"

	transport := testdata.MockHTTPRoundTripper{}

	transport.On(
		"RoundTrip",
		testdata.HTTPRequestWithJSONBody(t, graphql.Request{
			OpName: "GetUserByName",
			Query:  anilist.GetUserByName_Operation,
			Variables: &struct {
				Name string `json:"name"`
			}{
				Name: username,
			},
		}),
	).Return(
		//nolint:bodyclose // client transport closes it
		testdata.HTTPResponseWithJSONBody(t, graphql.Response{
			Data: &anilist.GetUserByNameResponse{
				User: anilist.GetUserByNameUser{},
			},
			Errors: gqlerror.List{
				&gqlerror.Error{
					Message: http.StatusText(http.StatusNotFound) + ".",
					Extensions: map[string]any{
						"status": http.StatusNotFound,
					},
					Locations: []gqlerror.Location{
						{
							Line:   2,
							Column: 2,
						},
					},
				},
			},
		}),
		nil,
	).Once()

	client := anilist.New(
		"http://example.com",
		anilist.WithClient(&http.Client{
			Transport: &transport,
		}),
		anilist.WithPageSize(10),
	)
	defer client.Close()

	got, err := client.GetUserID(ctx, username)
	require.ErrorIs(t, err, usecases.ErrStatusNotFound)

	assert.Empty(t, got)

	transport.AssertExpectations(t)
}

func TestTracker_GetUserID_unavailable(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	username := "foo"

	transport := testdata.MockHTTPRoundTripper{}

	transport.On(
		"RoundTrip",
		testdata.HTTPRequestWithJSONBody(t, graphql.Request{
			OpName: "GetUserByName",
			Query:  anilist.GetUserByName_Operation,
			Variables: &struct {
				Name string `json:"name"`
			}{
				Name: username,
			},
		}),
	).Return(
		//nolint:bodyclose // client transport closes it
		testdata.HTTPResponseWithJSONBody(t, graphql.Response{
			Data: &anilist.GetUserByNameResponse{
				User: anilist.GetUserByNameUser{},
			},
			Errors: gqlerror.List{
				&gqlerror.Error{
					Message: http.StatusText(http.StatusInternalServerError) + ".",
					Extensions: map[string]any{
						"status": http.StatusInternalServerError,
					},
					Locations: []gqlerror.Location{
						{
							Line:   2,
							Column: 2,
						},
					},
				},
			},
		}),
		nil,
	).Once()

	client := anilist.New(
		"http://example.com",
		anilist.WithClient(&http.Client{
			Transport: &transport,
		}),
		anilist.WithPageSize(10),
	)
	defer client.Close()

	got, err := client.GetUserID(ctx, username)
	require.ErrorIs(t, err, usecases.ErrStatusUnavailable)

	assert.Empty(t, got)

	transport.AssertExpectations(t)
}
