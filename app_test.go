package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	db      *DB
	app     *App
	handler http.Handler
)

func init() {
	var err error
	config := &AppConfig{
		dbDialect:           "sqlite3",
		dbURL:               ":memory:",
		downloadInterval:    time.Second,
		downloadConcurrency: 1}

	app, err = NewApp(config)
	if err != nil {
		panic(err)
	}

	db = app.db

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

func TestFeed_api_list(t *testing.T) {
	clearDatabase()

	feed := exampleFeed

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

func TestFeed_api_create(t *testing.T) {
	var (
		err   error
		b     []byte
		count int
	)
	clearDatabase()

	feed := exampleFeed

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

func TestFeed_api_show(t *testing.T) {
	clearDatabase()
	feed := exampleFeed
	err := app.db.createFeed(&feed)

	if err != nil {
		panic(err)
	}

	path := strings.Join([]string{
		"/feeds",
		strconv.Itoa(int(feed.ID))}, "/")

	rec := doRequest("GET", path, nil)

	if rec.Code != 200 {
		t.Error("Expected to see status 200, but seeing", rec.Code)
	}

	feed2 := Feed{}
	json.Unmarshal(rec.Body.Bytes(), &feed2)

	if feed.ID != feed2.ID {
		t.Error("Have wrong feed id", feed2.ID)
	}
}

func TestFeed_api_feed_items(t *testing.T) {
	clearDatabase()
	feed := exampleFeed
	err := app.db.createFeed(&feed)

	if err != nil {
		panic(err)
	}

	for _, item := range []Item{Item{Title: "Item1"}, Item{Title: "Item2"}} {
		err = db.createItem(&feed, &item)
		if err != nil {
			panic(err)
		}
	}

	path := strings.Join([]string{
		"/feeds",
		strconv.Itoa(int(feed.ID)),
		"items"}, "/")

	rec := doRequest("GET", path, nil)

	if rec.Code != 200 {
		t.Error("Expected to see code 200, but have", rec.Code)
		return
	}

	items := make([]Item, 0)
	json.Unmarshal(rec.Body.Bytes(), &items)

	if len(items) != 2 {
		t.Error("Expected to have 2 items, but have", len(items))
		return
	}

	if items[0].FeedID != feed.ID {
		t.Errorf("Expected item feedID to eq %d, but have %d", feed.ID, items[0].FeedID)
	}

	if strings.Contains(rec.Body.String(), "\"Feed\":") {
		t.Errorf("Feed should not be included in the json")
	}
}
