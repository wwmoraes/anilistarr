package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/wwmoraes/anilistarr/internal/entities"
	"github.com/wwmoraes/anilistarr/internal/telemetry"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

var (
	usernameRegex = regexp.MustCompile("^[[:word:]]+$")
)

type RestAPI struct {
	mediaLister usecases.MediaLister
}

func (face *RestAPI) GetList(w http.ResponseWriter, r *http.Request) {
	span := telemetry.SpanFromContext(r.Context())

	username := r.URL.Query().Get("username")
	if !usernameRegex.MatchString(username) {
		err := fmt.Errorf("invalid username")
		http.Error(w, err.Error(), http.StatusBadRequest)
		span.RecordError(err)
		return
	}

	customList, err := face.mediaLister.Generate(r.Context(), username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		span.RecordError(err)
		return
	}

	data, err := json.Marshal(customList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		span.RecordError(err)
		return
	}

	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	span.RecordError(err)
}

func (face *RestAPI) GetMap(w http.ResponseWriter, r *http.Request) {
	span := telemetry.SpanFromContext(r.Context())

	resp, err := http.Get("https://github.com/Fribb/anime-lists/raw/master/anime-list-full.json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		span.RecordError(err)
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		span.RecordError(err)
		return
	}

	var entries []entities.Media
	err = json.Unmarshal(data, &entries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		span.RecordError(err)
		return
	}

	telemetry.Int(span, "entries", len(entries))

	records := make(map[string]string, len(entries))
	for _, entry := range entries {
		if entry.SourceID == "" || entry.TargetID == "" {
			continue
		}

		records[entry.SourceID] = entry.TargetID
	}

	newData, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		span.RecordError(err)
		return
	}

	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(newData))
}

func (face *RestAPI) GetUser(w http.ResponseWriter, r *http.Request) {
	span := telemetry.SpanFromContext(r.Context())

	name := r.URL.Query().Get("name")
	if !usernameRegex.MatchString(name) {
		err := fmt.Errorf("invalid username")
		http.Error(w, err.Error(), http.StatusBadRequest)
		span.RecordError(err)
		return
	}

	userId, err := face.mediaLister.GetUserID(r.Context(), name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		span.RecordError(err)
		return
	}

	w.Header().Add("X-Anilist-User-Name", name)
	w.Header().Add("X-Anilist-User-Id", userId)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, userId)
}
