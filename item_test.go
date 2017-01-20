package main

import "testing"

func TestItem_create(t *testing.T) {
	var err error

	clearDatabase()

	f := Feed{
		Title: "The Name", Link: "The Link"}
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
