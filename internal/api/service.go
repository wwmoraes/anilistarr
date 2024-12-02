package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/goccy/go-json"
	telemetry "github.com/wwmoraes/gotell"

	"github.com/wwmoraes/anilistarr/internal/usecases"
)

var _ ServerInterface = (*Service)(nil)

// Service implements handlers to serve media lister as a REST API.
type Service struct {
	Unimplemented

	MediaLister usecases.MediaLister
}

// GetUserID retrieves an user ID for a given name. Responds with:
//   - 200 + plain-text ID + headers on success
//   - 404 if media lister cannot find the user
//   - 500 for any other errors
func (service *Service) GetUserID(w http.ResponseWriter, r *http.Request, name string) {
	span := telemetry.SpanFromContext(r.Context())

	userID, err := service.MediaLister.GetUserID(r.Context(), name)
	if errors.Is(err, usecases.ErrStatusNotFound) {
		span.RecordError(err)
		http.Error(w, err.Error(), http.StatusNotFound)

		return
	}

	if err != nil {
		span.RecordError(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Add("X-Anilist-User-Name", name)
	w.Header().Add("X-Anilist-User-Id", userID)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, userID)
}

// GetUserMedia retrieves media information from an user. Returns 200 on success
// with a marshaled [entities.CustomList] as JSON or a 502 otherwise.
func (service *Service) GetUserMedia(w http.ResponseWriter, r *http.Request, name string) {
	span := telemetry.SpanFromContext(r.Context())

	customList, err := service.MediaLister.Generate(r.Context(), name)
	if err != nil {
		span.RecordError(err)
		http.Error(w, err.Error(), http.StatusBadGateway)

		return
	}

	data, _ := json.Marshal(customList)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	span.RecordError(err)
}
