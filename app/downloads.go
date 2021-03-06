package app

import (
	"database/sql"
	"github.com/mmcdole/gofeed"
	"log"
)

func (a *App) enqueueDownloads() {
	feeds := []Feed{}
	err := a.db.getFeeds(SQLNoLimit, &feeds)

	if err != nil {
		log.Print("Failed to enqueue downloads due to error: ", err)
		return
	}

	for _, feed := range feeds {
		a.downloader.Enqueue(feed.Link)
	}
}

func (a *App) handleFeedDownload(url string, feed *gofeed.Feed, err error) {
	f := Feed{}

	if dberr := a.db.getFeedByLink(url, &f); dberr != nil {
		log.Print("Can't get the feed from database:", dberr)
		a.downloader.Discard(url)
		return
	}

	if err != nil {
		f.Error = err.Error()
		a.db.updateFeed(f.ID, &f)
		return
	}

	f.Title = feed.Title

	if feed.PublishedParsed != nil {
		f.Published = feed.PublishedParsed
	}

	if feed.UpdatedParsed != nil {
		f.Updated = feed.UpdatedParsed
	}

	if feed.Image != nil {
		f.Image = &Image{
			URL:   feed.Image.URL,
			Title: feed.Image.Title,
		}
	}

	if e := a.db.updateFeed(f.ID, &f); e != nil {
		log.Printf("Can't save feed: %v", e)
		return
	}

	for _, item := range feed.Items {
		guid := itemGUID(item)
		switch dberr := a.db.getItemByGUID(f.ID, guid, &Item{}); dberr {
		case sql.ErrNoRows:
			newItem := Item{}
			translateItem(item, &newItem)
			a.db.createItem(f.ID, &newItem)
		case nil:
		default:
			log.Printf("Failed to search a feed item: %v", dberr)
		}
	}
}

func translateItem(from *gofeed.Item, to *Item) {
	to.Title = from.Title
	to.Link = from.Link
	to.Description = from.Description
	to.Content = from.Content
	to.GUID = itemGUID(from)

	if from.PublishedParsed != nil {
		to.Published = from.PublishedParsed
	}

	if from.UpdatedParsed != nil {
		to.Updated = from.UpdatedParsed
	}

	if from.Image != nil {
		to.Image = &Image{
			Title: from.Image.Title,
			URL:   from.Image.URL}
	}
}

func itemGUID(item *gofeed.Item) string {
	if len(item.GUID) != 0 {
		return item.GUID
	}

	return item.Title
}
