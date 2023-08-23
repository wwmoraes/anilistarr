package test

import (
	"bytes"
	"io"
	"net/http"
)

type HTTPClient struct {
	Data map[string]string
}

func (client *HTTPClient) Get(url string) (*http.Response, error) {
	data, ok := client.Data[url]
	if !ok {
		return &http.Response{
			Status:     http.StatusText(http.StatusNotFound),
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(bytes.NewBuffer(nil)),
		}, nil
	}

	return &http.Response{
		Status:     http.StatusText(http.StatusOK),
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(data)),
	}, nil
}
