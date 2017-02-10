package app

func (db *DB) createItem(feedID int64, item *Item) (int64, error) {
	item.FeedID = feedID

	stmt := db.prepareNamed(
		"createItem",
		`INSERT INTO items
			(
				feed_id,
				guid,
				title,
				link,
				description,
				content
			)
		VALUES
			(
				:feed_id,
				:guid,
				:title,
				:link,
				:description,
				:content
			)`)

	result, err := stmt.Exec(item)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (db *DB) updateItem(id int64, item *Item) error {
	item.ID = id

	stmt := db.prepareNamed(
		"updateItem",
		`UPDATE items SET (title, link, description, content) =
		(:title, :link, :description, :content)`)

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
	stmt := db.prepare("getItems",
		"SELECT * FROM items_view WHERE feed_id = $1 ORDER BY id LIMIT $2")
	return stmt.Select(out, feedID, limit)
}

func (db *DB) getItemsCount(feedID int64, count *int) error {
	stmt := db.prepare("getItemsCount",
		"SELECT COUNT(id) FROM items WHERE feed_id = $1")
	return stmt.Get(count, feedID)
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

func (db *DB) getItemByGUID(feedID int64, guid string, item *Item) error {
	stmt := db.prepare("getItemByGUID",
		"SELECT * FROM items WHERE feed_id = $1 AND guid = $2 LIMIT 1")
	return stmt.Get(item, feedID, guid)
}
