package main

import (
	"github.com/jinzhu/gorm"
	"github.com/mmcdole/gofeed"
	"log"
	"time"
)

type downloader struct {
	period time.Duration
	db     *DB
}

func (d *downloader) start() {
}

func (d *downloader) stop() {
}

func (d *downloader) parseFeeds() {
	parser := gofeed.NewParser()
	feeds := []Feed{}

	err := d.db.allFeeds(&feeds)
	if err != nil {
		panic(err)
	}

	for i := range feeds {
		d.parseFeed(parser, &feeds[i])
	}
}

func (d *downloader) parseFeed(p *gofeed.Parser, f *Feed) {
	if !f.Active {
		return
	}

	feed, err := p.ParseURL(f.Link)
	if err != nil {
		f.Error = err.Error()
		f.Active = false
		d.db.saveFeed(f)
		return
	}

	for _, item := range feed.Items {
		i := Item{}
		switch err := d.db.findItem(f, item.Link, &i); err {
		case gorm.ErrRecordNotFound:
			err = d.db.createItem(f, translateItem(item))
			if err != nil {
				log.Print("Can't create item, database error:", err)
			}
		case nil:
			continue
		default:
			log.Print("Database error when creating feed item", err)
		}
	}
}

func translateItem(item *gofeed.Item) *Item {
	return &Item{
		Title:       item.Title,
		Description: item.Description,
		Content:     item.Content,
		Link:        item.Link,
		Updated:     item.Updated,
		Published:   item.Published}
}
