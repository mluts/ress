package main

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"net/http"
)

const addr = ":8080"
const dest = "./db.sqlite"

func main() {
	db, err := OpenDatabase("sqlite3", dest)
	if err != nil {
		log.Fatal(err)
	}

	app := &App{db}

	log.Print("Listening at", addr)
	log.Fatal(http.ListenAndServe(addr, app.handler()))
}
