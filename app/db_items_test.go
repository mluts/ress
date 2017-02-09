package app

import "testing"

func TestItem_get_items(t *testing.T) {
	clearDatabase()

	id1, err := db.createFeed("http://example.com/1")
	if err != nil {
		panic(err)
	}

	id2, err := db.createFeed("http://example.com/2")
	if err != nil {
		panic(err)
	}

	item1, item2 := itemExamples[0], itemExamples[1]

	_, err = db.createItem(id1, &item1)
	if err != nil {
		panic(err)
	}

	_, err = db.createItem(id2, &item2)
	if err != nil {
		panic(err)
	}

	items := []Item{}
	db.getItems(id2, SQLNoLimit, &items)

	if len(items) != 1 {
		t.Errorf("Expected to see 1 item, but have %d", len(items))
	}
}
