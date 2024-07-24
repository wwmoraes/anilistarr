package adapters_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/goccy/go-json"

	"github.com/wwmoraes/anilistarr/internal/adapters"
	"github.com/wwmoraes/anilistarr/internal/test"
	"github.com/wwmoraes/anilistarr/pkg/functional"
)

const (
	testLocalName = "test.json"
	testURI       = "mem://test"
)

var (
	testMetadata = []adapters.Metadata{
		memoryMetadata{
			AnilistID: 1,
			TheTvdbID: 91,
		},
		memoryMetadata{
			AnilistID: 2,
			TheTvdbID: 92,
		},
	}

	testSourceIDs = []string{testMetadata[0].GetSourceID()}
	testTargetIDs = []string{testMetadata[0].GetTargetID()}
)

type memoryGetter map[string]*http.Response

func (getter memoryGetter) Get(uri string) (*http.Response, error) {
	res, ok := getter[uri]
	if !ok {
		return nil, errors.New(http.StatusText(http.StatusNotFound))
	}

	return res, nil
}

func newJSONLocalProvider(tb testing.TB) *adapters.JSONLocalProvider[memoryMetadata] {
	tb.Helper()

	return &adapters.JSONLocalProvider[memoryMetadata]{
		Fs: &test.MemoryFS{
			testLocalName: functional.Unwrap(json.Marshal(testMetadata)),
		},
		Name: testLocalName,
	}
}

func newMemoryGetter(tb testing.TB) *memoryGetter {
	tb.Helper()

	testData, err := json.Marshal(testMetadata)
	if err != nil {
		tb.Fatal(err)
	}

	return &memoryGetter{
		testURI: &http.Response{
			StatusCode: http.StatusOK,
			Status:     http.StatusText(http.StatusOK),
			Header: http.Header{
				"Content-Type":   []string{"application/json"},
				"Content-Length": []string{strconv.Itoa(len(testData))},
			},
			Body: io.NopCloser(bytes.NewReader(testData)),
		},
	}
}
