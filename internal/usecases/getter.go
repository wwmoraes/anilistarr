package usecases

import "net/http"

// Getter provides a way to retrieve data through HTTP GET requests
type Getter interface {
	Get(string) (*http.Response, error)
}
