package main

import (
	// "database/sql"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mluts/ress/db/sqlite"
	"github.com/rubenv/sql-migrate"
	"time"
)

type DB struct {
	*sqlx.DB
}

type Feed struct {
	ID int64 `db:"id"`

	Title  string `db:"title"`
	Link   string `db:"link"`
	Author string `db:"author"`

	Error     string    `db:"error"`
	Active    bool      `db:"active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Item struct {
	ID     int64 `db:"id"`
	FeedID int64 `db:"feed_id"`

	Title       string `db:"title"`
	Link        string `db:"link"`
	Description string `db:"description"`
	Content     string `db:"content"`
	Author      string `db:"author"`

	Updated   *time.Time `db:"updated"`
	Published *time.Time `db:"published"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

func OpenDatabase(dialect, dest string) (db *DB, err error) {
	var (
		sqldb  *sql.DB
		sqlxdb *sqlx.DB
	)
	if dialect != "sqlite3" {
		db, err = nil, fmt.Errorf("%s is not supported", dialect)
	}

	sqldb, err = sql.Open(dialect, dest)

	// Foreign keys are disabled by default in sqlite3
	if _, err = sqldb.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return
	}

	if _, err = migrate.Exec(sqldb, dialect, sqlite.Migrations, migrate.Up); err != nil {
		return
	}

	sqlxdb = sqlx.NewDb(sqldb, dialect)
	db = &DB{sqlxdb}

	return
}

func (db *DB) feedPresent(link string) bool {
	var count int64
	db.Get(&count, "SELECT count(id) FROM feeds WHERE link = $1", link)
	return count > 0
}

func (db *DB) createFeed(feed *Feed) (id int64, err error) {
	var result sql.Result

	result, err = db.NamedExec(`
		INSERT INTO feeds (link, title, error, active)
			VALUES (:link, :title, :error, :active)
	`, feed)

	if err != nil {
		return
	}

	id, err = result.LastInsertId()

	return
}

func (db *DB) updateFeed(id int64, feed *Feed) error {
	feed.ID = id

	_, err := db.NamedExec(`
		UPDATE feeds SET
			(link, title, error, active)
			VALUES (:link, :title, :error, :active)
			WHERE id = :id
	`, feed)

	return err
}

func (db *DB) createOrUpdateFeedByLink(link string, feed *Feed) (int64, error) {
	var out Feed
	if err := db.getFeedByLink(link, &out); err != nil {
		return db.createFeed(feed)
	}

	return out.ID, db.updateFeed(out.ID, feed)
}

func (db *DB) deleteFeed(id int64) error {
	_, err := db.Exec("DELETE FROM feeds WHERE id = $1", id)
	return err
}

func (db *DB) getFeed(id int64, out *Feed) error {
	return db.Get(out, "SELECT * FROM feeds WHERE id = $1 ORDER BY id LIMIT 1", id)
}

func (db *DB) getFeedByLink(link string, out *Feed) error {
	return db.Get(out, "SELECT * FROM feeds WHERE link = $1 ORDER BY id LIMIT 1", link)
}

func (db *DB) getFeeds(limit int, out *[]Feed) error {
	return db.Select(out, "SELECT * FROM feeds ORDER BY id LIMIT $1", limit)
}

func (db *DB) createItem(feedID int64, item *Item) error {
	item.FeedID = feedID

	_, err := db.NamedExec(`
		INSERT INTO items (feed_id, title, link, description, content)
			VALUES (:feed_id, :title, :link, :description, :content)
	`, item)
	return err
}

func (db *DB) updateItem(item *Item) error {
	_, err := db.NamedExec(`
		UPDATE items SET (title, link, description, content)
			VALUES (link, description, content)
	`, item)
	return err
}

func (db *DB) deleteItem(id int64) error {
	_, err := db.Exec("DELETE FROM items WHERE id = $1", id)
	return err
}

func (db *DB) getItem(id int64, out *Item) error {
	return db.Get(out, "SELECT * FROM items WHERE id = $1 ORDER BY id LIMIT 1", id)
}

func (db *DB) getItems(feedID int64, limit int, out *[]Item) error {
	return db.Select(out, "SELECT * FROM items ORDER BY id LIMIT $1", limit)
}
