package main

import "testing"

var validateTests = []struct {
	feed  *Feed
	valid bool
}{
	{feed: &Feed{Name: "", URL: ""},
		valid: false},
	{feed: &Feed{Name: "Some Name", URL: ""},
		valid: false},
	{feed: &Feed{Name: "", URL: "Some Name"},
		valid: false},
	{feed: &Feed{Name: "Some Name", URL: "Some URL"},
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
	feed := Feed{Name: "The Name", URL: "url"}
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
