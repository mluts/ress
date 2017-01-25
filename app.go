package main

import (
	"github.com/gorilla/mux"
	"github.com/mluts/ress/downloader"
	"github.com/mmcdole/gofeed"
	"log"
	"net/http"
	"time"
)

// App represents the application API
type App struct {
	db         *DB
	downloader *downloader.Downloader
}

// AppConfig is an application configuration
type AppConfig struct {
	dbDialect           string
	dbURL               string
	downloadInterval    time.Duration
	downloadConcurrency uint
}

type errorResponse struct {
	Error string
}

type createFeedRequest struct {
	Title, Link string
}

// NewApp initializes new application
func NewApp(config *AppConfig) (*App, error) {
	db, err := OpenDatabase(config.dbDialect, config.dbURL)
	if err != nil {
		return nil, err
	}

	app := &App{db: db}

	d := downloader.New(
		config.downloadInterval,
		config.downloadConcurrency,
		app.handleFeedDownload)

	app.downloader = d

	return app, nil
}

// Run enqueues downloads and starts necessary goroutines
func (a *App) Run() {
	a.enqueueDownloads()
	go a.downloader.Serve()
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

func (a *App) enqueueDownloads() {
	feeds := []Feed{}
	err := a.db.allFeeds(&feeds)

	if err != nil {
		log.Print("Failed to enqueue downloads due to error: ", err)
		return
	}

	for _, feed := range feeds {
		a.downloader.Enqueue(feed.Link)
	}
}

func (a *App) handleFeedDownload(url string, feed *gofeed.Feed, err error) {
	var dberr error
	f := Feed{}
	dberr = a.db.db.Where("link = ?", url).First(&f).Error

	if dberr != nil {
		log.Print("Database error during feed fetching ", dberr)
		a.downloader.Discard(url)
		return
	}

	if err != nil {
		a.downloader.Discard(url)
		f.Error = err.Error()
		f.Active = false
		a.db.db.Save(&f)
		return
	}

	for _, item := range feed.Items {
		newItem := &Item{}
		translateItem(item, newItem)

		err := a.db.findOrCreateItem(&f, newItem)

		if err != nil {
			log.Print("Failed to create an item ", err)
		}
	}
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
