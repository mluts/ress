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

func feedsCount() (count int64) {
	err := db.Get(&count, "SELECT COUNT(id) FROM feeds")
	if err != nil {
		panic(err)
	}
	return
}

func itemsCount() (count int64) {
	err := db.Get(&count, "SELECT COUNT(id) FROM items")
	if err != nil {
		panic(err)
	}
	return
}

func clearDatabase() {
	db.MustExec("DELETE FROM feeds")
	db.MustExec("DELETE FROM items")
}

func doRequest(method, target string, body []byte) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, bytes.NewReader(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}

func TestFeed_api_list(t *testing.T) {
	clearDatabase()

	feed := feedExamples[0].feed

	db.createFeed(&feed)

	rec := doRequest("GET", "/feeds", nil)

	if rec.Code != 200 {
		t.Error("Code should be 200, but have", rec.Code)
		t.Error(rec.Body.String())
	}

	feeds := make([]Feed, 1)
	json.Unmarshal(rec.Body.Bytes(), &feeds)
	if len(feeds) != 1 {
		t.Fatal("Expected to have at least one feed")
	}

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
		count int64
	)
	clearDatabase()

	feed := feedExamples[0].feed

	b, err = json.Marshal(feed)
	if err != nil {
		panic(err)
	}

	count = feedsCount()
	if count != 0 {
		t.Error("Expected to have 0 feeds in db, but have", count)
	}

	rec := doRequest("POST", "/feeds", b)

	if rec.Code != 200 {
		t.Error("Expected to have code 200, but have", rec.Code)
	}

	count = feedsCount()
	if count != 1 {
		t.Error("Expected to have 1 created feed but have", count)
	}
}

func TestFeed_api_show(t *testing.T) {
	clearDatabase()
	feed := feedExamples[0].feed
	id, err := app.db.createFeed(&feed)

	if err != nil {
		panic(err)
	}

	path := strings.Join([]string{
		"/feeds",
		strconv.Itoa(int(id))}, "/")

	rec := doRequest("GET", path, nil)

	if rec.Code != 200 {
		t.Error("Expected to see status 200, but seeing", rec.Code)
		t.Error(rec.Body.String())
	}

	feed2 := Feed{}
	json.Unmarshal(rec.Body.Bytes(), &feed2)

	if id != feed2.ID {
		t.Error("Have wrong feed id", feed2.ID)
	}
}

func TestFeed_api_feed_items(t *testing.T) {
	clearDatabase()
	feed := feedExamples[0].feed
	id, err := app.db.createFeed(&feed)

	if err != nil {
		panic(err)
	}

	for _, item := range itemExamples {
		err = db.createItem(id, &item)
		if err != nil {
			panic(err)
		}
	}

	path := strings.Join([]string{
		"/feeds",
		strconv.Itoa(int(id)),
		"items"}, "/")

	rec := doRequest("GET", path, nil)

	if rec.Code != 200 {
		t.Error("Expected to see code 200, but have", rec.Code)
		t.Error(rec.Body.String())
		return
	}

	items := make([]Item, 0)
	json.Unmarshal(rec.Body.Bytes(), &items)

	if len(items) != len(itemExamples) {
		t.Error("Expected to have 2 items, but have", len(items))
		t.Error(rec.Body.String())
		return
	}

	if items[0].FeedID != id {
		t.Errorf("Expected item feedID to eq %d, but have %d", id, items[0].FeedID)
	}

	if strings.Contains(rec.Body.String(), "\"Feed\":") {
		t.Errorf("Feed should not be included in the json")
	}

	if !items[0].Unread {
		t.Error("Item should be unread initially")
	}
}

func TestFeed_api_delete_feed(t *testing.T) {
}
