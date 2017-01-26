package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

type createFeedRequest struct {
	Title, Link string
}

type route struct {
	method  string
	path    string
	handler func(http.ResponseWriter, *http.Request)
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
		return
	}

	a.enqueueDownloads()
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

func (a *App) deleteFeed(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	id := vars["id"]

	err = a.db.db.Delete(&Feed{}, "id = ?", id).Error
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.db.db.Delete(&Item{}, "feed_id = ?", id).Error
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
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
	var routes = []route{
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
