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
					id 				 INTEGER PRIMARY KEY,
					link 			 TEXT 		NOT NULL CHECK(length(link) > 0),
					title 		 TEXT 		NOT NULL CHECK(length(title) > 0),
					author 		 TEXT 		NOT NULL DEFAULT "",
					active 		 BOOLEAN 	NOT NULL DEFAULT TRUE,
					error 		 TEXT 		NOT NULL DEFAULT "",
					created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
				);

			CREATE TABLE IF NOT EXISTS items
				(
					id 					INTEGER PRIMARY KEY,
					feed_id 		INTEGER REFERENCES feeds(id) ON DELETE CASCADE,
					link 				TEXT NOT NULL CHECK(length(link) > 0),
					title 			TEXT NOT NULL CHECK(length(title) > 0),
					content 		TEXT NOT NULL DEFAULT "",
					description TEXT NOT NULL DEFAULT "",
					author 			TEXT NOT NULL DEFAULT "",
					updated 		DATETIME NOT NULL DEFAULT 0,
					published 	DATETIME NOT NULL DEFAULT 0,
					created_at 	DATETIME DEFAULT CURRENT_TIMESTAMP,
					updated_at 	DATETIME DEFAULT CURRENT_TIMESTAMP
				);

			CREATE TABLE IF NOT EXISTS item_reads
				(
					id 			INTEGER PRIMARY KEY,
					item_id INTEGER REFERENCES items(id) ON DELETE CASCADE
				);

			CREATE VIEW IF NOT EXISTS items_view AS
				SELECT items.*, NOT ifnull(item_reads.id, 0) AS unread FROM items
					LEFT JOIN item_reads ON item_reads.item_id = items.id;

			CREATE UNIQUE INDEX IF NOT EXISTS
				feed_item_link ON items ( feed_id, link );

			CREATE UNIQUE INDEX IF NOT EXISTS
				feed_link ON feeds(link);

			CREATE UNIQUE INDEX IF NOT EXISTS
				item_read ON item_reads(id, item_id);

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
