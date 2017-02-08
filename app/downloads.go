package app

import (
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
		a.db.updateFeed(f.ID, &Feed{
			Error:  err.Error(),
			Active: false,
		})
		return
	}

	if e := a.db.updateFeed(f.ID, &Feed{
		Title: feed.Title,
	}); e != nil {
		log.Printf("Can't save feed: %v", e)
	}

	for _, item := range feed.Items {
		newItem := Item{}
		translateItem(item, &newItem)
		_, err := a.db.createItem(f.ID, &newItem)
		if err != nil {
			log.Printf("Failed to create an item: %v", err)
		}
	}
}

func translateItem(from *gofeed.Item, to *Item) {
	to.Title = from.Title
	to.Link = from.Link
	to.Description = from.Description
	to.Content = from.Content
}
