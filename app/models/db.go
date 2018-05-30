package models

import "github.com/jinzhu/gorm"

const connectionString = "postgresql://postgres:development@db:5432/stereodose?sslmode=disable"

var db *gorm.DB

func GetDB() *gorm.DB {
	if db != nil {
		return db
	}
	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		panic(err.Error())
	}
	return db
}
