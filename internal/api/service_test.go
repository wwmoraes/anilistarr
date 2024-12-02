package api_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/wwmoraes/anilistarr/internal/api"
	"github.com/wwmoraes/anilistarr/internal/entities"
	"github.com/wwmoraes/anilistarr/internal/testdata"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

func TestService_GetUserID(t *testing.T) {
	t.Parallel()

	username := "foo"
	userID := "bar"

	ctx := context.TODO()
	mediaLister := testdata.MockMediaLister{}

	mediaLister.
		On("GetUserID", mock.Anything, username).
		Return(userID, nil).
		Once()

	service := api.Service{
		MediaLister: &mediaLister,
	}

	r := httptest.NewRequest(
		http.MethodGet,
		"http://example.com/",
		http.NoBody,
	).WithContext(ctx)
	w := httptest.NewRecorder()

	service.GetUserID(w, r, username)

	res := w.Result()
	defer res.Body.Close()

	gotBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	gotUserID := string(bytes.Trim(gotBody, " \r\n"))

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, userID, gotUserID)
	assert.Subset(t, res.Header, http.Header{
		"X-Anilist-User-Name": []string{username},
		"X-Anilist-User-Id":   []string{userID},
		"Content-Type":        []string{"text/plain; charset=utf-8"},
	})
	mediaLister.AssertExpectations(t)
}

func TestService_GetUserID_error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		wantError   error
		wantHeaders http.Header
		name        string
		wantStatus  int
	}{
		{
			name:       "not found",
			wantError:  usecases.ErrStatusNotFound,
			wantStatus: http.StatusNotFound,
			wantHeaders: http.Header{
				"Content-Type": []string{"text/plain; charset=utf-8"},
			},
		},
		{
			name:       "unknown",
			wantError:  errors.New("bar"),
			wantStatus: http.StatusInternalServerError,
			wantHeaders: http.Header{
				"Content-Type": []string{"text/plain; charset=utf-8"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			username := "foo"
			ctx := context.TODO()

			mediaLister := testdata.MockMediaLister{}

			mediaLister.
				On("GetUserID", mock.Anything, username).
				Return("", tt.wantError).Once()

			service := api.Service{
				MediaLister: &mediaLister,
			}

			r := httptest.NewRequest(
				http.MethodGet,
				"http://example.com/",
				http.NoBody,
			).WithContext(ctx)
			w := httptest.NewRecorder()

			service.GetUserID(w, r, username)

			res := w.Result()
			defer res.Body.Close()

			gotBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			gotMessage := string(bytes.Trim(gotBody, " \r\n"))

			assert.Equal(t, tt.wantStatus, res.StatusCode)
			assert.Subset(t, res.Header, tt.wantHeaders)
			assert.Equal(t, tt.wantError.Error(), gotMessage)
			mediaLister.AssertExpectations(t)
		})
	}
}

func TestService_GetUserMedia(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	username := "foo"
	medias := entities.CustomList{
		entities.CustomEntry{
			TvdbID: 91,
		},
	}

	mediaLister := testdata.MockMediaLister{}

	mediaLister.
		On("Generate", mock.Anything, username).
		Return(medias, nil).
		Once()

	r := httptest.NewRequest(
		http.MethodGet,
		"http://example.com/",
		http.NoBody,
	).WithContext(ctx)
	resWriter := httptest.NewRecorder()

	service := api.Service{
		MediaLister: &mediaLister,
	}

	service.GetUserMedia(resWriter, r, username)

	res := resWriter.Result()
	defer res.Body.Close()

	gotBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	var gotMedias entities.CustomList

	err = json.Unmarshal(gotBody, &gotMedias)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Subset(t, res.Header, http.Header{
		"Content-Type": []string{"application/json"},
	})
	assert.Equal(t, medias, gotMedias)
	mediaLister.AssertExpectations(t)
}

func TestService_GetUserMedia_error(t *testing.T) {
	t.Parallel()

	username := "foo"
	ctx := context.TODO()
	wantErr := errors.New("bar")

	mediaLister := testdata.MockMediaLister{}

	mediaLister.
		On("Generate", mock.Anything, username).
		Return(
			entities.CustomList(nil),
			wantErr,
		).Once()

	r := httptest.NewRequest(
		http.MethodGet,
		"http://example.com/",
		http.NoBody,
	).WithContext(ctx)
	resWriter := httptest.NewRecorder()

	service := api.Service{
		MediaLister: &mediaLister,
	}

	service.GetUserMedia(resWriter, r, username)

	res := resWriter.Result()
	defer res.Body.Close()

	gotBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	gotMessage := string(bytes.Trim(gotBody, " \r\n"))

	assert.Equal(t, http.StatusBadGateway, res.StatusCode)
	assert.Equal(t, wantErr.Error(), gotMessage)
}
