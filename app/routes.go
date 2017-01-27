package app

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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
	feeds := []Feed{}
	err := a.db.getFeeds(-1, &feeds)

	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, &feeds)
}

func (a *App) createFeed(w http.ResponseWriter, r *http.Request) {
	req := createFeedRequest{}
	err := decodeJSONRequest(&req, r)

	if err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = a.db.createFeed(&Feed{
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
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}

	err = a.db.getFeed(int64(id), &feed)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}

	jsonResponse(w, feed)
}

func (a *App) deleteFeed(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.db.deleteFeed(int64(id))
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (a *App) feedItems(w http.ResponseWriter, r *http.Request) {
	var (
		items []Item
		err   error
		id    int
	)

	vars := mux.Vars(r)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}

	err = a.db.getItems(int64(id), -1, &items)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}

	jsonResponse(w, &items)
}

func (a *App) markItemRead(w http.ResponseWriter, r *http.Request) {
	var (
		itemID int
		err    error
	)
	vars := mux.Vars(r)

	if itemID, err = strconv.Atoi(vars["id"]); err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = a.db.markItemRead(int64(itemID), true); err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (a *App) unmarkItemRead(w http.ResponseWriter, r *http.Request) {
	var (
		itemID int
		err    error
	)

	vars := mux.Vars(r)

	if itemID, err = strconv.Atoi(vars["id"]); err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = a.db.markItemRead(int64(itemID), false); err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (a *App) Handler() http.Handler {
	var routes = []route{
		{http.MethodGet, "/feeds", a.listFeeds},
		{http.MethodPost, "/feeds", a.createFeed},
		{http.MethodGet, "/feeds/{id:[0-9]+}", a.showFeed},
		{http.MethodDelete, "/feeds/{id:[0-9]+}", a.deleteFeed},
		{http.MethodGet, "/feeds/{id:[0-9]+}/items", a.feedItems},
		{http.MethodPost, "/feeds/{feed_id:[0-9]+}/items/{id:[0-9]+}/read", a.markItemRead},
		{http.MethodDelete, "/feeds/{feed_id:[0-9]+}/items/{id:[0-9]+}/read", a.unmarkItemRead},
	}

	r := mux.NewRouter()

	for _, route := range routes {
		r.NewRoute().
			Path(route.path).
			Methods(route.method).
			HandlerFunc(route.handler)
	}

	return r
}
