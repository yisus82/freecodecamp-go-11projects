package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DB *gorm.DB

func connect() *gorm.DB {
	db, err := gorm.Open("sqlite3", "./leads.db")
	if err != nil {
		panic("failed to connect database")
	}
	DB = db
	return db
}

func GetDatabase() *gorm.DB {
	if DB == nil {
		connect()
	}
	return DB
}
