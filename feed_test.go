package main

import "testing"

var exampleFeed = Feed{
	Title: "Example Feed",
	Link:  "http://example.com/rss"}

var validateTests = []struct {
	feed  *Feed
	valid bool
}{
	{feed: &Feed{Title: "", Link: ""},
		valid: false},
	{feed: &Feed{Title: "Some Title", Link: ""},
		valid: false},
	{feed: &Feed{Title: "", Link: "Some Title"},
		valid: false},
	{feed: &Feed{Title: "Some Title", Link: "Some Link"},
		valid: true}}

func TestFeed_validate(t *testing.T) {
	for _, tt := range validateTests {
		err := tt.feed.validate()

		if err == nil && !tt.valid {
			t.Errorf("%v should be invalid", tt.feed)
		} else if err != nil && tt.valid {
			t.Errorf("%v should be valid", tt.feed)
		}
	}
}

func TestFeed_create(t *testing.T) {
	feed := Feed{Title: "The Title", Link: "url"}
	err := db.createFeed(&feed)
	if err != nil {
		panic(err)
	}

	if db.db.NewRecord(feed) {
		t.Error("Feed should be persisted")
	}

	if !feed.Active {
		t.Error("Feed should be Active initially")
	}

	if len(feed.Error) != 0 {
		t.Error("Feed should not contain errors")
	}
}
