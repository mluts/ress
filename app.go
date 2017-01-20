package main

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

// App represents
type App struct {
	db *DB
}

type errorResponse struct {
	Error string
}

type createFeedRequest struct {
	Name, URL string
}

func parseJSONRequest(v interface{}, r *http.Request) error {
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(v)
	return err
}

func jsonError(w http.ResponseWriter, msg string) {
	b, err := json.Marshal(errorResponse{msg})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	io.Copy(w, bytes.NewReader(b))
}

func jsonResponse(w http.ResponseWriter, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.Copy(w, bytes.NewReader(b))
}

func (a *App) listFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := a.db.allFeeds()
	if err != nil {
		jsonError(w, err.Error())
		return
	}

	jsonResponse(w, feeds)
}

func (a *App) createFeed(w http.ResponseWriter, r *http.Request) {
	req := createFeedRequest{}
	err := parseJSONRequest(&req, r)

	if err != nil {
		jsonError(w, err.Error())
		return
	}

	err = a.db.createFeed(&Feed{
		Name: req.Name,
		URL:  req.URL})

	if err != nil {
		jsonError(w, err.Error())
	}
}

func (a *App) handler() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/feeds", a.listFeeds).Methods(http.MethodGet)
	r.HandleFunc("/feeds", a.createFeed).Methods(http.MethodPost)
	return r
}
