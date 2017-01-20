package main

import "testing"

var validateTests = []struct {
	feed  *Feed
	valid bool
}{
	{feed: &Feed{Name: "", URL: ""},
		valid: false},
	{feed: &Feed{Name: "Some Name", URL: "Some URL"},
		valid: false}}

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
