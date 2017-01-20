package main

import (
	"github.com/jinzhu/gorm"
	"testing"
)

func TestItem_create(t *testing.T) {
	var err error

	clearDatabase()

	f := exampleFeed
	i := Item{}

	err = db.createItem(&f, &i)
	if err == nil {
		t.Error("Can't create item without a feed")
		return
	}

	db.createFeed(&f)
	err = db.createItem(&f, &i)

	if err != nil {
		t.Error("Should have created an item")
		return
	}

	var count int

	db.db.Model(&Item{}).Count(&count)

	if count != 1 {
		t.Error("Expected count to be 1, but have ", count)
	}
}

func TestItem_find(t *testing.T) {
	var err error

	clearDatabase()

	f := exampleFeed
	err = db.createFeed(&f)
	if err != nil {
		panic(err)
	}

	link := f.Link
	i := Item{}

	err = db.findItem(&f, link, &i)
	if err != gorm.ErrRecordNotFound {
		t.Error("Expected 'not found' eror, but have", err)
	}

	err = db.createItem(&f, &Item{})

	if err != nil {
		t.Error("Expected not to have error, but have", err)
	}
}
