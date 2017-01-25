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
	Link        string `gorm:"not null;index"`
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
	return db.db.Model(f).Related(&f.Items).First(i, "link = ?", link).Error
}

func (db *DB) findOrCreateItem(f *Feed, i *Item) (err error) {
	err = db.findItem(f, i.Link, &Item{})

	if err == gorm.ErrRecordNotFound {
		err = db.createItem(f, i)
	}

	return
}

func translateItem(from *gofeed.Item, to *Item) {
	to.Title = from.Title
	to.Description = from.Description
	to.Content = from.Content
	to.Link = from.Link
	to.Updated = from.Updated
	to.Published = from.Published
}

func (db *DB) validateItem(i *Item) (err error) {
	if len(i.Title) == 0 {
		return errors.New("Title can't be blank")
	}

	if len(i.Link) == 0 {
		return errors.New("Link can't be blank")
	}

	return
}
