package main

import (
	"log"
	"net/http"
	"time"
)

const addr = ":8080"

func main() {
	config := &AppConfig{
		dbDialect:           "sqlite3",
		dbURL:               "./db.sqlite",
		downloadInterval:    time.Second * 30,
		downloadConcurrency: 10}

	app, err := NewApp(config)
	if err != nil {
		log.Fatal("Can't initialize the app:", err)
	}

	app.Run()
	log.Print("Listening at", addr)
	log.Fatal(http.ListenAndServe(addr, app.handler()))
}
