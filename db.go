package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// DB provides a database interface
type DB struct {
	db *gorm.DB
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
