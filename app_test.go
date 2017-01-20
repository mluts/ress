package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	tt "testing"
)

var (
	db      *DB
	app     *App
	handler http.Handler
)

func init() {
	var err error

	db, err = OpenDatabase("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	app = &App{db}
	handler = app.handler()
}

func clearDatabase() {
	db.db.Delete(&Feed{})
	db.db.Delete(&Item{})
}

func doRequest(method, target string, body []byte) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, bytes.NewReader(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}

func TestFeed_api_list(t *tt.T) {
	clearDatabase()

	feed := Feed{
		Title: "The Title",
		Link:  "http://example.com/rss"}

	db.createFeed(&feed)

	rec := doRequest("GET", "/feeds", nil)

	if rec.Code != 200 {
		t.Error("Code should be 200, but have", rec.Code)
	}

	feeds := make([]Feed, 1)
	json.Unmarshal(rec.Body.Bytes(), &feeds)
	feed2 := feeds[0]

	if feed.Title != feed2.Title {
		t.Error("Title is not equal", feed.Title, feed2.Title)
	}

	if feed.Link != feed2.Link {
		t.Error("Link is not equal", feed.Link, feed2.Link)
	}
}

func TestFeed_api_create(t *tt.T) {
	var (
		err   error
		b     []byte
		count int
	)
	clearDatabase()

	feed := Feed{
		Title: "The Title",
		Link:  "http://example.com/rss"}

	b, err = json.Marshal(feed)
	if err != nil {
		panic(err)
	}

	app.db.feedsCount(&count)

	if count != 0 {
		t.Error("Expected to have 0 feeds in db, but have", count)
	}

	rec := doRequest("POST", "/feeds", b)

	if rec.Code != 200 {
		t.Error("Expected to have code 200, but have", rec.Code)
	}

	app.db.feedsCount(&count)

	if count != 1 {
		t.Error("Expected to have 1 created feed but have", count)
	}
}
