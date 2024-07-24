package api

import (
	"fmt"
	"net/http"

	"github.com/goccy/go-json"
	telemetry "github.com/wwmoraes/gotell"

	"github.com/wwmoraes/anilistarr/internal/usecases"
)

type apiService struct {
	Unimplemented

	mediaLister usecases.MediaLister
}

func NewService(mediaLister usecases.MediaLister) (ServerInterface, error) {
	if mediaLister == nil {
		return nil, fmt.Errorf("Server needs a valid Media Lister instance")
	}

	return &apiService{
		mediaLister: mediaLister,
	}, nil
}

func (service *apiService) GetUserID(w http.ResponseWriter, r *http.Request, name string) {
	span := telemetry.SpanFromContext(r.Context())

	userId, err := service.mediaLister.GetUserID(r.Context(), name)
	if err != nil {
		span.RecordError(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Add("X-Anilist-User-Name", name)
	w.Header().Add("X-Anilist-User-Id", userId)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, userId)
}

func (service *apiService) GetUserMedia(w http.ResponseWriter, r *http.Request, name string) {
	span := telemetry.SpanFromContext(r.Context())

	customList, err := service.mediaLister.Generate(r.Context(), name)
	if err != nil {
		span.RecordError(err)
		http.Error(w, err.Error(), http.StatusBadGateway)

		return
	}

	data, err := json.Marshal(customList)
	if err != nil {
		span.RecordError(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	span.RecordError(err)
}
