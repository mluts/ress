package main

import (
	// "database/sql"
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mluts/ress/db/sqlite"
	"github.com/rubenv/sql-migrate"
	"time"
)

// DB provides a database interface
type DB struct {
	db *gorm.DB
}

type _DB struct {
	*sqlx.DB
}

type _Feed struct {
	ID uint `db:"id"`

	Title  sql.NullString `db:"title"`
	Link   sql.NullString `db:"link"`
	Author sql.NullString `db:"author"`

	Error     sql.NullString `db:"error"`
	Active    bool           `db:"active"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

type _Item struct {
	ID     uint `db:"id"`
	FeedID uint `db:"feed_id"`

	Title       string `db:"title"`
	Link        string `db:"link"`
	Description string `db:"description"`
	Content     string `db:"content"`
}

// OpenDatabase opens and initializes database
func OpenDatabase(dialect, dest string) (*DB, error) {
	db, err := gorm.Open(dialect, dest)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Feed{})
	db.AutoMigrate(&Item{})

	return &DB{db}, nil
}

func _OpenDatabase(dialect, dest string) (db *_DB, err error) {
	var (
		sqldb  *sql.DB
		sqlxdb *sqlx.DB
	)
	if dialect != "sqlite3" {
		db, err = nil, fmt.Errorf("%s is not supported", dialect)
	}

	sqldb, err = sql.Open(dialect, dest)
	_, err = migrate.Exec(sqldb, dialect, sqlite.Migrations, migrate.Up)
	if err != nil {
		return
	}

	sqlxdb = sqlx.NewDb(sqldb, dialect)
	db = &_DB{sqlxdb}

	return
}

func (db *_DB) feedPresent(link string) bool {
	var count int64
	db.Get(&count, "SELECT count(id) FROM feeds WHERE link = $1", link)
	return count > 0
}

func (db *_DB) createFeed(feed *_Feed) (id int64, err error) {
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

func (db *_DB) updateFeed(feed *_Feed) error {
	_, err := db.NamedExec(`
		UPDATE feeds SET
			(link, title, error, active)
			VALUES (:link, :title, :error, :active)
	`, feed)

	return err
}

func (db *_DB) deleteFeed(id int64) error {
	_, err := db.Exec("DELETE FROM feeds WHERE id = $1", id)
	return err
}

func (db *_DB) getFeed(id int64, out *_Feed) error {
	return db.Get(out, "SELECT * FROM feeds WHERE id = $1 ORDER BY id LIMIT 1", id)
}

func (db *_DB) getFeeds(limit int, out *[]_Feed) error {
	return db.Select(out, "SELECT * FROM feeds ORDER BY id LIMIT $1", limit)
}

func (db *_DB) createItem(feed *_Feed, item *_Item) error {
	item.FeedID = feed.ID

	_, err := db.NamedExec(`
		INSERT INTO items (feed_id, title, link, description, content)
			VALUES (:feed_id, :title, :link, :description, :content)
	`, item)
	return err
}
