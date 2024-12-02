package main

import (
	"bytes"
	"io"
	"net/http"

	"github.com/wwmoraes/anilistarr/internal/usecases"
)

type httpClient struct {
	Data map[string]string
}

// Do handles HTTP requests, simulating GET and ignoring all other methods.
func (client *httpClient) Do(req *http.Request) (*http.Response, error) {
	switch req.Method {
	case http.MethodGet:
		return client.Get(req.URL.String())
	default:
		return nil, usecases.ErrStatusUnimplemented
	}
}

// Get simulates a network round-trip to handle an HTTP GET request.
func (client *httpClient) Get(uri string) (*http.Response, error) {
	data, ok := client.Data[uri]
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
