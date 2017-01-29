package app

import (
	"database/sql"
)

func (db *DB) createFeed(feed *Feed) (id int64, err error) {
	var result sql.Result

	stmt := db.prepareNamed("createFeed", `
		INSERT INTO feeds (link, title, error, active)
			VALUES (:link, :title, :error, :active)
	`)

	result, err = stmt.Exec(feed)

	if err != nil {
		return
	}

	id, err = result.LastInsertId()

	return
}

func (db *DB) updateFeed(id int64, feed *Feed) error {
	feed.ID = id

	stmt := db.prepareNamed("updateFeed", `
		UPDATE feeds SET
			(link, title, error, active)
			VALUES (:link, :title, :error, :active)
			WHERE id = :id
	`)

	_, err := stmt.Exec(feed)

	return err
}

func (db *DB) deleteFeed(id int64) error {
	stmt := db.prepare("deleteFeed", "DELETE FROM feeds WHERE id = $1")
	_, err := stmt.Exec(id)
	return err
}

func (db *DB) getFeed(id int64, out *Feed) error {
	stmt := db.prepare(
		"getFeed",
		"SELECT * FROM feeds WHERE id = $1 ORDER BY id LIMIT 1")
	return stmt.Get(out, id)
}

func (db *DB) getFeedByLink(link string, out *Feed) error {
	stmt := db.prepare(
		"getFeedByLink",
		"SELECT * FROM feeds WHERE link = $1 ORDER BY id LIMIT 1")
	return stmt.Get(out, link)
}

func (db *DB) getFeeds(limit int, out *[]Feed) error {
	stmt := db.prepare(
		"getFeeds",
		"SELECT * FROM feeds ORDER BY id LIMIT $1")
	return stmt.Select(out, limit)
}
