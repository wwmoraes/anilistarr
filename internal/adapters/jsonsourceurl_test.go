package adapters_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/wwmoraes/anilistarr/internal/adapters"
)

type mockClient struct {
	data map[string]string
}

func (client *mockClient) Get(url string) (*http.Response, error) {
	data, ok := client.data[url]
	if !ok {
		return &http.Response{
			Status:     http.StatusText(http.StatusNotFound),
			StatusCode: http.StatusNotFound,
		}, nil
	}

	return &http.Response{
		Status:     http.StatusText(http.StatusOK),
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(data)),
	}, nil
}

type mockData struct {
	AnilistID string `json:"anilist_id,omitempty"`
	TvdbID    string `json:"tvdb_id,omitempty"`
}

func (data mockData) GetAnilistID() string {
	return data.AnilistID
}

func (data mockData) GetTvdbID() string {
	return data.TvdbID
}

func TestJSONSourceURL_Fetch(t *testing.T) {
	expectedData := []adapters.Metadata{
		mockData{
			AnilistID: "123",
			TvdbID:    "456",
		},
	}

	bytesData, err := json.Marshal(expectedData)
	if err != nil {
		t.Error(err)
	}

	ctx := context.Background()

	provider := adapters.JSONSourceURL[mockData]("/anime-lists.json")

	metadata, err := provider.Fetch(ctx, &mockClient{
		data: map[string]string{
			"/anime-lists.json": string(bytesData),
		},
	})
	if err != nil {
		t.Error(err)
	}

	if len(metadata) != len(expectedData) {
		t.Errorf("metadata length mismatch: got %d, wanted %d", len(metadata), len(expectedData))
	}

	for index, entry := range metadata {
		if entry.GetAnilistID() != expectedData[index].GetAnilistID() {
			t.Errorf("metadata anilist ID mismatch: got %s, expected %s", entry.GetAnilistID(), expectedData[index].GetAnilistID())
		}

		if entry.GetTvdbID() != expectedData[index].GetTvdbID() {
			t.Errorf("metadata tvdb ID mismatch: got %s, expected %s", entry.GetTvdbID(), expectedData[index].GetTvdbID())
		}
	}
}
