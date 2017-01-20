package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

// App represents the application API
type App struct {
	db *DB
}

type errorResponse struct {
	Error string
}

type createFeedRequest struct {
	Title, Link string
}

func (a *App) listFeeds(w http.ResponseWriter, r *http.Request) {
	feeds := make([]Feed, 0)
	err := a.db.allFeeds(&feeds)

	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, feeds)
}

func (a *App) createFeed(w http.ResponseWriter, r *http.Request) {
	req := createFeedRequest{}
	err := decodeJSONRequest(&req, r)

	if err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.db.createFeed(&Feed{
		Title: req.Title,
		Link:  req.Link})

	if err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
	}
}

func (a *App) showFeed(w http.ResponseWriter, r *http.Request) {
	var feed Feed

	vars := mux.Vars(r)
	id := vars["id"]

	err := a.db.feed(&feed, id)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}

	jsonResponse(w, feed)
}

func (a *App) feedItems(w http.ResponseWriter, r *http.Request) {
	var (
		feed Feed
		err  error
	)

	vars := mux.Vars(r)
	id := vars["id"]

	err = a.db.feed(&feed, id)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}

	err = a.db.feedItems(&feed)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}

	jsonResponse(w, feed.Items)
}

func (a *App) handler() http.Handler {
	var routes = []struct {
		method  string
		path    string
		handler func(w http.ResponseWriter, r *http.Request)
	}{
		{http.MethodGet, "/feeds", a.listFeeds},
		{http.MethodPost, "/feeds", a.createFeed},
		{http.MethodGet, "/feeds/{id:[0-9]+}", a.showFeed},
		{http.MethodGet, "/feeds/{id:[0-9]+}/items", a.feedItems}}

	r := mux.NewRouter()

	for _, route := range routes {
		r.NewRoute().
			Path(route.path).
			Methods(route.method).
			HandlerFunc(route.handler)
	}

	return r
}
