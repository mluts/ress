package main

import (
	"errors"
	"github.com/jinzhu/gorm"
)

// Feed stores useful information about subscribed feed
type Feed struct {
	gorm.Model

	Title string
	Link  string

	// Indicates feed error
	Error  string
	Active bool

	Items []Item
}

func (f *Feed) validate() error {
	if len(f.Title) == 0 {
		return errors.New("Title can't be blank")
	}

	if len(f.Link) == 0 {
		return errors.New("Link can't be blank")
	}

	return nil
}

func (db *DB) allFeeds() ([]Feed, error) {
	feeds := []Feed{}
	err := db.db.Find(&feeds).Error
	return feeds, err
}

func (db *DB) createFeed(f *Feed) error {
	err := f.validate()
	if err != nil {
		return err
	}

	f.Active = true

	return db.db.Create(f).Error
}

func (db *DB) feed(f *Feed, id interface{}) error {
	return db.db.First(f, id).Error
}

func (db *DB) feedsCount(c *int) error {
	return db.db.Model(&Feed{}).Count(c).Error
}
