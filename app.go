package main

import (
	"github.com/mluts/ress/downloader"
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
