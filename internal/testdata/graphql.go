package testdata

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
)

var _ http.Handler = (*MockGraphqlHandler)(nil)
var _ http.RoundTripper = (*MockHTTPRoundTripper)(nil)

type MockGraphqlHandler struct {
	mock.Mock
}

func (handler *MockGraphqlHandler) AddRequestResponse(tb testing.TB, req, res any) *mock.Call {
	reqBody, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}

	resBody, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}

	handler.TestData().Set(string(reqBody), resBody)

	return handler.On(
		"ServeHTTP",
		implements[http.ResponseWriter](),
		httpRequestWithBody(reqBody),
	)
}

func (handler *MockGraphqlHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler.Called(w, req)
	defer req.Body.Close()

	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	if !handler.TestData().Has(string(reqBody)) {
		w.WriteHeader(http.StatusNotFound)

		return
	}

	res := handler.TestData().Get(string(reqBody)).Data().([]byte)

	_, err = w.Write(res)
	if err != nil {
		panic(err)
	}
}

type MockHTTPRoundTripper struct {
	mock.Mock
}

func (transport *MockHTTPRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	args := transport.Called(req)

	return args.Get(0).(*http.Response), args.Error(1)
}
