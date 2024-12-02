package testdata

import (
	"net/http"

	"github.com/stretchr/testify/mock"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

var _ usecases.Doer = (*MockDoer)(nil)

type MockDoer struct {
	mock.Mock
}

func (doer *MockDoer) Do(req *http.Request) (*http.Response, error) {
	args := doer.Called(req)

	return args.Get(0).(*http.Response), args.Error(1)
}
