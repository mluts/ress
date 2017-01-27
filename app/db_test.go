package app

import (
	"database/sql"
	"testing"
)

func TestDB_migrations_work(t *testing.T) {
	_, err := OpenDatabase("sqlite3", ":memory:")
	if err != nil {
		t.Error("DB was not initialized:", err)
	}
}

func opendb() *DB {
	db, err := OpenDatabase("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	return db
}

func nullstring(str string) sql.NullString {
	return sql.NullString{str, true}
}

var feedExamples = []struct {
	feed Feed
	ok   bool
}{
	{Feed{
		Title: "The title1",
		Link:  "The link1"}, true},

	{Feed{
		Title: "The title1"}, false},

	{Feed{
		Link: "The title1"}, false},

	{Feed{
		Title: "The title1",
		Link:  "The link1"}, false},

	{Feed{
		Title: "The title1",
		Link:  "The link2"}, true},

	{Feed{
		Title: "The title2",
		Link:  "The link3"}, true},
}

var itemExamples = []Item{
	{
		Title: "The Title 1",
		Link:  "The Link 1",
	},
}

func TestDB_createFeed(t *testing.T) {
	db := opendb()

	for _, example := range feedExamples {
		id, err := db.createFeed(&example.feed)
		t.Logf("Created feed with id %d", id)

		if example.ok && err != nil {
			t.Errorf("Expected %v to be persisted: %v", example.feed, err)
		} else if !example.ok && err == nil {
			t.Errorf("Expected %v not to be persisted: %v", example.feed, err)
		}
	}
}

func TestDB_deleteFeed(t *testing.T) {
	db := opendb()

	id, err := db.createFeed(&feedExamples[0].feed)
	if err != nil {
		t.Fatal(err)
	}

	err = db.deleteFeed(id)
	if err != nil {
		t.Errorf("Expected feed %d to be deleted: %v", id, err)
	}
}

func TestDB_getFeed(t *testing.T) {
	var example = feedExamples[0].feed
	db := opendb()

	id, err := db.createFeed(&example)
	if err != nil {
		t.Fatal(err)
	}
	out := Feed{}

	err = db.getFeed(id, &out)
	if err != nil {
		t.Fatal(err)
	}

	if out.Link != example.Link {
		t.Error("Expected to have link %s, but have %s", example.Link, out.Link)
	}
}
