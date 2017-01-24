package main

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/mmcdole/gofeed"
)

// Item stores information aboud rss feed item
type Item struct {
	gorm.Model

	Feed   Feed `json:"-"`
	FeedID uint

	Title       string
	Description string
	Content     string
	Link        string
	Updated     string
	Published   string
}

func (db *DB) createItem(f *Feed, i *Item) error {
	if db.db.NewRecord(f) {
		return errors.New("Feed is not persisted")
	}

	i.FeedID = f.ID
	return db.db.Create(i).Error
}

func (db *DB) feedItems(f *Feed) error {
	return db.db.Model(f).Related(&f.Items).Error
}

func (db *DB) findItem(f *Feed, link string, i *Item) error {
	return db.db.Model(f).Related(&f.Items).First(i, "link = ?", i.Link).Error
}

func (db *DB) findOrCreateItem(f *Feed, i *Item) (err error) {
	err = db.db.Model(f).Related(&f.Items).First(&Item{}, "link = ?", i.Link).Error

	if err == gorm.ErrRecordNotFound {
		err = db.createItem(f, i)
	}

	return
}

func (db *DB) translateItem(from *gofeed.Item, to *Item) {
	to.Title = from.Title
	to.Description = from.Description
	to.Content = from.Content
	to.Link = from.Link
	to.Updated = from.Updated
	to.Published = from.Published
}
