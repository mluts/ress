package main

import (
	"errors"
	"github.com/jinzhu/gorm"
)

// Feed stores useful information about subscribed feed
type Feed struct {
	gorm.Model

	Name string
	URL  string
}

func (f *Feed) validate() error {
	if len(f.Name) == 0 {
		return errors.New("Name can't be blank")
	}

	if len(f.URL) == 0 {
		return errors.New("URL can't be blank")
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

	return db.db.Create(f).Error
}
