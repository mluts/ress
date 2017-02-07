package app

import (
	"errors"
	"github.com/mmcdole/gofeed"
	"testing"
)

func TestHandleFeedDownload_updates_feed(t *testing.T) {
	clearDatabase()
	url := "http://example.com/rss"

	id, err := app.db.createFeed(&Feed{Link: url})
	if err != nil {
		panic(err)
	}

	feed := &gofeed.Feed{
		Title: "The Title",
	}
	app.handleFeedDownload(url, feed, nil)

	out := Feed{}
	app.db.getFeed(id, &out)

	if out.Title != feed.Title {
		t.Error("Expected to see an updated title")
	}
}

func TestHandleFeedDownload_saves_error(t *testing.T) {
	var dberr error

	clearDatabase()
	err := errors.New("Something bad happened")

	url := "http://example.com/rss"

	id, dberr := app.db.createFeed(&Feed{Link: url})
	if dberr != nil {
		panic(dberr)
	}

	app.handleFeedDownload(url, nil, err)

	out := Feed{}
	app.db.getFeed(id, &out)

	if out.Error != err.Error() {
		t.Error("Expected to see the error saved to db")
	}
}
