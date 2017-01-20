package main

import (
	"errors"
	"github.com/jinzhu/gorm"
)

// Item stores information aboud rss feed item
type Item struct {
	gorm.Model

	Feed   Feed
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

func (db *DB) itemsCount(c *int) error {
	return db.db.Model(&Item{}).Count(c).Error
}

func (db *DB) feedItems(f *Feed) error {
	return db.db.Model(f).Related(&f.Items).Error
}

func (db *DB) findItem(f *Feed, link string, i *Item) error {
	return db.db.Model(f).Related(&f.Items).First(i, "link = ?", i.Link).Error
}
