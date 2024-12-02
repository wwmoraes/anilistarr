package sources_test

import (
	"context"
	"testing"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/wwmoraes/anilistarr/internal/adapters/sources"
	"github.com/wwmoraes/anilistarr/internal/testdata"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

func TestJSONProvider(t *testing.T) {
	t.Parallel()

	testURI := "mem://test"

	testMetadata := []usecases.Metadata{
		testdata.Metadata{
			SourceID: "1",
			TargetID: "91",
		},
		testdata.Metadata{
			SourceID: "2",
			TargetID: "92",
		},
	}

	testData, err := json.Marshal(testMetadata)
	if err != nil {
		t.Error(err)
	}

	provider := sources.JSON[testdata.Metadata](testURI)

	getter := testdata.MockGetter{}

	getter.On("Get", mock.Anything, testURI).
		Return(testData, nil).Once()

	gotURL := provider.String()
	assert.Equal(t, testURI, gotURL)

	gotMetadata, err := provider.Fetch(context.TODO(), &getter)
	require.NoError(t, err)

	assert.Equal(t, testMetadata, gotMetadata)
	getter.AssertExpectations(t)
}

func TestJSONProvider_nilGetter(t *testing.T) {
	t.Parallel()

	provider := sources.JSON[testdata.Metadata]("")

	got, err := provider.Fetch(context.Background(), nil)
	require.ErrorIs(t, err, usecases.ErrStatusInternal)

	assert.Nil(t, got)
}

func TestJSONProvider_notFound(t *testing.T) {
	t.Parallel()

	var data []byte

	testURI := "mem://test"

	getter := testdata.MockGetter{}

	getter.On("Get", mock.Anything, testURI).
		Return(data, usecases.ErrStatusNotFound).Once()

	provider := sources.JSON[testdata.Metadata](testURI)

	got, err := provider.Fetch(context.Background(), &getter)
	require.ErrorIs(t, err, usecases.ErrStatusNotFound)

	assert.Nil(t, got)
	getter.AssertExpectations(t)
}

func TestJSONProvider_invalid(t *testing.T) {
	t.Parallel()

	testURI := "mem://test"
	testData := []byte("test")

	getter := testdata.MockGetter{}
	getter.On("Get", mock.Anything, testURI).
		Return(testData, nil).Once()

	provider := sources.JSON[testdata.Metadata](testURI)

	got, err := provider.Fetch(context.TODO(), &getter)
	require.ErrorIs(t, err, usecases.ErrStatusFailedPrecondition)

	assert.Nil(t, got)
	getter.AssertExpectations(t)
}
