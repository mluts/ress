package main

import (
	"flag"
	"github.com/mluts/ress/app"
	"log"
	"net/http"
	"time"
)

var (
	addr                string
	databaseURL         string
	downloadInterval    time.Duration
	downloadConcurrency uint
	staticContentPath   string
	apiPrefix           string
)

func init() {
	flag.StringVar(&addr, "addr", ":8080", "Service address")
	flag.StringVar(&databaseURL, "db", "./db.sqlite", "Database URL")
	flag.DurationVar(&downloadInterval, "interval", time.Second*30, "Download interval")
	flag.UintVar(&downloadConcurrency, "workers", 10, "Amount of workers")
	flag.StringVar(&staticContentPath, "static", "./static", "Static content path")
	flag.StringVar(&apiPrefix, "prefix", "/api", "API prefix")
}

func main() {
	flag.Parse()

	config := &app.AppConfig{
		DBDialect:           "sqlite3",
		DBURL:               databaseURL,
		DownloadInterval:    downloadInterval,
		DownloadConcurrency: downloadConcurrency}

	a, err := app.NewApp(config)
	if err != nil {
		log.Fatal("Can't initialize the app:", err)
	}

	a.Run()
	mux := http.NewServeMux()
	mux.Handle(apiPrefix, http.StripPrefix(apiPrefix, a.Handler()))
	mux.Handle("/", http.FileServer(http.Dir(staticContentPath)))

	log.Printf(`
	Listening at %s %s
	Static content is served from %s
	Database URL is %s
	Download interval/concurrency %v/%d
`, addr, apiPrefix, staticContentPath,
		databaseURL, downloadInterval, downloadConcurrency)

	log.Fatal(http.ListenAndServe(addr, mux))
}
