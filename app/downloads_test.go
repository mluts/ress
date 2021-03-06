package app

import (
	"errors"
	"github.com/mmcdole/gofeed"
	"testing"
)

func TestHandleFeedDownload_updates_feed(t *testing.T) {
	clearDatabase()
	url := "http://example.com/rss"

	id, err := app.db.createFeed(url)
	if err != nil {
		panic(err)
	}

	image := &gofeed.Image{
		URL:   "http://example.com/logo.png",
		Title: "The Feed Logo",
	}

	feed := &gofeed.Feed{
		Title: "The Title",
		Image: image,
	}
	app.handleFeedDownload(url, feed, nil)

	out := Feed{}
	app.db.getFeed(id, &out)

	if out.Title != feed.Title {
		t.Error("Expected to see an updated title")
	}

	if out.Image == nil {
		t.Error("Expected to see a downloaded image")
	}

	if out.Image.URL != image.URL {
		t.Errorf("Expected to have image url %s, but have %s", image.URL, out.Image.URL)
	}

	if out.Image.Title != image.Title {
		t.Errorf("Expected to have image url %s, but have %s", image.Title, out.Image.Title)
	}
}

func TestHandleFeedDownload_saves_error(t *testing.T) {
	var dberr error

	clearDatabase()
	err := errors.New("Something bad happened")

	url := "http://example.com/rss"

	id, dberr := app.db.createFeed(url)
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

func TestHandleFeedDownload_saves_items(t *testing.T) {
	var count int

	clearDatabase()

	url := "http://example.com/rss"

	id, err := app.db.createFeed(url)
	if err != nil {
		panic(err)
	}

	if app.db.getItemsCount(id, &count); count != 0 {
		t.Errorf("Expected to have 0 items in DB, but have %d", count)
	}

	image := &gofeed.Image{
		Title: "The title",
		URL:   "http://example.com/logo.png"}

	feed := &gofeed.Feed{
		Title: "The Title",
		Items: []*gofeed.Item{
			{Title: "The title 1",
				Link:  "http://example.com/1",
				GUID:  "1",
				Image: image},
			{Title: "The title 1",
				Link:  "http://example.com/2",
				GUID:  "2",
				Image: image},
			{Title: "The title 1",
				Link:  "http://example.com/3",
				GUID:  "3",
				Image: image},
		},
	}

	app.handleFeedDownload(url, feed, nil)

	if app.db.getItemsCount(id, &count); count != 3 {
		t.Errorf("Expected to have 3 items in DB, but have %d", count)
	}

	app.handleFeedDownload(url, feed, nil)

	if app.db.getItemsCount(id, &count); count != 3 {
		t.Errorf("Expected to have 3 items in DB, but have %d", count)
	}

	items := make([]Item, 0)
	err = app.db.getItems(id, SQLNoLimit, &items)
	if err != nil {
		panic(err)
	}

	for _, item := range items {
		if item.Image.Title != image.Title {
			t.Errorf("item.Image.Title should be %s, bur have %s", image.Title, item.Image.Title)
		}

		if item.Image.URL != image.URL {
			t.Errorf("item.Image.URL should be %s, bur have %s", image.URL, item.Image.URL)
		}
	}
}
