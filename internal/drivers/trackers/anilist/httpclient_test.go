package anilist_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"

	"github.com/wwmoraes/anilistarr/internal/drivers/trackers/anilist"
	"github.com/wwmoraes/anilistarr/internal/testdata"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

func TestRatedClient_Do(t *testing.T) {
	t.Parallel()

	t.Skip()

	type fields struct {
		Doer    usecases.Doer
		Limiter *rate.Limiter
	}

	type args struct {
		req *http.Request
	}

	tests := []struct {
		fields    fields
		assertion assert.ErrorAssertionFunc
		want      *http.Response
		args      args
		name      string
	}{
		{
			name: "200",
			fields: fields{
				Doer:    http.DefaultClient,
				Limiter: rate.NewLimiter(rate.Limit(time.Nanosecond), 1),
			},
			args: args{
				req: &http.Request{},
			},
			assertion: assert.NoError,
			want:      &http.Response{},
		},
		{
			name: "429",
			fields: fields{
				Doer:    http.DefaultClient,
				Limiter: rate.NewLimiter(rate.Every(0), 0),
			},
			args: args{
				req: &http.Request{},
			},
			assertion: assert.Error,
			want:      &http.Response{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := anilist.RatedClient{
				Doer:    tt.fields.Doer,
				Limiter: tt.fields.Limiter,
			}

			got, err := client.Do(tt.args.req)
			tt.assertion(t, err)

			err = got.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRatedClient_200(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(
		http.MethodGet,
		"http://example.com/foo",
		http.NoBody,
	)

	recorder := httptest.NewRecorder()
	recorder.Header().Set("X-Ratelimit-Limit", "1")
	recorder.Header().Set("X-Ratelimit-Remaining", "0")
	recorder.WriteHeader(http.StatusOK)
	_, err := recorder.WriteString("bar")
	require.NoError(t, err)

	res := recorder.Result()
	res.Request = req
	defer res.Body.Close()

	doer := testdata.MockDoer{}

	doer.On("Do", req).Return(res, nil).Once()

	limiter := rate.NewLimiter(rate.Every(time.Nanosecond), 1)

	client := anilist.RatedClient{
		Doer:    &doer,
		Limiter: limiter,
	}

	got, err := client.Do(req)
	require.NoError(t, err)
	defer got.Body.Close()

	assert.Equal(t, res, got)

	doer.AssertExpectations(t)
}

func TestRatedClient_local_429(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(
		http.MethodGet,
		"http://example.com/foo",
		http.NoBody,
	)

	doer := testdata.MockDoer{}

	limiter := rate.NewLimiter(rate.Every(time.Hour), 1)

	client := anilist.RatedClient{
		Doer:    &doer,
		Limiter: limiter,
	}

	// consume limit before call
	limiter.SetBurst(0)
	limiter.SetBurstAt(time.Now().Add(time.Hour), 1)

	got, err := client.Do(req)
	require.NoError(t, err)
	defer got.Body.Close()

	assert.Equal(t, http.StatusTooManyRequests, got.StatusCode)
	assert.Equal(t, "3600", got.Header.Get("Retry-After"))

	doer.AssertExpectations(t)
}

func TestRatedClient_remote_429(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(
		http.MethodGet,
		"http://example.com/foo",
		http.NoBody,
	)

	retryAt := time.Now().Add(time.Hour)

	recorder := httptest.NewRecorder()
	recorder.Header().Set("X-Ratelimit-Remaining", "0")
	recorder.Header().Set("X-Ratelimit-Limit", "2")
	recorder.Header().Set("X-Ratelimit-Reset", strconv.FormatInt(retryAt.Unix(), 10))
	recorder.Header().Set("Retry-After", "3600")
	recorder.WriteHeader(http.StatusTooManyRequests)

	res := recorder.Result()
	res.Request = req
	defer res.Body.Close()

	doer := testdata.MockDoer{}

	doer.On("Do", req).Return(res, nil).Once()

	limiter := rate.NewLimiter(rate.Every(time.Hour), 1)

	client := anilist.RatedClient{
		Doer:    &doer,
		Limiter: limiter,
	}

	got, err := client.Do(req)
	require.NoError(t, err)
	defer got.Body.Close()

	assert.Equal(t, res, got)
	assert.Equal(t, 2, limiter.Burst())

	doer.AssertExpectations(t)
}
