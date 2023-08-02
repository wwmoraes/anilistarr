package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/wwmoraes/anilistarr/internal/entities"
	"github.com/wwmoraes/anilistarr/internal/telemetry"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

type RestAPI interface {
	GetList(http.ResponseWriter, *http.Request)
	GetMap(http.ResponseWriter, *http.Request)
	GetUser(http.ResponseWriter, *http.Request)
}

type restAPI struct {
	mapper *usecases.MediaLinker
}

func NewRestAPI(mapper *usecases.MediaLinker) (RestAPI, error) {
	return &restAPI{
		mapper: mapper,
	}, nil
}

func (face *restAPI) GetList(w http.ResponseWriter, r *http.Request) {
	span := telemetry.SpanFromContext(r.Context())

	username := r.URL.Query().Get("username")
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	customList, err := face.mapper.GenerateCustomList(r.Context(), username)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		span.RecordError(err)
		return
	}

	data, err := json.Marshal(customList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err.Error())
		span.RecordError(err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	span.RecordError(err)
}

func (face *restAPI) GetMap(w http.ResponseWriter, r *http.Request) {
	span := telemetry.SpanFromContext(r.Context())

	resp, err := http.Get("https://github.com/Fribb/anime-lists/raw/master/anime-list-full.json")
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		span.RecordError(err)
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err.Error())
		span.RecordError(err)
		return
	}

	var entries []entities.Media
	err = json.Unmarshal(data, &entries)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err.Error())
		span.RecordError(err)
		return
	}

	telemetry.Int(span, "entries", len(entries))

	records := make(map[string]string, len(entries))
	for _, entry := range entries {
		if entry.AnilistID == "" || entry.TvdbID == "" {
			continue
		}

		records[entry.AnilistID] = entry.TvdbID
	}

	newData, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err.Error())
		span.RecordError(err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(newData))
}

func (face *restAPI) GetUser(w http.ResponseWriter, r *http.Request) {
	span := telemetry.SpanFromContext(r.Context())

	name := r.URL.Query().Get("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userId, err := face.mapper.GetUserID(r.Context(), name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err.Error())
		span.RecordError(err)
		return
	}

	w.Header().Add("X-Anilist-User-Name", name)
	w.Header().Add("X-Anilist-User-Id", userId)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, userId)
}
