package main

import (
	"github.com/jinzhu/gorm"
	"testing"
)

var exampleItem = Item{
	Title: "The item title",
	Link:  "https://example.com/blog/1"}

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
		t.Fatal("Can't create feed ", err)
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

func TestItem_find_or_create_item_creates_item(t *testing.T) {
	var err error
	clearDatabase()

	f := exampleFeed
	err = db.createFeed(&f)
	if err != nil {
		t.Fatal("Can't create feed ", err)
	}

	i := exampleItem

	if !db.db.NewRecord(&i) {
		t.Fatalf("%v should be new record", i)
	}

	err = db.findOrCreateItem(&f, &i)

	if err != nil {
		t.Error("Failed to create an item ", err)
	}

	if db.db.NewRecord(&i) {
		t.Errorf("%v should be persisted", i)
	}

	newItem := i
	newItem.Title = "Title 2"

	err = db.findOrCreateItem(&f, &newItem)
	db.db.First(&newItem, "link = ?", i.Link)

	if err != nil {
		t.Fatal("Something went wrong", err)
	}

	if newItem.Title != i.Title {
		t.Error("Item title should not be updated")
	}
}

func TestItem_create_passes_validations(t *testing.T) {
	var err error
	clearDatabase()

	f := exampleFeed
	i := exampleItem

	i.Title = ""
	err = db.createItem(&f, &i)

	if err == nil {
		t.Error("Item should not be created without a title")
	}

	i = exampleItem
	i.Link = ""
	err = db.createItem(&f, &i)

	if err == nil {
		t.Error("Item should not be created without a link")
	}
}
