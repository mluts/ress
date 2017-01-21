package main

import (
	"github.com/mluts/ress/downloader"
	"log"
	"net/http"
	"time"
)

const addr = ":8080"
const dest = "./db.sqlite"

func main() {
	db, err := OpenDatabase("sqlite3", dest)
	if err != nil {
		log.Fatal(err)
	}

	app := &App{db: db}
	app.downloader = downloader.New(
		time.Minute*5,
		5,
		app.handleFeedDownload)

	app.enqueueDownloads()
	ticker := time.NewTicker(time.Minute * 5)
	go func() {
		for {
			<-ticker.C
			app.enqueueDownloads()
		}
	}()

	go app.downloader.Serve()

	log.Print("Listening at", addr)
	log.Fatal(http.ListenAndServe(addr, app.handler()))
}
