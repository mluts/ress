package sqlite

import (
	"github.com/rubenv/sql-migrate"
)

// Migrations holds sqlite migrations required for github.com/mluts/ress
var Migrations = &migrate.MemoryMigrationSource{
	Migrations: []*migrate.Migration{
		&migrate.Migration{
			Id: "1",
			Up: []string{`
			CREATE TABLE IF NOT EXISTS feeds
				(
					id INTEGER PRIMARY KEY,
					link TEXT NOT NULL CHECK(length(link) > 0),
					title TEXT NOT NULL DEFAULT "",
					author TEXT NOT NULL DEFAULT "",
					active BOOLEAN	NOT NULL DEFAULT TRUE,
					error TEXT NOT NULL DEFAULT "",
					published DATETIME DEFAULT 0,
					updated DATETIME DEFAULT 0,
					created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
				);

			-- represents a feed image
			CREATE TABLE IF NOT EXISTS feed_images
				(
					id INTEGER PRIMARY KEY,
					feed_id INTEGER REFERENCES feeds(id) ON DELETE CASCADE,
					url TEXT NOT NULL CHECK(length(url) > 0),
					title TEXT NOT NULL,
					created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
				);

			-- feed -> items, should be deleted if feed was deleted
			CREATE TABLE IF NOT EXISTS items
				(
					id INTEGER PRIMARY KEY,
					feed_id INTEGER REFERENCES feeds(id) ON DELETE CASCADE,
					link TEXT NOT NULL CHECK(length(link) > 0),
					title TEXT NOT NULL CHECK(length(title) > 0),
					content TEXT NOT NULL DEFAULT "",
					description TEXT NOT NULL DEFAULT "",
					guid TEXT NOT NULL,
					author TEXT NOT NULL DEFAULT "",
					updated DATETIME DEFAULT 0,
					published DATETIME DEFAULT 0,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
				);

			-- Holds items read/uread marks
			CREATE TABLE IF NOT EXISTS item_reads
				(
					id INTEGER PRIMARY KEY,
					item_id INTEGER REFERENCES items(id) ON DELETE CASCADE
				);

			-- Holds image for item
			CREATE TABLE IF NOT EXISTS item_images
				(
					id INTEGER PRIMARY KEY,
					item_id INTEGER REFERENCES items(id) ON DELETE CASCADE,
					url TEXT NOT NULL CHECK(length(url) > 0),
					title TEXT NOT NULL DEFAULT "",
					created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
				);

			-- Abstraction on the top of "items"
			CREATE VIEW IF NOT EXISTS items_view AS
				SELECT items.*,
					NOT ifnull(item_reads.id, 0) AS unread,
					ifnull(item_images.id, 0) AS "item_image.id",
					ifnull(item_images.url, "") AS "item_image.url",
					ifnull(item_images.title, "") AS "item_image.title"
					FROM items
					LEFT JOIN item_reads ON item_reads.item_id = items.id
					LEFT JOIN item_images ON item_images.item_id = items.id;

			CREATE VIEW IF NOT EXISTS feeds_view AS
				SELECT feeds.*,
					ifnull(feed_images.id, 0) AS "feed_image.id",
					ifnull(feed_images.url, "") AS "feed_image.url",
					ifnull(feed_images.title, "") AS "feed_image.title"
				FROM feeds
				LEFT JOIN feed_images ON feed_images.feed_id = feeds.id;


			-- items.link should not be duplicated within same feed
			CREATE UNIQUE INDEX IF NOT EXISTS
				feed_item_link ON items (feed_id, link);

			-- items.guid should not be duplicated within same feed
			CREATE UNIQUE INDEX IF NOT EXISTS
				feed_item_guid ON items(feed_id, guid);

			CREATE UNIQUE INDEX IF NOT EXISTS
				item_image_item ON item_images(item_id);

			-- feeds.link should not be duplicated
			CREATE UNIQUE INDEX IF NOT EXISTS
				feed_link ON feeds(link);

			-- item read/unread should not be duplicated
			CREATE UNIQUE INDEX IF NOT EXISTS
				item_read ON item_reads(id, item_id);

			CREATE UNIQUE INDEX IF NOT EXISTS
				feed_image_feed ON feed_images(feed_id);

			-- feeds.updated_at trigger
			CREATE TRIGGER IF NOT EXISTS
				feeds_updated_at AFTER UPDATE ON feeds
				BEGIN
					UPDATE feeds SET updated_at = CURRENT_TIMESTAMP
						WHERE id = NEW.id;
				END;

			-- items.updated_at trigger
			CREATE TRIGGER IF NOT EXISTS
				items_updated_at AFTER UPDATE ON items
				BEGIN
					UPDATE items SET updated_at = CURRENT_TIMESTAMP
						WHERE id = NEW.id;
				END;

			-- feed_images.updated_at trigger
			CREATE TRIGGER IF NOT EXISTS
				feed_images_updated_at AFTER UPDATE ON feed_images
				BEGIN
					UPDATE feed_images SET updated_at = CURRENT_TIMESTAMP
						WHERE id = NEW.id;
				END;
			`},
		},
	},
}
