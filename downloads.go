package main

import (
	"github.com/mmcdole/gofeed"
	"log"
)

func (a *App) enqueueDownloads() {
	feeds := []Feed{}
	err := a.db.allFeeds(&feeds)

	if err != nil {
		log.Print("Failed to enqueue downloads due to error: ", err)
		return
	}

	for _, feed := range feeds {
		a.downloader.Enqueue(feed.Link)
	}
}

func (a *App) handleFeedDownload(url string, feed *gofeed.Feed, err error) {
	var dberr error
	f := Feed{}
	dberr = a.db.db.Where("link = ?", url).First(&f).Error

	if dberr != nil {
		log.Print("Database error during feed fetching ", dberr)
		a.downloader.Discard(url)
		return
	}

	if err != nil {
		a.downloader.Discard(url)
		f.Error = err.Error()
		f.Active = false
		a.db.db.Save(&f)
		return
	}

	for _, item := range feed.Items {
		newItem := &Item{}
		translateItem(item, newItem)

		err := a.db.findOrCreateItem(&f, newItem)

		if err != nil {
			log.Print("Failed to create an item ", err)
		}
	}
}
