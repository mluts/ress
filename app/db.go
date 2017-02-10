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

// SQLNoLimit used as LIMIT argument means ALL records
const SQLNoLimit = -1

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

	Error  string `db:"error"`
	Active bool   `db:"active"`

	Published *time.Time `db:"published"`
	Updated   *time.Time `db:"updated"`

	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`

	Image *Image `db:"feed_image"`
}

// Image represents a simple picture with url and optional title
type Image struct {
	ID    int64  `db:"id"`
	Title string `db:"title"`
	URL   string `db:"url"`
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

	GUID string `db:"guid"`

	Updated   *time.Time `db:"updated"`
	Published *time.Time `db:"published"`

	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`

	Image *Image `db:"item_image"`
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
