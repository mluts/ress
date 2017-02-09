package app

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
		DBDialect:           "sqlite3",
		DBURL:               ":memory:",
		DownloadInterval:    time.Second,
		DownloadConcurrency: 1}

	app, err = NewApp(config)
	if err != nil {
		panic(err)
	}

	db = app.db

	handler = app.Handler()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func feedsCount() (count int64) {
	must(db.db.Get(&count, "SELECT COUNT(id) FROM feeds"))
	return
}

func itemsCount() (count int64) {
	must(db.db.Get(&count, "SELECT COUNT(id) FROM items"))
	return
}

func clearDatabase() {
	db.db.MustExec("DELETE FROM feeds")
	db.db.MustExec("DELETE FROM items")
}

func doRequest(method, target string, body []byte) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, bytes.NewReader(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}

func assertRequestCode(t *testing.T, r *httptest.ResponseRecorder, code int) {
	if r.Code != code {
		t.Errorf("Expected to see code %d, but have %d", code, r.Code)
		t.Error(r.Body.String())
	}
}

func assertIntEq(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Error("Expected %d to equal to %d", actual, expected)
	}
}

func TestAPI_feeds(t *testing.T) {
	clearDatabase()

	feed := feedExamples[0].feed

	db.createFeed(feed.Link)

	rec := doRequest("GET", "/feeds", nil)

	assertRequestCode(t, rec, 200)

	feeds := make([]Feed, 1)
	json.Unmarshal(rec.Body.Bytes(), &feeds)
	if len(feeds) != 1 {
		t.Fatal("Expected to have at least one feed")
	}

	feed2 := feeds[0]

	if feed.Link != feed2.Link {
		t.Error("Link is not equal", feed.Link, feed2.Link)
	}
}

func TestAPI_create_feed(t *testing.T) {
	var (
		err   error
		b     []byte
		count int64
	)
	clearDatabase()

	b, err = json.Marshal(&Feed{Link: "https://example.com/feed"})
	if err != nil {
		panic(err)
	}

	count = feedsCount()
	if count != 0 {
		t.Error("Expected to have 0 feeds in db, but have", count)
	}

	rec := doRequest("POST", "/feeds", b)

	assertRequestCode(t, rec, 200)

	count = feedsCount()
	if count != 1 {
		t.Error("Expected to have 1 created feed but have", count)
	}
}

func TestAPI_show_feed(t *testing.T) {
	clearDatabase()
	feed := feedExamples[0].feed
	id, err := app.db.createFeed(feed.Link)

	if err != nil {
		panic(err)
	}

	path := strings.Join([]string{
		"/feeds",
		strconv.Itoa(int(id))}, "/")

	rec := doRequest("GET", path, nil)

	assertRequestCode(t, rec, 200)

	feed2 := Feed{}
	json.Unmarshal(rec.Body.Bytes(), &feed2)

	if id != feed2.ID {
		t.Error("Have wrong feed id", feed2.ID)
	}
}

func TestAPI_feed_items(t *testing.T) {
	clearDatabase()
	feed := feedExamples[0].feed
	id, err := app.db.createFeed(feed.Link)

	if err != nil {
		panic(err)
	}

	for _, item := range itemExamples {
		_, err = db.createItem(id, &item)
		if err != nil {
			panic(err)
		}
	}

	path := strings.Join([]string{
		"/feeds",
		strconv.Itoa(int(id)),
		"items"}, "/")

	rec := doRequest("GET", path, nil)

	assertRequestCode(t, rec, 200)

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

func TestAPI_delete_feed(t *testing.T) {
	clearDatabase()

	feed := feedExamples[0].feed
	id, err := db.createFeed(feed.Link)
	if err != nil {
		t.Fatal(err)
	}

	if len(itemExamples) == 0 {
		t.Fatal("Don't have item examples")
	}

	for _, item := range itemExamples {
		_, err := db.createItem(id, &item)
		if err != nil {
			t.Errorf("Expected to save feed, but had an error: %v", err)
		}
	}

	assertIntEq(t, 1, int(feedsCount()))
	assertIntEq(t, len(itemExamples), int(itemsCount()))

	path := strings.Join([]string{
		"/feeds",
		strconv.Itoa(int(id))}, "/")

	rec := doRequest("DELETE", path, nil)

	assertRequestCode(t, rec, 200)
	assertIntEq(t, 0, int(feedsCount()))
	assertIntEq(t, 0, int(itemsCount()))
}

func TestAPI_mark_item_read(t *testing.T) {
	var (
		err        error
		id, itemID int64
	)
	clearDatabase()

	feed := feedExamples[0].feed
	id, err = db.createFeed(feed.Link)
	if err != nil {
		t.Fatal(err)
	}

	item := itemExamples[0]

	itemID, err = db.createItem(id, &item)
	if err != nil {
		t.Fatal(err)
	}

	if err := db.getItem(id, &item); err != nil {
		t.Fatal(err)
	}

	if !item.Unread {
		t.Error("Item should be unread")
	}

	path := strings.Join([]string{
		"/feeds",
		strconv.Itoa(int(id)),
		"items",
		strconv.Itoa(int(itemID)),
		"read"}, "/")

	rec := doRequest("POST", path, nil)

	assertRequestCode(t, rec, 200)

	i := Item{}
	err = db.getItem(itemID, &i)
	if err != nil {
		t.Fatal(err)
	}

	if i.Unread {
		t.Error("Item should be marked as read")
	}

	rec = doRequest("DELETE", path, nil)
	assertRequestCode(t, rec, 200)

	err = db.getItem(itemID, &i)
	if err != nil {
		t.Fatal(err)
	}

	if !i.Unread {
		t.Error("Item should be marked as unread")
	}
}
