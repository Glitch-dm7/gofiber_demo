package models

import "gorm.io/gorm"

type Book struct{
	ID					uint8			`gorm:"primary key;autoincrement" json:"id"`
	Author			*string		`json:"author"`	
	Title				*string		`json:"title"`
	Publisher		*string		`json:"publisher"`
}

func MigrateBooks(db *gorm.DB) error {
	err := db.AutoMigrate(&Book{})
	return err
}