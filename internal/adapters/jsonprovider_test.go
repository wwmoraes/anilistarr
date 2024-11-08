package adapters_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"testing"

	"github.com/goccy/go-json"

	"github.com/wwmoraes/anilistarr/internal/adapters"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

func TestLocalProvider(t *testing.T) {
	t.Parallel()

	testData, err := json.Marshal(testMetadata)
	if err != nil {
		t.Error(err)
	}

	provider := adapters.JSONProvider[memoryMetadata](testURI)

	getter := memoryGetter{
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

	gotURL := provider.String()

	if gotURL != testURI {
		t.Errorf("want %v, got %v", testURI, gotURL)
	}

	gotMetadata, err := provider.Fetch(context.TODO(), usecases.HTTPGetterAsGetter(getter))
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if !reflect.DeepEqual(gotMetadata, testMetadata) {
		t.Errorf("want %v, got %v", testMetadata, gotMetadata)
	}
}

func TestLocalProvider_nilGetter(t *testing.T) {
	t.Parallel()

	provider := adapters.JSONProvider[memoryMetadata]("")

	gotValue, gotErr := provider.Fetch(context.Background(), nil)

	if !errors.Is(gotErr, adapters.ErrNoGetter) {
		t.Errorf("got error %v, want %v", gotErr, adapters.ErrNoGetter)
	}

	if gotValue != nil {
		t.Errorf("unexpected value %v", gotValue)
	}
}

func TestLocalProvider_notFound(t *testing.T) {
	t.Parallel()

	testURI := "mem://test"

	provider := adapters.JSONProvider[memoryMetadata](testURI)

	gotValue, gotErr := provider.Fetch(context.Background(), usecases.HTTPGetterAsGetter(&memoryGetter{}))
	if gotErr == nil {
		t.Error("got no error, expected some")
	}

	if gotValue != nil {
		t.Errorf("got value %v, expected nil", gotValue)
	}
}

func TestLocalProvider_invalid(t *testing.T) {
	t.Parallel()

	testURI := "mem://test"
	testData := []byte("test")

	provider := adapters.JSONProvider[memoryMetadata](testURI)

	gotValue, gotErr := provider.Fetch(context.TODO(), usecases.HTTPGetterAsGetter(&memoryGetter{
		testURI: &http.Response{
			StatusCode: http.StatusOK,
			Status:     http.StatusText(http.StatusOK),
			Header: http.Header{
				"Content-Type":   []string{"application/json"},
				"Content-Length": []string{strconv.Itoa(len(testData))},
			},
			Body: io.NopCloser(bytes.NewReader(testData)),
		},
	}))
	if gotErr == nil {
		t.Error("got no error, expected some")
	}

	if gotValue != nil {
		t.Errorf("got value %v, expected nil", gotValue)
	}
}
