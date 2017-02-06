package app

func (db *DB) createItem(feedID int64, item *Item) (int64, error) {
	item.FeedID = feedID

	stmt := db.prepareNamed(
		"createItem",
		`INSERT INTO items (feed_id, title, link, description, content)
		 VALUES (:feed_id, :title, :link, :description, :content)`)

	result, err := stmt.Exec(item)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (db *DB) updateItem(item *Item) error {
	stmt := db.prepareNamed(
		"updateItem",
		`UPDATE items SET (title, link, description, content) =
		 (:link, :description, :content)`)

	_, err := stmt.Exec(item)
	return err
}

func (db *DB) deleteItem(id int64) error {
	stmt := db.prepare("deleteItem", "DELETE FROM items WHERE id = $1")
	_, err := stmt.Exec(id)
	return err
}

func (db *DB) getItem(id int64, out *Item) error {
	stmt := db.prepare("getItem", "SELECT * FROM items_view WHERE id = $1 ORDER BY id LIMIT 1")
	return stmt.Get(out, id)
}

func (db *DB) getItems(feedID int64, limit int, out *[]Item) error {
	stmt := db.prepare("getItems", "SELECT * FROM items_view ORDER BY id LIMIT $1")
	return stmt.Select(out, limit)
}

func (db *DB) markItemRead(itemID int64, read bool) (err error) {
	if read {
		stmt := db.prepare("markItemRead", "INSERT INTO item_reads ( item_id ) VALUES ( $1 )")
		_, err = stmt.Exec(itemID)
	} else {
		stmt := db.prepare("markItemUnread", "DELETE FROM item_reads WHERE item_id = $1")
		_, err = stmt.Exec(itemID)
	}

	return
}
