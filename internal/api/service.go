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

type Service struct {
	Unimplemented

	mediaLister usecases.MediaLister
}

func NewService(mediaLister usecases.MediaLister) (*Service, error) {
	if mediaLister == nil {
		return nil, errors.New("server needs a valid media lister instance")
	}

	return &Service{
		mediaLister: mediaLister,
	}, nil
}

func (service *Service) GetUserID(w http.ResponseWriter, r *http.Request, name string) {
	span := telemetry.SpanFromContext(r.Context())

	userID, err := service.mediaLister.GetUserID(r.Context(), name)
	if errors.Is(err, usecases.ErrNotFound) {
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

func (service *Service) GetUserMedia(w http.ResponseWriter, r *http.Request, name string) {
	span := telemetry.SpanFromContext(r.Context())

	customList, err := service.mediaLister.Generate(r.Context(), name)
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
