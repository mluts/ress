package app

import (
	"database/sql"
)

func (db *DB) createFeed(link string) (id int64, err error) {
	var result sql.Result

	stmt := db.prepare("createFeed",
		"INSERT INTO feeds (link) VALUES ($1)")

	result, err = stmt.Exec(link)

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
			(title, error, active) =
			(:title, :error, :active)
			WHERE id = :id
	`)

	_, err := stmt.Exec(feed)
	if err != nil {
		return err
	}

	if feed.Image != nil {
		err = db.createOrUpdateFeedImage(id, feed.Image.URL, feed.Image.Title)
	}

	return err
}

func (db *DB) createOrUpdateFeedImage(feedID int64, url, title string) error {
	stmt := db.prepare(
		"createFeedImage",
		`INSERT INTO feed_images (feed_id, url, title) VALUES ($1, $2, $3)`)

	switch err := db.getFeedImage(feedID, &Image{}); err {
	case sql.ErrNoRows:
		_, err = stmt.Exec(feedID, url, title)
		return nil
	case nil:
		return db.updateFeedImage(feedID, url, title)
	default:
		return err
	}
}

func (db *DB) updateFeedImage(feedID int64, url, title string) error {
	stmt := db.prepare(
		"updateFeedImage",
		`UPDATE feed_images SET (url, title) = ($1, $2) WHERE feed_id = $3`)
	_, err := stmt.Exec(url, title, feedID)
	return err
}

func (db *DB) getFeedImage(feedID int64, out *Image) error {
	stmt := db.prepare(
		"getFeedImage",
		`SELECT id, url, title FROM feed_images
			WHERE feed_id = $1 LIMIT 1`)
	err := stmt.Get(out, feedID)
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
		"SELECT * FROM feeds_view WHERE id = $1 ORDER BY id LIMIT 1")
	return stmt.Get(out, id)
}

func (db *DB) getFeedByLink(link string, out *Feed) error {
	stmt := db.prepare(
		"getFeedByLink",
		"SELECT * FROM feeds_view WHERE link = $1 ORDER BY id LIMIT 1")
	return stmt.Get(out, link)
}

func (db *DB) getFeeds(limit int, out *[]Feed) error {
	stmt := db.prepare(
		"getFeeds",
		"SELECT * FROM feeds_view ORDER BY id LIMIT $1")
	return stmt.Select(out, limit)
}
