package main

import (
	"github.com/mluts/ress/app"
	"log"
	"net/http"
	"time"
)

const addr = ":8080"

func main() {
	config := &app.AppConfig{
		DBDialect:           "sqlite3",
		DBURL:               "./db.sqlite",
		DownloadInterval:    time.Second * 30,
		DownloadConcurrency: 10}

	app, err := NewApp(config)
	if err != nil {
		log.Fatal("Can't initialize the app:", err)
	}

	app.Run()
	log.Print("Listening at", addr)
	log.Fatal(http.ListenAndServe(addr, app.handler()))
}
