package main

import (
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
}

func doRequest(method, target string, body []byte) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", "/feeds", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}

func TestFeed_list(t *tt.T) {
	feed := Feed{
		Name: "The Name",
		URL:  "http://example.com/rss"}

	db.createFeed(&feed)

	rec := doRequest("GET", "/feeds", nil)

	if rec.Code != 200 {
		t.Error("Code should be 200, but have", rec.Code)
	}

	feeds := make([]Feed, 1)
	json.Unmarshal(rec.Body.Bytes(), &feeds)
	feed2 := feeds[0]

	if feed.Name != feed2.Name {
		t.Error("Name is not equal", feed.Name, feed2.Name)
	}

	if feed.URL != feed2.URL {
		t.Error("URL is not equal", feed.URL, feed2.URL)
	}

	clearDatabase()
}
