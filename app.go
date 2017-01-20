package main

import (
	"github.com/gorilla/mux"
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
	err := decodeJSONRequest(&req, r)

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
