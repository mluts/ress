package app

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	// sqlite3 support
	_ "github.com/mattn/go-sqlite3"
	"github.com/mluts/ress/db/sqlite"
	"github.com/rubenv/sql-migrate"
	"time"
)

// DB is a database wrapper to provide application specific methods for it
type DB struct {
	db    *sqlx.DB
	named map[string]*sqlx.NamedStmt
	stmt  map[string]*sqlx.Stmt
}

// Feed represents a rss feed in our database
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

// Item represents a feed item in our database
type Item struct {
	ID     int64 `db:"id"`
	FeedID int64 `db:"feed_id"`

	Title       string `db:"title"`
	Link        string `db:"link"`
	Description string `db:"description"`
	Content     string `db:"content"`
	Author      string `db:"author"`

	Unread bool `db:"unread"`

	Updated   *time.Time `db:"updated"`
	Published *time.Time `db:"published"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

// OpenDatabase returns a common database connection for given dialect
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
	db = &DB{db: sqlxdb}

	return
}

func (db *DB) prepareNamed(name, query string) (s *sqlx.NamedStmt) {
	if s, ok := db.named[name]; ok {
		return s
	}

	s, err := db.db.PrepareNamed(query)
	if err != nil {
		panic(err)
	}
	return
}

func (db *DB) prepare(name, query string) (s *sqlx.Stmt) {
	if s, ok := db.stmt[name]; ok {
		return s
	}

	s, err := db.db.Preparex(query)
	if err != nil {
		panic(err)
	}
	return
}

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

func (db *DB) createOrUpdateFeedByLink(link string, feed *Feed) (int64, error) {
	var out Feed
	if err := db.getFeedByLink(link, &out); err != nil {
		return db.createFeed(feed)
	}

	return out.ID, db.updateFeed(out.ID, feed)
}

func (db *DB) deleteFeed(id int64) error {
	stmt := db.prepare("deleteFeed", `
		DELETE FROM feeds WHERE id = $1
	`)
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
		`UPDATE items SET (title, link, description, content)
		VALUES (:link, :description, :content)`)

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
