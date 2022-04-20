package config

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func ConnectDatabase() {
	d, err := gorm.Open("mysql", "root:liceo@tcp(127.0.0.1:3306)/bookstore?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	db = d
}

func GetDatabase() *gorm.DB {
	if db == nil {
		ConnectDatabase()
	}
	return db
}
