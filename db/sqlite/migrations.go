package sqlite

import (
	"github.com/rubenv/sql-migrate"
)

var Migrations = &migrate.MemoryMigrationSource{
	Migrations: []*migrate.Migration{
		&migrate.Migration{
			Id: "1",
			Up: []string{`
			CREATE TABLE IF NOT EXISTS feeds
				(
					id INTEGER PRIMARY KEY,
					link TEXT NOT NULL CHECK(length(link) > 0),
					title TEXT NOT NULL CHECK(length(title) > 0),
					author TEXT,
					active BOOLEAN DEFAULT TRUE,
					error TEXT,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
				);

			CREATE TABLE IF NOT EXISTS items
				(
					id INTEGER PRIMARY KEY,
					feed_id REFERENCES feeds(id) ON DELETE CASCADE,
					link TEXT NOT NULL CHECK(length(link) > 0),
					title TEXT NOT NULL CHECK(length(title) > 0),
					content TEXT,
					author TEXT,
					updated DATETIME,
					published DATETIME,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
				);

			CREATE UNIQUE INDEX IF NOT EXISTS
				feed_item_link ON items ( feed_id, link );

			CREATE UNIQUE INDEX IF NOT EXISTS
				feed_link ON feeds(link);

			CREATE TRIGGER IF NOT EXISTS
				feeds_updated_at AFTER UPDATE ON feeds
				BEGIN
					UPDATE feeds SET updated_at = CURRENT_TIMESTAMP
						WHERE id = NEW.id;
				END;

			CREATE TRIGGER IF NOT EXISTS
				items_updated_at AFTER UPDATE ON items
				BEGIN
					UPDATE items SET updated_at = CURRENT_TIMESTAMP
						WHERE id = NEW.id;
				END;
			`},
		},
	},
}
