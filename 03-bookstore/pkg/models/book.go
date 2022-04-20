package models

import (
	"03-bookstore/pkg/config"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

type Book struct {
	gorm.Model
	Title     string `json:"title"`
	Author    string `json:"author"`
	Publisher string `json:"publisher"`
}

func init() {
	db = config.GetDatabase()
	db.AutoMigrate(&Book{})
}

func (b *Book) CreateBook() *Book {
	db.Create(&b)
	return b
}

func GetBooks() []Book {
	var books []Book
	db.Find(&books)
	return books
}

func GetBook(id int64) *Book {
	var book Book
	db.First(&book, id)
	return &book
}

func UpdateBook(id int64, newBook Book) *Book {
	var book Book
	db.First(&book, id)
	if book.ID == 0 {
		return nil
	}
	newBook.ID = book.ID
	newBook.CreatedAt = book.CreatedAt
	db.Save(&newBook)
	return &newBook
}

func DeleteBook(id int64) {
	db.Where("ID = ?", id).Delete(&Book{})
}
